package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/theoremoon/kosenctfx/scoreserver/model"

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
	return func(cc echo.Context) error {
		c := cc.(*loginContext)
		if err := s.app.LogoutUser(c.User.ID); err != nil {
			return errorHandle(c, err)
		}
		c.SetCookie(s.removeTokenCookie())
		return messageHandle(c, LogoutMessage)
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
			NewPassword string
		})
		if err := c.Bind(req); err != nil {
			return errorHandle(c, err)
		}
		if err := s.app.PasswordReset(req.Token, req.NewPassword); err != nil {
			return errorHandle(c, err)
		}
		return messageHandle(c, PasswordUpdateMessage)
	}
}

func (s *server) passwordUpdateHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		lc := c.(*loginContext)
		req := new(struct {
			NewPassword string
		})
		if err := c.Bind(req); err != nil {
			return errorHandle(c, err)
		}
		if err := s.app.PasswordUpdate(lc.User, req.NewPassword); err != nil {
			return errorHandle(c, err)
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
		team, err := s.app.GetUserTeam(lc.User.ID)
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
		for i, _ := range challenges {
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

func (s *server) qualificationsHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusNotImplemented, NotImplementedMessage)
	}
}

func (s *server) notificationsHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusNotImplemented, NotImplementedMessage)
	}
}

func (s *server) teamHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		teamIDstr := c.Param("id")
		teamID, err := strconv.ParseUint(teamIDstr, 10, 32)
		if err != nil {
			return errorHandle(c, err)
		}
		team, err := s.app.GetTeamByID(uint(teamID))
		if err != nil {
			return errorHandle(c, err)
		}
		return c.JSON(http.StatusOK, team)
	}
}

func (s *server) userHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userIDstr := c.Param("id")
		userID, err := strconv.ParseUint(userIDstr, 10, 32)
		if err != nil {
			return errorHandle(c, err)
		}
		user, err := s.app.GetUserByID(uint(userID))
		if err != nil {
			return errorHandle(c, err)
		}
		return c.JSON(http.StatusOK, user)
	}
}

func (s *server) doQualificationHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		lc := c.(*loginContext)
		req := new(struct {
			Content string
		})
		if err := c.Bind(req); err != nil {
			return errorHandle(c, err)
		}
		qual, err := s.app.NewQualification(lc.User, req.Content)
		if err != nil {
			return errorHandle(c, err)
		}

		// FIXME
		go s.adminWebhook.Post(fmt.Sprintf(
			"New qualification is created: %s/admin/qualifications/%d\n```\n%s```",
			c.Request().Host,
			qual.ID,
			qual.Content,
		))
		return messageHandle(c, QualificationSentMesssage)
	}
}

func (s *server) submitHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		lc := c.(*loginContext)
		req := new(struct {
			Flag string
		})
		if err := c.Bind(req); err != nil {
			return errorHandle(c, err)
		}
		challenge, correct, valid, err := s.app.SubmitFlag(lc.User, req.Flag)
		if err != nil {
			return errorHandle(c, err)
		}

		team, err := s.app.GetUserTeam(lc.User.ID)
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
			return messageHandle(c, "wrong flag")
		}
	}
}

func (s *server) initializeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusNotImplemented, NotImplementedMessage)
	}
}

func (s *server) qualificationUpdateHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusNotImplemented, NotImplementedMessage)
	}
}

func (s *server) openChallengeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusNotImplemented, NotImplementedMessage)
	}
}

func (s *server) updateChallengeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusNotImplemented, NotImplementedMessage)
	}
}

func (s *server) newChallengeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusNotImplemented, NotImplementedMessage)
	}
}

func (s *server) newNotificationHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusNotImplemented, NotImplementedMessage)
	}
}

func (s *server) adminQualificationHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusNotImplemented, NotImplementedMessage)
	}
}
