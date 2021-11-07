package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
	"github.com/theoremoon/kosenctfx/scoreserver/util"
)

const (
	cacheDuration         = 1 * time.Minute
	challengesJSONKey     = "challengesJSONKey"
	rankingJSONKey        = "rankingJSONKey"
	ClientSeriesMaxTeams  = 20
	sessionSetKey         = "sessionSetKey"
	sessionActiveDuration = 10 * time.Minute
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

		team, _ := s.app.GetLoginTeam(token.Token)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":  LoginMessage,
			"teamname": team.Teamname,
			"team_id":  team.ID,
			"country":  team.CountryCode,
		})
	}
}

func (s *server) logoutHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.SetCookie(s.removeTokenCookie())
		return messageHandle(c, LogoutMessage)
	}
}

func (s *server) accountHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		team, err := s.getLoginTeam(c)
		if err != nil {
			return c.JSON(http.StatusOK, nil)
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"teamname": team.Teamname,
			"team_id":  team.ID,
			"country":  team.CountryCode,
			"is_admin": team.IsAdmin,
		})
	}
}

func (s *server) ctfHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		conf, err := s.app.GetCTFConfig()
		if err != nil {
			return errorHandle(c, err)
		}
		status := service.CalcCTFStatus(conf)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"start_at":      conf.StartAt,
			"end_at":        conf.EndAt,
			"register_open": conf.RegisterOpen,
			"is_open":       conf.CTFOpen,
			"is_running":    status == service.CTFRunning,
			"is_over":       status == service.CTFEnded,
		})
	}
}

func (s *server) scoreboardHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		conf, err := s.app.GetCTFConfig()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		scoreboard, err := s.getScoreboard(conf)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		return c.JSON(http.StatusOK, scoreboard)
	}
}

func (s *server) tasksHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		conf, err := s.app.GetCTFConfig()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		status := service.CalcCTFStatus(conf)

		// CTF始まってないとき問題見せない
		if status == service.CTFNotStarted {
			return errorMessageHandle(c, http.StatusForbidden, CTFNotStartedMessage)
		}

		// CTFがnow-runningでloginしてないとき、non sensitiveな情報しかみせない
		t, _ := s.getLoginTeam(c)
		if status == service.CTFRunning && t == nil {
			challenges, err := s.getNonsensitiveChallenges(conf)
			if err != nil {
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}
			return c.JSON(http.StatusOK, challenges)
		}

		challenges, err := s.getChallenges(conf)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		return c.JSON(http.StatusOK, challenges)
	}
}

/// チームの点数のグラフ出すやつ
/// 長くなりうる配列パラメータを受け取るためにPOST
func (s *server) seriesHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(struct {
			Teams []string `json:"teams"`
		})
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": InvalidRequestMessage,
			})
		}

		conf, err := s.app.GetCTFConfig()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		scoreboard, err := s.getScoreboard(conf)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		// 通常のクライアントにはClientSeriesMaxTeams個までしかエントリを返さない
		topTeamSeries := make([]TeamScoreSeries, 0, ClientSeriesMaxTeams)
		for i, t := range req.Teams {
			if i >= ClientSeriesMaxTeams {
				break
			}
			team_idx := -1
			for idx, team := range scoreboard {
				if team.Teamname == t {
					team_idx = idx
					break
				}
			}
			if team_idx == -1 {
				continue
			}

			series, err := s.getTeamSeries(conf, scoreboard[team_idx].TeamID)
			if err != nil {
				if xerrors.Is(err, redis.Nil) {
					topTeamSeries = append(topTeamSeries, TeamScoreSeries{})
					continue
				}
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}
			topTeamSeries = append(topTeamSeries, series)
		}

		return c.JSON(http.StatusOK, topTeamSeries)
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

		if err := s.app.UpdateCountry(lc.Team, req.CountryCode); err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
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
			"team_id":  team.ID,
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
		challenge, correct, valid, err := s.app.SubmitFlag(lc.Team, lc.RealIP(), flag, ctfStatus == service.CTFRunning, time.Now().Unix())
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

			// 非同期である必要はないが非同期でもいいし、
			// エラーの扱いが楽になるので非同期にしている
			go func() {
				// time seriesを更新する
				_, scoreboard, err := s.refreshCache(conf)
				if err != nil {
					log.Printf("%+v\n", err)
					return
				}
				if err := s.appendScoreSeries(conf, scoreboard, time.Now()); err != nil {
					log.Printf("%+v\n", err)
					return
				}
			}()

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
			Name         string `json:"ctf_name"`
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
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": ConfigUpdateMessage,
		})
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
		s.refreshCache(conf)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": fmt.Sprintf(ChallengeOpenTemplate, chal.Name),
		})
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

		conf, err := s.app.GetCTFConfig()
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
		s.refreshCache(conf)
		s.AdminWebhook.Post(fmt.Sprintf(ChallengeClosedAdminMessage, chal.Name))
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": fmt.Sprintf(ChallengeCloseTemplate, chal.Name),
		})
	}
}

