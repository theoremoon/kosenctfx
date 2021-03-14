package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/xerrors"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
	"github.com/theoremoon/kosenctfx/scoreserver/util"
)

const (
	cacheDuration     = 1 * time.Minute
	challengesJSONKey = "challengesJSONKey"
	rankingJSONKey    = "rankingJSONKey"
)

func (s *server) registerHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(struct {
			Teamname    string
			Email       string
			Password    string
			CountryCode string `json:"country"`
		})
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": InvalidRequestMessage,
			})
		}

		if _, err := s.app.RegisterTeam(req.Teamname, req.Password, req.Email, req.CountryCode); err != nil {
			return errorHandle(c, err)
		}
		return messageHandle(c, RegisteredMessage)
	}

}

func (s *server) loginHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(struct {
			Teamname string
			Password string
		})
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": InvalidRequestMessage,
			})
		}

		token, err := s.app.Login(req.Teamname, req.Password, c.RealIP())
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		c.SetCookie(s.tokenCookie(token))
		return messageHandle(c, LoginMessage)
	}
}

func (s *server) logoutHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.SetCookie(s.removeTokenCookie())
		return messageHandle(c, LogoutMessage)
	}
}

func (s *server) infoHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ret := make(map[string]interface{})
		team, _ := s.getLoginTeam(c)
		if team != nil {
			ret["teamname"] = team.Teamname
			ret["teamid"] = team.ID
		}
		conf, err := s.app.GetCTFConfig()
		if err != nil {
			return errorHandle(c, err)
		}
		ret["ctf_start"] = conf.StartAt
		ret["ctf_end"] = conf.EndAt
		ret["ctf_name"] = conf.CTFName

		return c.JSON(http.StatusOK, ret)
	}
}

/// ログインしているかどうか、CTF開催中かどうかで挙動が変わる
/// ログインしていない -> notificationとrankingを返す
/// ログインしている   -> 問題情報も返す
/// CTF開催中          -> ranking / 問題情報は redisにcacheしておいて、expiredになったら計算してcacheし直す
/// CTF開催中でない    -> 毎回計算する
func (s *server) infoUpdateHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		t, _ := s.getLoginTeam(c)
		refresh := c.QueryParam("refresh")

		conf, err := s.app.GetCTFConfig()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		status := service.CalcCTFStatus(conf)

		ret := make(map[string]interface{})

		// TODO notification
		// cache を使う
		if refresh == "" && status == service.CTFRunning {
			challenges, ranking, err := s.getCacheInfo()
			if err != nil {
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}

			if challenges != "" && ranking != "" {
				var cs []*service.Challenge
				var scoreboard *service.Scoreboard
				err1 := json.Unmarshal([]byte(challenges), &cs)
				err2 := json.Unmarshal([]byte(ranking), &scoreboard)
				if err1 == nil && err2 == nil {
					ret["challenges"] = cs
					ret["ranking"] = scoreboard
				}
			}
		}

		_, exist1 := ret["challenges"]
		_, exist2 := ret["ranking"]
		if !exist1 || !exist2 {
			// CTF開催中またはCTF終了後なら、公開されている問題を読み込むが、そうでない（invalid or 開催前）なら読み込まない
			chals := make([]*model.Challenge, 0)
			if status == service.CTFRunning || status == service.CTFEnded {
				chals, err = s.app.ListOpenedRawChallenges()
				if err != nil {
					return errorHandle(c, xerrors.Errorf(": %w", err))
				}
			}
			teams, err := s.app.ListTeams()
			if err != nil {
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}

			challenges, ranking, err := s.app.ScoreFeed(chals, teams)
			if err != nil {
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}
			for i := range challenges {
				challenges[i].Flag = ""
			}
			ret["challenges"] = challenges
			ret["ranking"] = ranking

			// cacheする
			if status == service.CTFRunning {
				bytes1, err1 := json.Marshal(challenges)
				bytes2, err2 := json.Marshal(ranking)
				if err1 == nil && err2 == nil {
					s.setCacheInfo(string(bytes1), string(bytes2))
				}
			}
		}

		if t == nil || status == service.CTFNotStarted {
			delete(ret, "challenges")
		}
		return c.JSON(http.StatusOK, ret)
	}
}

func (s *server) passwordresetRequestHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(struct {
			Email string
		})
		if err := c.Bind(req); err != nil {
			return errorHandle(c, err)
		}

		if err := s.app.PasswordResetRequest(req.Email); err != nil {
			return errorHandle(c, err)
		}
		return messageHandle(c, PasswordResetEmailSentMessage)
	}
}

