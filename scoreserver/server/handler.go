package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/xerrors"

	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/service"

	"github.com/labstack/echo/v4"
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
		ret := make(map[string]interface{})

		// TODO notification / ranking

		// 問題情報を返す
		// TODO REDIS
		if u != nil {
			challenges, err := s.app.ListOpenChallenges()
			for i := range challenges {
				challenges[i].Flag = ""
			}
			if err != nil {
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}
			ret["challenges"] = challenges
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

func (s *server) teamnameUpdateHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		lc := c.(*loginContext)
		req := new(struct {
			NewTeamname string
		})
		if err := c.Bind(req); err != nil {
			return errorHandle(c, err)
		}
		team, err := s.app.GetTeamByID(lc.User.TeamId)
		if err != nil {
			return errorHandle(c, err)
		}

		if err := s.app.UpdateTeamname(team.ID, req.NewTeamname); err != nil {
			return errorHandle(c, err)
		}
		return messageHandle(c, TeamnameUpdateMessage)
	}
}

func (s *server) challengesHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		challenges, err := s.app.ListOpenChallenges()
		for i := range challenges {
			challenges[i].Flag = ""
		}
		if err != nil {
			return errorHandle(c, err)
		}
		return c.JSON(http.StatusOK, challenges)
	}
}

func (s *server) rankingHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusNotImplemented, NotImplementedMessage)
	}
}

func (s *server) notificationsHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		notifications, err := s.app.ListNotifications()
		if err != nil {
			return errorHandle(c, err)
		}
		return c.JSON(http.StatusOK, notifications)
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

		ctfStatus, err := s.app.CurrentCTFStatus()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		flag := strings.Trim(req.Flag, " ")
		challenge, correct, valid, err := s.app.SubmitFlag(lc.User, flag, ctfStatus == service.CTFRunning)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		team, err := s.app.GetTeamByID(lc.User.TeamId)
		if err != nil {
			log.Println(err)
			team = &model.Team{
				Teamname: "",
			}
		}

		if valid {
			s.systemWebhook.Post(fmt.Sprintf(
				"`%s@%s` solved `%s` :100:",
				lc.User.Username,
				team.Teamname,
				challenge.Name,
			))
			return messageHandle(c, fmt.Sprintf("correct! solved `%s` and got score", challenge.Name))
		} else if correct {
			s.adminWebhook.Post(fmt.Sprintf(
				"`%s@%s` solved `%s`.",
				lc.User.Username,
				team.Teamname,
				challenge.Name,
			))
			return messageHandle(c, fmt.Sprintf("correct. solved `%s`", challenge.Name))
		} else {
			s.adminWebhook.Post(fmt.Sprintf(
				"`%s@%s` submit flag `%s`, but wrong.",
				lc.User.Username,
				team.Teamname,
				req.Flag,
			))
			return errorHandle(c, service.NewErrorMessage("wrong flag"))
		}
	}
}

func (s *server) ctfConfigHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(struct {
			Name          string `json:"name"`
			StartAt       int64  `json:"start_at"`
			EndAt         int64  `json:"end_at"`
			RegisterOpen  bool   `json:"register_open"`
			CTFOpen       bool   `json:"ctf_open"`
			LockCount     int    `json:"lock_count"`
			LockFrequency int    `json:"lock_frequency"`
			LockDuration  int    `json:"lock_duration"`
			ScoreExpr     string `json:"score_expr"`
		})
		if err := c.Bind(req); err != nil {
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
		conf.LockFrequency = req.LockFrequency
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
		// TODO add notification?
		s.systemWebhook.Post(fmt.Sprintf("Challenge `%s` opened!", chal.Name))
		return c.JSON(http.StatusOK, ChallengeOpenMessage)
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

		chal, err := s.app.GetChallengeByName(req.Name)
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		return c.JSON(http.StatusOK, chal)
	}
}

func (s *server) newNotificationHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(struct {
			Content string
		})
		if err := c.Bind(req); err != nil {
			return errorHandle(c, err)
		}
		notification, err := s.app.AddNotification(req.Content)
		if err != nil {
			return errorHandle(c, err)
		}
		s.systemWebhook.Post(fmt.Sprintf("Notification: ```\n%s\n```", notification.Content))
		return c.JSON(http.StatusOK, AddNotificationMessage)
	}
}

func (s *server) listChallengesHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		challenges, err := s.app.ListAllChallenges()
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
		//TODO
		log.Printf("[+] challenge-status: %v %v\n", req.Result, req.Name)
		return c.NoContent(http.StatusOK)
	}
}