func (s *server) updateChallengeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(struct {
			ID          uint32
			Name        string
			Flag        string
			Category    string
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
				Category:    req.Category,
				Description: req.Description,
				Author:      req.Author,
				IsSurvey:    req.IsSurvey,
				Tags:        req.Tags,
				Attachments: req.Attachments,
			})
		if err != nil {
			return errorHandle(c, err)
		}
		conf, err := s.app.GetCTFConfig()
		if err != nil {
			log.Printf("%+v\n", err)
			return errorHandle(c, err)
		}
		s.refreshCache(conf)
		return c.JSON(http.StatusOK, fmt.Sprintf(ChallengeUpdateTemplate, req.Name))
	}
}

func (s *server) newChallengeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(struct {
			Name        string
			Flag        string
			Category    string
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
				Category:    req.Category,
				Description: req.Description,
				Author:      req.Author,
				IsSurvey:    req.IsSurvey,
				Tags:        req.Tags,
				Attachments: req.Attachments,
				IsOpen:      chal.IsOpen,
			}); err != nil {
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}
			conf, err := s.app.GetCTFConfig()
			if err != nil {
				log.Printf("%+v\n", err)
				return xerrors.Errorf(": %w", err)
			}
			s.refreshCache(conf)
			return c.JSON(http.StatusOK, fmt.Sprintf(ChallengeUpdateTemplate, req.Name))
		} else {
			// ADD
			if err := s.app.AddChallenge(&service.Challenge{
				Name:        req.Name,
				Flag:        req.Flag,
				Category:    req.Category,
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
		submissions, err := s.app.ListValidSubmissions()

		challenges, _, err := s.app.ScoreFeed(chals, teams, submissions)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		return c.JSON(http.StatusOK, challenges)
	}
}

func (s *server) adminTeamHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		teamName := c.QueryParam("team")
		team, err := s.app.GetTeamByName(teamName)
		if err != nil && xerrors.Is(err, gorm.ErrRecordNotFound) {
			return errorHandle(c, service.NewErrorMessage(NoSuchTeamMessage))
		} else if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		// get team submissions
		submissions, err := s.app.ListTeamSubmissions(team.ID)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		res := map[string]interface{}{
			"teamname":    team.Teamname,
			"team_id":     team.ID,
			"email":       team.Email,
			"country":     team.CountryCode,
			"submissions": submissions,
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (s *server) listTeamHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		teams, err := s.app.ListAllTeams()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		return c.JSON(http.StatusOK, teams)
	}
}

func (s *server) updateTeamEmail() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(struct {
			ID    uint32 `json:"id"`
			Email string `json:"email"`
		})
		if err := c.Bind(req); err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		team, err := s.app.GetTeamByID(req.ID)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		if err := s.app.UpdateEmail(team, req.Email); err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		return messageHandle(c, ProfileUpdateMessage)
	}
}

// こんなところにロジックを書くなんてと思いつつ書く
func (s *server) recalcSeries() echo.HandlerFunc {
	return func(c echo.Context) error {
		// valid submissions を全部拾ってきて順番に適用していく
		chals, err := s.app.ListAllRawChallenges()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		teams, err := s.app.ListTeams()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		submissions, err := s.app.ListValidSubmissions()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		// submitを提出順にsort
		sort.Slice(submissions, func(i, j int) bool {
			return submissions[i].SubmittedAt < submissions[j].SubmittedAt
		})

		conf, err := s.app.GetCTFConfig()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		// 既存のseries全部消す
		if err := s.removeAllSeries(conf); err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		// series全部計算し直す（めっちゃおもそう……
		for i := 0; i < len(submissions); i++ {
			_, scoreboard, err := s.app.ScoreFeed(chals, teams, submissions[:i+1])
			if err != nil {
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}
			if err := s.appendScoreSeries(conf, scoreboard, time.Unix(submissions[i].SubmittedAt, 0)); err != nil {
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}
		}

		return messageHandle(c, "Recalc Score")
	}
}