func (s *server) passwordresetHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(struct {
			Token       string
			NewPassword string `json:"new_password"`
		})
		if err := c.Bind(req); err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		if err := s.app.PasswordReset(req.Token, req.NewPassword); err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		return messageHandle(c, PasswordUpdateMessage)
	}
}

func (s *server) profileUpdateHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		lc := c.(*loginContext)
		req := new(struct {
			Teamname    string `json:"teamname"`
			Password    string `json:"password"`
			CountryCode string `json:"country"`
		})
		if err := c.Bind(req); err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		if req.Teamname != "" {
			if err := s.app.UpdateTeamname(lc.Team, req.Teamname); err != nil {
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}
		}

		if req.Password != "" {
			if err := s.app.PasswordUpdate(lc.Team, req.Password); err != nil {
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}
		}

		if req.CountryCode != "" {
			if err := s.app.UpdateCountry(lc.Team, req.CountryCode); err != nil {
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}
		}
		return messageHandle(c, ProfileUpdateMessage)
	}
}

func (s *server) teamHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		teamIDstr := c.Param("id")
		teamID, err := strconv.ParseUint(teamIDstr, 10, 32)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		team, err := s.app.GetTeamByID(uint32(teamID))
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		res := map[string]interface{}{
			"teamname": team.Teamname,
			"teamid":   team.ID,
			"country":  team.CountryCode,
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (s *server) submitHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		lc := c.(*loginContext)
		req := new(struct {
			Flag string `json:"flag"`
		})
		if err := c.Bind(req); err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		// check Submission Lock
		submittable, err := s.app.CheckSubmittable(lc.Team.ID)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		if !submittable {
			return errorHandle(c, xerrors.Errorf(": %w", service.NewErrorMessage(SubmissionLockedMessage)))
		}

		conf, err := s.app.GetCTFConfig()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		ctfStatus := service.CalcCTFStatus(conf)

		// flag submission
		flag := strings.Trim(req.Flag, " ")
		challenge, correct, valid, err := s.app.SubmitFlag(lc.Team, lc.RealIP(), flag, ctfStatus == service.CTFRunning)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		if valid {
			s.SolveLogWebhook.Post(fmt.Sprintf(
				ValidSubmissionSystemMessage,
				util.DiscordString(lc.Team.Teamname),
				challenge.Name,
			))
			s.AdminWebhook.Post(fmt.Sprintf(
				ValidSubmissionAdminMessage,
				util.DiscordString(lc.Team.Teamname),
				challenge.Name,
				util.DiscordString(req.Flag),
			))
			return messageHandle(c, fmt.Sprintf(ValidSubmissionMessage, challenge.Name))
		} else if correct {
			s.AdminWebhook.Post(fmt.Sprintf(
				CorrectSubmissionAdminMessage,
				util.DiscordString(lc.Team.Teamname),
				challenge.Name,
				util.DiscordString(req.Flag),
			))
			return messageHandle(c, fmt.Sprintf(CorrectSubmissionMessage, challenge.Name))
		} else {
			// wrong count
			count, err := s.app.GetWrongCount(lc.Team.ID, time.Duration(conf.LockDuration)*time.Second)
			if err != nil {
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}
			if count >= conf.LockCount {
				if err := s.app.LockSubmission(lc.Team.ID, time.Duration(conf.LockSecond)*time.Second); err != nil {
					return errorHandle(c, xerrors.Errorf(": %w", err))
				}
			}

			s.AdminWebhook.Post(fmt.Sprintf(
				WrongSubmissionAdminMessage,
				util.DiscordString(lc.Team.Teamname),
				util.DiscordString(req.Flag),
			))
			return errorHandle(c, service.NewErrorMessage(WrongSubmissionMessage))
		}
	}
}

func (s *server) scoreEmulateHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		maxCountStr := c.QueryParam("maxCount")
		maxCount, err := strconv.Atoi(maxCountStr)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		if maxCount < 0 {
			return errorHandle(c, xerrors.Errorf(": %w", service.NewErrorMessage(ScoreEmulateMaxCountTooSmallMessage)))
		}

		expr := c.QueryParam("expr")
		if expr == "" {
			conf, err := s.app.GetCTFConfig()
			if err != nil {
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}
			expr = conf.ScoreExpr
		}

		scores := make([]int, maxCount+1)
		for i := 0; i <= maxCount; i++ {
			scores[i], err = service.CalcChallengeScore(i, expr)
			if err != nil {
				return errorHandle(c, xerrors.Errorf(": %w", service.NewErrorMessage(err.Error())))
			}
		}
		return c.JSON(http.StatusOK, scores)
	}
}

