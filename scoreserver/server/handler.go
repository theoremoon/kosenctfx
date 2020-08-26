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
	"github.com/labstack/echo/v4"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
)

const (
	cacheDuration      = 1 * time.Minute
	challengesJSONKey  = "challengesJSONKey"
	rankingJSONKey     = "rankingJSONKey"
	userRankingJSONKey = "userRankingJSONKey"
)

func (s *server) registerWithTeamHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(struct {
			Teamname string
			Username string
			Email    string
			Password string
		})
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": InvalidRequestMessage,
			})
		}

		if err := s.app.RegisterUserWithTeam(req.Username, req.Password, req.Email, req.Teamname); err != nil {
			return errorHandle(c, err)
		}
		return messageHandle(c, RegisteredMessage)
	}
}

func (s *server) registerAndJoinTeamHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(struct {
			Token    string
			Username string
			Email    string
			Password string
		})
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": InvalidRequestMessage,
			})
		}

		if err := s.app.RegisterUserAndJoinToTeam(req.Username, req.Password, req.Email, req.Token); err != nil {
			return errorHandle(c, err)
		}
		return messageHandle(c, RegisteredMessage)
	}
}

func (s *server) loginHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(struct {
			Username string
			Password string
		})
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": InvalidRequestMessage,
			})
		}

		token, err := s.app.LoginUser(req.Username, req.Password)
		if err != nil {
			return errorHandle(c, err)
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
		user, _ := s.getLoginUser(c)
		if user != nil {
			team, err := s.app.GetTeamByID(user.TeamId)
			if err != nil {
				return errorHandle(c, err)
			}
			ret["username"] = user.Username
			ret["teamname"] = team.Teamname
			ret["userid"] = user.ID
			ret["teamid"] = team.ID
		}
		conf, err := s.app.GetCTFConfig()
		if err != nil {
			return errorHandle(c, err)
		}
		ret["ctf_start"] = conf.StartAt.Unix()
		ret["ctf_end"] = conf.EndAt.Unix()
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
		u, _ := s.getLoginUser(c)
		refresh := c.QueryParam("refresh")

		conf, err := s.app.GetCTFConfig()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		status := service.CalcCTFStatus(conf)
		if status == service.CTFNotStarted {
			return messageHandle(c, CTFNotStartedMessage)
		}

		ret := make(map[string]interface{})

		// TODO notification
		// cache を使う
		if refresh == "" && status == service.CTFRunning {
			challenges, ranking, userRanking, err := s.getCacheInfo()
			if err != nil {
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}

			if challenges != "" && ranking != "" && userRanking != "" {
				var cs []*service.Challenge
				var scoreboard *service.Scoreboard
				var userScoreboard *service.Scoreboard
				err1 := json.Unmarshal([]byte(challenges), &cs)
				err2 := json.Unmarshal([]byte(ranking), &scoreboard)
				err3 := json.Unmarshal([]byte(userRanking), &userScoreboard)
				if err1 == nil && err2 == nil && err3 == nil {
					ret["challenges"] = cs
					ret["ranking"] = scoreboard
					ret["userRanking"] = userScoreboard
				}
			}
		}

		_, exist1 := ret["challenges"]
		_, exist2 := ret["ranking"]
		_, exist3 := ret["userRanking"]
		if !exist1 || !exist2 || !exist3 {
			challenges, ranking, userRanking, err := s.app.ScoreFeed()
			if err != nil {
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}
			for i := range challenges {
				challenges[i].Flag = ""
			}
			ret["challenges"] = challenges
			ret["ranking"] = ranking
			ret["userRanking"] = userRanking

			// cacheする
			if status == service.CTFRunning {
				bytes1, err1 := json.Marshal(challenges)
				bytes2, err2 := json.Marshal(ranking)
				bytes3, err3 := json.Marshal(userRanking)
				if err1 == nil && err2 == nil && err3 == nil {
					s.setCacheInfo(string(bytes1), string(bytes2), string(bytes3))
				}
			}
		}

		if u == nil {
			delete(ret, "challenges")
		}
		return c.JSON(http.StatusOK, ret)
	}
}