func (s *server) allTeamSeries() echo.HandlerFunc {
	return func(c echo.Context) error {
		conf, err := s.app.GetCTFConfig()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		teams, err := s.app.ListTeams() // ここでは順位表にのるチームだけでいい
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		topTeamSeries := make(map[string]TeamScoreSeries)
		for _, t := range teams {
			series, err := s.getTeamSeries(conf, t.ID)
			if err != nil {
				if xerrors.Is(err, redis.Nil) {
					continue
				}
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}
			topTeamSeries[t.Teamname] = series
		}
		return c.JSON(http.StatusOK, topTeamSeries)
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

// metricsHandler はprometheus exporterとしてのエンドポイント
// CTFに関する集計された値を返す
// sensitiveな情報を扱うのでadmin only
func (s *server) metricsHandler() echo.HandlerFunc {
	reg := prometheus.NewRegistry()

	// 計測する物一覧（先に箱を作って、あとで実際の値を取る）
	flagCounter := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "number_of_submitted_flags",
	})
	reg.MustRegister(flagCounter)

	validFlagCounter := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "number_of_valid_flags",
	})
	reg.MustRegister(validFlagCounter)

	teamsCounter := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "number_of_registered_teams",
	})
	reg.MustRegister(teamsCounter)

	activeSessionCounter := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "number_of_active_sessions",
	})
	reg.MustRegister(activeSessionCounter)

	solveCollector := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "solve",
	}, []string{
		"name",
		"category",
	})
	reg.MustRegister(solveCollector)

	// 別にmetricsとして見たいかと言われればそうでもないので
	// scoreCollector := prometheus.NewGaugeVec(prometheus.GaugeOpts{
	// 	Name: "score",
	// }, []string{
	// 	"team",
	// 	"country",
	// })
	// reg.MustRegister(scoreCollector)

	h := promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	})

	return func(c echo.Context) error {
		numofSubmissions, err := s.app.CountSubmissions()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		flagCounter.Set(float64(numofSubmissions))

		numofValidSubmissions, err := s.app.CountValidSubmissions()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		validFlagCounter.Set(float64(numofValidSubmissions))

		numofTeams, err := s.app.CountTeams()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		teamsCounter.Set(float64(numofTeams))

		numofSessions, err := s.countActiveSessions()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		activeSessionCounter.Set(float64(numofSessions))

		solves, err := s.app.TaskSolves()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		for chal, count := range solves {
			solveCollector.With(prometheus.Labels{
				"name":     chal.Name,
				"category": chal.Category,
			}).Set(float64(count))
		}

		// conf, err := s.app.GetCTFConfig()
		// if err != nil {
		// 	return errorHandle(c, xerrors.Errorf(": %w", err))
		// }
		// scoreboard, err := s.getScoreboard(conf)
		// if err != nil {
		// 	return errorHandle(c, xerrors.Errorf(": %w", err))
		// }
		// for _, t := range scoreboard {
		// 	scoreCollector.With(prometheus.Labels{
		// 		"team":    t.Teamname,
		// 		"country": t.Country,
		// 	}).Set(float64(t.Score))
		// }

		// serve by promhttp and wrap
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

// ---

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

func (s *server) refreshCache(config *model.Config) ([]*service.Challenge, []*service.ScoreFeedEntry, error) {
	var chals []*model.Challenge
	var err error
	if service.CalcCTFStatus(config) == service.CTFNotStarted {
		chals = []*model.Challenge{}
	} else {
		chals, err = s.app.ListOpenedRawChallenges()
	}
	if err != nil {
		return nil, nil, xerrors.Errorf(": %w", err)
	}
	teams, err := s.app.ListTeams()
	if err != nil {
		return nil, nil, xerrors.Errorf(": %w", err)
	}

	submissions, err := s.app.ListValidSubmissions()
	if err != nil {
		return nil, nil, xerrors.Errorf(": %w", err)
	}

	challenges, scoreboard, err := s.app.ScoreFeed(chals, teams, submissions)
	if err != nil {
		return nil, nil, xerrors.Errorf(": %w", err)
	}
	if err := s.setChallenges(config, challenges); err != nil {
		return nil, nil, xerrors.Errorf(": %w", err)
	}
	if err := s.setScoreboard(config, scoreboard); err != nil {
		return nil, nil, xerrors.Errorf(": %w", err)
	}
	return challenges, scoreboard, nil
}

func (s *server) setChallenges(config *model.Config, challenges []*service.Challenge) error {
	challengesJson, err := json.Marshal(challenges)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	key := challengesKey(config.CTFName)

	if err := s.redis.Set(key, string(challengesJson), cacheDuration).Err(); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (s *server) getRawChallenges(config *model.Config) ([]*service.Challenge, error) {
	key := challengesKey(config.CTFName)
	challengesStr, err := s.redis.Get(key).Result()

	if err != nil {
		// nilのときrefreshする
		if xerrors.Is(err, redis.Nil) {
			challenges, _, err := s.refreshCache(config)
			if err != nil {
				return nil, xerrors.Errorf(": %w", err)
			}
			return challenges, nil
		}
		return nil, xerrors.Errorf(": %w", err)
	}

	var challenges []*service.Challenge
	if err := json.Unmarshal([]byte(challengesStr), &challenges); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return challenges, nil
}

/// 非ログインユーザ向けに色々を消したchallengesを渡す
func (s *server) getNonsensitiveChallenges(config *model.Config) ([]*service.Challenge, error) {
	challenges, err := s.getRawChallenges(config)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	for i := 0; i < len(challenges); i++ {
		challenges[i].Description = ""
		challenges[i].Tags = []string{}
		challenges[i].Flag = ""
		challenges[i].Author = ""
		challenges[i].Attachments = []service.Attachment{}
	}
	return challenges, nil
}

/// ユーザ向けにflagを消したchallengesを渡す
func (s *server) getChallenges(config *model.Config) ([]*service.Challenge, error) {
	challenges, err := s.getRawChallenges(config)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	for i := 0; i < len(challenges); i++ {
		challenges[i].Flag = ""
	}
	return challenges, nil
}

func challengesKey(ctfname string) string {
	return fmt.Sprintf("%s_challenges", ctfname)
}

func (s *server) setScoreboard(config *model.Config, scoreboard []*service.ScoreFeedEntry) error {
	scoreboardJson, err := json.Marshal(scoreboard)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	key := scoreboardKey(config.CTFName)
	if err := s.redis.Set(key, string(scoreboardJson), cacheDuration).Err(); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (s *server) getScoreboard(config *model.Config) ([]*service.ScoreFeedEntry, error) {
	key := scoreboardKey(config.CTFName)
	scoreboardStr, err := s.redis.Get(key).Result()
	if err != nil {
		// nilのときrefreshする
		if xerrors.Is(err, redis.Nil) {
			_, scoreboard, err := s.refreshCache(config)
			if err != nil {
				return nil, xerrors.Errorf(": %w", err)
			}
			return scoreboard, nil
		}
		return nil, xerrors.Errorf(": %w", err)
	}

	var scoreboard []*service.ScoreFeedEntry
	if err := json.Unmarshal([]byte(scoreboardStr), &scoreboard); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return scoreboard, nil
}

func scoreboardKey(ctfname string) string {
	return fmt.Sprintf("%s_scorefeed", ctfname)
}

type TeamScoreSeriesEntry struct {
	Teamname string `json:"teamname"`
	Score    int    `json:"score"`
	Pos      int    `json:"pos"`
	Time     int64  `json:"time"`
}

type TeamScoreSeries []*TeamScoreSeriesEntry

/// チーム毎に時系列ランキングを更新
func (s *server) appendScoreSeries(config *model.Config, standings []*service.ScoreFeedEntry, now time.Time) error {
	for _, team := range standings {
		series := TeamScoreSeriesEntry{
			Teamname: team.Teamname,
			Score:    team.Score,
			Pos:      team.Pos,
			Time:     now.Unix(),
		}
		seriesJson, err := json.Marshal(series)
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}

		key := rankingSeriesKey(config.CTFName, team.TeamID)
		err = s.redis.RPush(key, string(seriesJson)).Err()
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
	}

	return nil
}

func (s *server) getTeamSeries(config *model.Config, teamID uint32) (TeamScoreSeries, error) {
	key := rankingSeriesKey(config.CTFName, teamID)
	seriesStr, err := s.redis.LRange(key, 0, -1).Result()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	var series []*TeamScoreSeriesEntry
	if err := json.Unmarshal([]byte("["+strings.Join(seriesStr, ",")+"]"), &series); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return series, nil
}

func (s *server) removeAllSeries(config *model.Config) error {
	keys, err := s.redis.Keys(fmt.Sprintf("%s_rankingseries_*", config.CTFName)).Result()
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	if err := s.redis.Del(keys...).Err(); err != nil {
	}
	return nil
}

func rankingSeriesKey(ctfname string, team uint32) string {
	return fmt.Sprintf("%s_rankingseries_%d", ctfname, team)
}

func (s *server) countActiveSessions() (int64, error) {
	now := time.Now().Unix()
	nowStr := strconv.FormatInt(now, 10)
	// expiredなセッションを消す
	if err := s.redis.ZRemRangeByScore(sessionSetKey, "0", nowStr).Err(); err != nil {
		return 0, xerrors.Errorf(": %w", err)
	}

	cnt, err := s.redis.ZCount(sessionSetKey, nowStr, "inf").Result()
	if err != nil {
		return 0, xerrors.Errorf(": %w", err)
	}
	return cnt, nil
}