func (s *server) getConfigHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		conf, err := s.app.GetCTFConfig()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		ret := make(map[string]interface{})
		ret["ctf_name"] = conf.CTFName
		ret["start_at"] = conf.StartAt
		ret["end_at"] = conf.EndAt
		ret["score_expr"] = conf.ScoreExpr
		ret["register_open"] = conf.RegisterOpen
		ret["ctf_open"] = conf.CTFOpen
		ret["lock_second"] = conf.LockSecond
		ret["lock_duration"] = conf.LockDuration
		ret["lock_count"] = conf.LockCount

		return c.JSON(http.StatusOK, ret)
	}
}

func (s *server) ctfConfigHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(struct {
			Name         string `json:"name"`
			StartAt      int64  `json:"start_at"`
			EndAt        int64  `json:"end_at"`
			RegisterOpen bool   `json:"register_open"`
			CTFOpen      bool   `json:"ctf_open"`
			LockCount    int64  `json:"lock_count"`
			LockSecond   int    `json:"lock_second"`
			LockDuration int    `json:"lock_duration"`
			ScoreExpr    string `json:"score_expr"`
		})
		if err := c.Bind(req); err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		// score_exprのチェック
		_, err := service.CalcChallengeScore(10, req.ScoreExpr)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		conf, err := s.app.GetCTFConfig()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		conf.CTFName = req.Name
		conf.StartAt = req.StartAt
		conf.EndAt = req.EndAt
		conf.RegisterOpen = req.RegisterOpen
		conf.CTFOpen = req.CTFOpen
		conf.LockCount = req.LockCount
		conf.LockSecond = req.LockSecond
		conf.LockDuration = req.LockDuration
		conf.ScoreExpr = req.ScoreExpr
		if err := s.app.SetCTFConfig(conf); err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		return c.JSON(http.StatusOK, ConfigUpdateMessage)
	}
}

func (s *server) openChallengeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(struct {
			Name string `json:"name"`
		})
		if err := c.Bind(req); err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		chal, err := s.app.GetRawChallengeByName(req.Name)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		if chal.IsOpen {
			return c.JSON(http.StatusOK, fmt.Sprintf(ChallengeAlreadyOpenedTemplate, chal.Name))
		}

		if err := s.app.OpenChallenge(chal.ID); err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		conf, err := s.app.GetCTFConfig()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		status := service.CalcCTFStatus(conf)
		if status == service.CTFRunning {
			s.TaskOpenWebhook.Post(fmt.Sprintf(ChallengeOpenSystemMessage, chal.Name))
		}
		s.AdminWebhook.Post(fmt.Sprintf(ChallengeOpenAdminMessage, chal.Name))
		return c.JSON(http.StatusOK, fmt.Sprintf(ChallengeOpenTemplate, chal.Name))
	}
}

func (s *server) closeChallengeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(struct {
			Name string `json:"name"`
		})
		if err := c.Bind(req); err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		chal, err := s.app.GetRawChallengeByName(req.Name)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		if !chal.IsOpen {
			return c.JSON(http.StatusOK, fmt.Sprintf(ChallengeAlreadyClosedTemplate, chal.Name))
		}

		if err := s.app.CloseChallenge(chal.ID); err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		s.AdminWebhook.Post(fmt.Sprintf(ChallengeClosedAdminMessage, chal.Name))
		return c.JSON(http.StatusOK, fmt.Sprintf(ChallengeCloseTemplate, chal.Name))
	}
}

func (s *server) updateChallengeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(struct {
			ID          uint32
			Name        string
			Flag        string
			Description string
			Author      string
			IsSurvey    bool `json:"is_survey"`
			Tags        []string
			Attachments []service.Attachment
		})
		if err := c.Bind(req); err != nil {
			return errorHandle(c, err)
		}
		err := s.app.UpdateChallenge(
			req.ID,
			&service.Challenge{
				Name:        req.Name,
				Flag:        req.Flag,
				Description: req.Description,
				Author:      req.Author,
				IsSurvey:    req.IsSurvey,
				Tags:        req.Tags,
				Attachments: req.Attachments,
			})
		if err != nil {
			return errorHandle(c, err)
		}
		return c.JSON(http.StatusOK, fmt.Sprintf(ChallengeUpdateTemplate, req.Name))
	}
}

