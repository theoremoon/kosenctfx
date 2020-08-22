package server

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
	"github.com/theoremoon/kosenctfx/scoreserver/webhook"
	"golang.org/x/xerrors"
)

var (
	InvalidRequestMessage         = "invalid request"
	RegisteredMessage             = "registered"
	LoginMessage                  = "logged in"
	LogoutMessage                 = "logged out"
	AlreadyAuthorizedMessage      = "you are already login"
	UnauthorizedMessage           = "login required"
	AdminUnauthorizedMessage      = "you are not the admin"
	PasswordResetEmailSentMessage = "password reset email has been sent"
	PasswordUpdateMessage         = "password is successfully reset"
	TeamnameUpdateMessage         = "temaname is successfully updated"
	RenewTeamTokenMessage         = "teamtoken is updated"

	CTFAlreadyStartedMessage = "CTF has already started"
	CTFNotStartedMessage     = "CTF has not started yet"
	CTFNotRunningMessage     = "CTF not running"

	CTFClosedMessage          = "Competition is closed now"
	RegistrationClosedMessage = "Registraction is closed now"

	ChallengeOpenMessage   = "Open the challenge"
	ChallengeAddMessage    = "Added the challenge"
	ChallengeUpdateMessage = "Updated the challenge"

	AddNotificationMessage = "Add New Notification"

	NotImplementedMessage = "Not Implemented"
)

type Server interface {
	Start(addr string) error
}
type server struct {
	app           service.App
	Token         string
	SessionKey    string
	FrontendURL   string
	adminWebhook  webhook.Webhook
	systemWebhook webhook.Webhook
}

func New(app service.App, frontendURL, token string) Server {
	return &server{
		app:           app,
		SessionKey:    "kosenctfx",
		Token:         token,
		FrontendURL:   frontendURL,
		adminWebhook:  webhook.Dummy("ADMIN"),
		systemWebhook: webhook.Dummy("SYSTEM"),
	}
}

func (s *server) Start(addr string) error {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{s.FrontendURL},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPost},
	}))

	e.POST("/register-with-team", s.registerWithTeamHandler(), s.notLoginMiddleware, s.registerableMiddleware)
	e.POST("/register-and-join-team", s.registerAndJoinTeamHandler(), s.notLoginMiddleware, s.registerableMiddleware)
	e.POST("/login", s.loginHandler())
	e.POST("/logout", s.logoutHandler())
	e.GET("/info", s.infoHandler())
	e.GET("/info-update", s.infoUpdateHandler())

	e.POST("/renew-teamtoken", s.renewTeamTokenHanler(), s.loginMiddleware)
	e.POST("/passwordreset-request", s.passwordresetRequestHandler(), s.notLoginMiddleware)
	e.POST("/passwordreset", s.passwordresetHandler(), s.notLoginMiddleware)
	e.POST("/password-update", s.passwordUpdateHandler(), s.loginMiddleware)
	e.POST("/teamname-update", s.teamnameUpdateHandler(), s.loginMiddleware, s.ctfNotStartedMiddleware)

	e.GET("/challenges", s.challengesHandler(), s.ctfStartedMiddleware, s.ctfPlayableMiddleware)
	e.GET("/ranking", s.rankingHandler(), s.ctfStartedMiddleware)
	e.GET("/notifications", s.notificationsHandler())
	e.GET("/team/:id", s.teamHandler())
	e.GET("/user/:id", s.userHandler())

	e.POST("/submit", s.submitHandler(), s.loginMiddleware, s.ctfStartedMiddleware, s.ctfPlayableMiddleware)

	e.POST("/admin/ctf-config", s.ctfConfigHandler(), s.adminMiddleware)
	e.POST("/admin/open-challenge", s.openChallengeHandler(), s.adminMiddleware)
	e.POST("/admin/update-challenge", s.updateChallengeHandler(), s.adminMiddleware)
	e.POST("/admin/new-challenge", s.newChallengeHandler(), s.adminMiddleware)
	e.POST("/admin/new-notification", s.newNotificationHandler(), s.adminMiddleware)
	e.GET("/admin/list-challenges", s.listChallengesHandler(), s.adminMiddleware)
	e.GET("/admin/list-challenges", s.listChallengesHandler(), s.adminMiddleware)
	e.POST("/admin/set-challenge-status", s.setChallengeStatusHandler(), s.adminMiddleware)

	return e.Start(addr)
}

func errorHandle(c echo.Context, err error) error {
	var errMsg service.ErrorMessage
	if xerrors.As(err, &errMsg) {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": errMsg.Error(),
		})
	}

	log.Printf("%+v\n", err)
	c.Logger().Error(err)
	return c.NoContent(http.StatusInternalServerError)
}
func messageHandle(c echo.Context, msg string) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": msg,
	})
}

func (s *server) tokenCookie(token *model.LoginToken) *http.Cookie {
	return &http.Cookie{
		Name:     s.SessionKey,
		Value:    token.Token,
		Expires:  token.ExpiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

func (s *server) removeTokenCookie() *http.Cookie {
	return &http.Cookie{
		Name:     s.SessionKey,
		Value:    "",
		Expires:  time.Time{},
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

type loginContext struct {
	echo.Context
	User *model.User
}

func (s *server) getLoginUser(c echo.Context) (*model.User, error) {
	cookie, err := c.Cookie(s.SessionKey)
	if err != nil {
		return nil, err
	}

	user, err := s.app.GetLoginUser(cookie.Value)
	if err != nil {
		return nil, err
	}
	return user, nil
}