func (s *server) renewTeamTokenHanler() echo.HandlerFunc {
	return func(c echo.Context) error {
		lc := c.(*loginContext)
		t, err := s.app.GetTeamByID(lc.User.TeamId)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		if err := s.app.RenewTeamToken(t); err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		return messageHandle(c, RenewTeamTokenMessage)
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

func (s *server) passwordUpdateHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		lc := c.(*loginContext)
		req := new(struct {
			NewPassword string `json:"new_password"`
		})
		if err := c.Bind(req); err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		if err := s.app.PasswordUpdate(lc.User, req.NewPassword); err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		return messageHandle(c, PasswordUpdateMessage)
	}
}

func (s *server) teamHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		teamIDstr := c.Param("id")
		teamID, err := strconv.ParseUint(teamIDstr, 10, 32)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		team, err := s.app.GetTeamByID(uint(teamID))
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		members, err := s.app.GetTeamMembers(team.ID)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		res := map[string]interface{}{
			"teamname": team.Teamname,
			"teamid":   team.ID,
			"members":  s.app.UsersFilter(members),
		}

		// TODO: チームで解いた問題を追加
		// ログインしていて自チームの場合Tokenもつく
		user, _ := s.getLoginUser(c)
		if user != nil && user.TeamId == uint(teamID) {
			res["token"] = team.Token
		}
		return c.JSON(http.StatusOK, res)
	}
}

func (s *server) userHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userIDstr := c.Param("id")
		userID, err := strconv.ParseUint(userIDstr, 10, 32)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		user, err := s.app.GetUserByID(uint(userID))
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		team, err := s.app.GetTeamByID(user.TeamId)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		// TODO: ユーザで解いた問題を追加
		return c.JSON(http.StatusOK, map[string]interface{}{
			"username": user.Username,
			"teamname": team.Teamname,
			"userid":   user.ID,
			"teamid":   team.ID,
		})
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
		submittable, err := s.app.CheckSubmittable(lc.User.TeamId)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		if !submittable {
			return errorHandle(c, xerrors.Errorf(": %w", service.NewErrorMessage("Your submission is currently locked")))
		}

		conf, err := s.app.GetCTFConfig()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		ctfStatus := service.CalcCTFStatus(conf)

		// flag submission
		flag := strings.Trim(req.Flag, " ")
		challenge, correct, valid, err := s.app.SubmitFlag(lc.User, flag, ctfStatus == service.CTFRunning)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		team, err := s.app.GetTeamByID(lc.User.TeamId)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		if valid {
			s.SystemWebhook.Post(fmt.Sprintf(
				"`%s@%s` solved `%s` :100:, `%s`",
				lc.User.Username,
				team.Teamname,
				challenge.Name,
				req.Flag,
			))
			return messageHandle(c, fmt.Sprintf("correct! solved `%s` and got score", challenge.Name))
		} else if correct {
			s.AdminWebhook.Post(fmt.Sprintf(
				"`%s@%s` solved `%s`: `%s`",
				lc.User.Username,
				team.Teamname,
				challenge.Name,
				req.Flag,
			))
			return messageHandle(c, fmt.Sprintf("correct. solved `%s`", challenge.Name))
		} else {
			// wrong count
			count, err := s.app.GetWrongCount(team.ID, time.Duration(conf.LockDuration)*time.Second)
			if err != nil {
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}
			if count >= conf.LockCount {
				if err := s.app.LockSubmission(team.ID, time.Duration(conf.LockSecond)*time.Second); err != nil {
					return errorHandle(c, xerrors.Errorf(": %w", err))
				}
			}

			s.AdminWebhook.Post(fmt.Sprintf(
				"`%s@%s` submit flag `%s`, but wrong.",
				lc.User.Username,
				team.Teamname,
				req.Flag,
			))
			return errorHandle(c, service.NewErrorMessage("wrong flag"))
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
			return errorHandle(c, xerrors.Errorf(": %w", service.NewErrorMessage("maxCount should be larger than 0")))
		}

		conf, err := s.app.GetCTFConfig()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		scores := make([]int, maxCount+1)
		for i := 0; i <= maxCount; i++ {
			scores[i], err = service.CalcChallengeScore(i, conf.ScoreExpr)
			if err != nil {
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}
		}
		return c.JSON(http.StatusOK, scores)
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
			LockCount    int    `json:"lock_count"`
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
		conf.StartAt = time.Unix(req.StartAt, 0)
		conf.EndAt = time.Unix(req.EndAt, 0)
		conf.RegisterOpen = req.RegisterOpen
		conf.CTFOpen = req.CTFOpen
		conf.LockCount = req.LockCount
		conf.LockSecond = req.LockSecond
		conf.LockDuration = req.LockDuration
		conf.ScoreExpr = req.ScoreExpr
		if err := s.app.SetCTFConfig(conf); err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		return c.NoContent(http.StatusOK)
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
		if err := s.app.OpenChallenge(chal.ID); err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		s.SystemWebhook.Post(fmt.Sprintf("Challenge `%s` opened!", chal.Name))
		return c.JSON(http.StatusOK, ChallengeOpenMessage)
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
		if err := s.app.CloseChallenge(chal.ID); err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		s.SystemWebhook.Post(fmt.Sprintf("Challenge `%s` closed!", chal.Name))
		return c.JSON(http.StatusOK, ChallengeCloseMessage)
	}
}