func (s *server) newChallengeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(struct {
			Name        string
			Flag        string
			Description string
			Author      string
			IsSurvey    bool `json:"is_survey"`
			Tags        []string
			Attachments []service.Attachment
		})
		if err := c.Bind(req); err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		if chal, err := s.app.GetRawChallengeByName(req.Name); err == nil {
			// UPDATE
			if err := s.app.UpdateChallenge(chal.ID, &service.Challenge{
				Name:        req.Name,
				Flag:        req.Flag,
				Description: req.Description,
				Author:      req.Author,
				IsSurvey:    req.IsSurvey,
				Tags:        req.Tags,
				Attachments: req.Attachments,
				IsOpen:      chal.IsOpen,
			}); err != nil {
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}
			return c.JSON(http.StatusOK, fmt.Sprintf(ChallengeUpdateTemplate, req.Name))
		} else {
			// ADD
			if err := s.app.AddChallenge(&service.Challenge{
				Name:        req.Name,
				Flag:        req.Flag,
				Description: req.Description,
				Author:      req.Author,
				IsSurvey:    req.IsSurvey,
				Tags:        req.Tags,
				Attachments: req.Attachments,
			}); err != nil {
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}
			return c.JSON(http.StatusOK, fmt.Sprintf(ChallengeAddTemplate, req.Name))
		}
	}
}

func (s *server) listChallengesHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		chals, err := s.app.ListAllRawChallenges()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		teams, err := s.app.ListTeams()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		challenges, ranking, err := s.app.ScoreFeed(chals, teams)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"challenges": challenges,
			"ranking":    ranking,
		})
	}
}

func (s *server) getPresignedURLHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(struct {
			Key string `json:"key"`
		})
		if err := c.Bind(req); err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		if req.Key == "" {
			return errorHandle(c, xerrors.Errorf(": %w", service.NewErrorMessage(PresignedURLKeyRequiredMessage)))
		}

		if s.Bucket == nil {
			return errorHandle(c, xerrors.Errorf(": %w", service.NewErrorMessage(BucketNullMessage)))
		}

		key := uuid.New().String() + "/" + req.Key
		presignedURL, downloadURL, err := s.Bucket.GeneratePresignedURL(key)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"presignedURL": presignedURL,
			"downloadURL":  downloadURL,
		})
	}
}

func (s *server) sqlHandler() echo.HandlerFunc {

	return func(c echo.Context) error {
		req := new(struct {
			Query string `json:"query"`
		})
		if err := c.Bind(req); err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		if req.Query == "" {
			req.Query = "SELECT * FROM information_schema.tables WHERE table_schema=database()"
		}

		cols, rows, err := s.doQuery(req.Query)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"columns": cols,
			"rows":    rows,
		})
	}
}

func (s *server) getChallengesJSON() (string, error) {
	r, err := s.redis.Get(challengesJSONKey).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", xerrors.Errorf(": %w", err)
	}
	return r, nil
}

func (s *server) setChallngesJSON(value string) error {
	err := s.redis.Set(challengesJSONKey, value, cacheDuration).Err()
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (s *server) getCacheInfo() (string, string, error) {
	c, err := s.redis.Get(challengesJSONKey).Result()
	if err == redis.Nil {
		return "", "", nil
	} else if err != nil {
		return "", "", xerrors.Errorf(": %w", err)
	}

	r, err := s.redis.Get(rankingJSONKey).Result()
	if err == redis.Nil {
		return "", "", nil
	} else if err != nil {
		return "", "", xerrors.Errorf(": %w", err)
	}

	return c, r, nil
}

func (s *server) setCacheInfo(challenges, ranking string) error {
	err := s.redis.Set(challengesJSONKey, challenges, cacheDuration).Err()
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	err = s.redis.Set(rankingJSONKey, ranking, cacheDuration).Err()
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	return nil
}

func (s *server) doQuery(query string) ([]string, []map[string]interface{}, error) {
	rows, err := s.db.Raw(query).Rows()
	if err != nil {
		return nil, nil, xerrors.Errorf(": %w", err)
	}

	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, xerrors.Errorf(": %w", err)
	}
	result := make([]map[string]interface{}, 0)
	for rows.Next() {
		row := make([]interface{}, len(cols))
		row_ptr := make([]interface{}, len(cols))
		for i := 0; i < len(row); i++ {
			row_ptr[i] = &row[i]
		}
		if err := rows.Scan(row_ptr...); err != nil {
			return nil, nil, xerrors.Errorf(": %w", err)
		}

		result_row := make(map[string]interface{})
		for i := 0; i < len(cols); i++ {
			switch v := (*row_ptr[i].(*interface{})).(type) {
			case nil:
				result_row[cols[i]] = nil
			case []byte:
				result_row[cols[i]] = string(v)
			default:
				result_row[cols[i]] = v
			}
		}

		result = append(result, result_row)
	}
	return cols, result, nil
}