func (s *server) updateChallengeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(struct {
			ID          uint
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
		return c.JSON(http.StatusOK, ChallengeUpdateMessage)
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
			}); err != nil {
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}
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
		}

		chal, err := s.app.GetRawChallengeByName(req.Name)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		return c.JSON(http.StatusOK, chal)
	}
}

func (s *server) listChallengesHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		challenges, err := s.app.ListAllRawChallenges()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		return c.JSON(http.StatusOK, challenges)
	}
}

func (s *server) setChallengeStatusHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(struct {
			Name   string `json:"name"`
			Result bool   `json:"result"`
		})
		if err := c.Bind(req); err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		if req.Result {
			s.AdminWebhook.Post(fmt.Sprintf(":heavy_check_mark: Solvability Checked: `%s`", req.Name))
		} else {
			s.SystemWebhook.Post(fmt.Sprintf(":heavy_multiplication_x: failed to solve: `%s`", req.Name))
		}
		return c.NoContent(http.StatusOK)
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

func (s *server) getCacheInfo() (string, string, string, error) {
	c, err := s.redis.Get(challengesJSONKey).Result()
	if err == redis.Nil {
		return "", "", "", nil
	} else if err != nil {
		return "", "", "", xerrors.Errorf(": %w", err)
	}

	r, err := s.redis.Get(rankingJSONKey).Result()
	if err == redis.Nil {
		return "", "", "", nil
	} else if err != nil {
		return "", "", "", xerrors.Errorf(": %w", err)
	}

	u, err := s.redis.Get(userRankingJSONKey).Result()
	if err == redis.Nil {
		return "", "", "", nil
	} else if err != nil {
		return "", "", "", xerrors.Errorf(": %w", err)
	}

	return c, r, u, nil
}

func (s *server) setCacheInfo(challenges, ranking, userRanking string) error {
	err := s.redis.Set(challengesJSONKey, challenges, cacheDuration).Err()
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	err = s.redis.Set(rankingJSONKey, ranking, cacheDuration).Err()
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	err = s.redis.Set(userRankingJSONKey, userRanking, cacheDuration).Err()
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}
