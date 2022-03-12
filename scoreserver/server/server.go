package server

import (
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/theoremoon/kosenctfx/scoreserver/bucket"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
	"github.com/theoremoon/kosenctfx/scoreserver/webhook"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

var (
	AdminUnauthorizedMessage            = "You are not admin"
	AlreadyAuthorizedMessage            = "You are already logged in"
	BucketNullMessage                   = "Bucket information is not registered to the server"
	CTFAlreadyStartedMessage            = "CTF has already started"
	CTFClosedMessage                    = "Competition is closed now"
	CTFNotRunningMessage                = "CTF is not running now"
	CTFNotStartedMessage                = "CTF has not started yet"
	ChallengeAddTemplate                = "Add challenge: `%s`"
	ChallengeAlreadyClosedTemplate      = "`%s` is already closed"
	ChallengeAlreadyOpenedTemplate      = "`%s` is already opened"
	ChallengeCloseTemplate              = "`%s` is closed"
	ChallengeClosedAdminMessage         = "Challenge `%s` closed"
	ChallengeOpenAdminMessage           = "Challenge `%s` opened!"
	ChallengeOpenSystemMessage          = "Challenge `%s` opened!"
	ChallengeOpenTemplate               = "`%s` is opened"
	ChallengeUpdateTemplate             = "Updated the challenge: `%s`"
	ConfigUpdateMessage                 = "Config is updated"
	CorrectSubmissionAdminMessage       = "`%s` solved `%s`: `%s`"
	CorrectSubmissionMessage            = "Correct! You solved `%s`"
	InvalidRequestMessage               = "Invalid request"
	LoginMessage                        = "Logged in"
	LogoutMessage                       = "Logged out"
	NotImplementedMessage               = "Not Implemented"
	PasswordResetEmailSentMessage       = "We've sent you the password reset token"
	PasswordUpdateMessage               = "Password is successfully reset"
	PresignedURLKeyRequiredMessage      = "Key is required"
	ProfileUpdateMessage                = "Team profile is successfully updated"
	RegisteredMessage                   = "Registered!"
	RegistrationClosedMessage           = "Registration is closed now"
	ScoreEmulateMaxCountTooSmallMessage = "maxCount should be larger than 0"
	SolvabilityCheckedSolveMessage      = ":heavy_check_mark: `%s`"
	SolvabilityFailedSystemMessage      = ":warning: `%s`"
	SubmissionLockedMessage             = "Your submission is currently locked. Please wait for minutes."
	UnauthorizedMessage                 = "Login is required"
	ValidSubmissionAdminMessage         = "`%s` solved `%s` :100:, `%s`"
	ValidSubmissionMessage              = "Correct! You solved `%s`"
	ValidSubmissionSystemMessage        = "`%s` solved `%s` :100:"
	WrongSubmissionAdminMessage         = "`%s` submits a wrong flag: `%s`"
	WrongSubmissionMessage              = "Wrong flag..."
	NoSuchTeamMessage                   = "No such team"
)

const (
	GRAPHQL_COMPLEXITY_LIMIT = 200
)

type Server interface {
	Start(addr string) error
}
type server struct {
	app             service.App
	db              *gorm.DB
	Token           string
	SessionKey      string
	FrontendURL     string
	redis           *redis.Client
	AdminWebhook    webhook.Webhook
	SolveLogWebhook webhook.Webhook
	TaskOpenWebhook webhook.Webhook
	Bucket          bucket.Bucket
}

func New(app service.App, db *gorm.DB, redis *redis.Client, frontendURL, token string) *server {
	return &server{
		app:             app,
		db:              db,
		SessionKey:      "kosenctfx",
		Token:           token,
		FrontendURL:     frontendURL,
		redis:           redis,
		AdminWebhook:    webhook.Dummy("ADMIN"),
		TaskOpenWebhook: webhook.Dummy("TASK OPEN"),
		SolveLogWebhook: webhook.Dummy("SOLVE"),
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

	e.POST("/register", s.registerHandler(), s.notLoginMiddleware, s.registerableMiddleware)
	e.POST("/login", s.loginHandler())
	e.POST("/logout", s.logoutHandler())
	e.GET("/ctf", s.ctfHandler())
	e.GET("/account", s.accountHandler())
	e.GET("/scoreboard", s.scoreboardHandler())
	e.GET("/tasks", s.tasksHandler())
	e.POST("/series", s.seriesHandler())

	e.POST("/passwordreset-request", s.passwordresetRequestHandler(), s.notLoginMiddleware)
	e.POST("/passwordreset", s.passwordresetHandler(), s.notLoginMiddleware)
	e.POST("/update-profile", s.profileUpdateHandler(), s.loginMiddleware)

	e.GET("/team/:id", s.teamHandler())

	e.POST("/submit", s.submitHandler(), s.loginMiddleware, s.ctfStartedMiddleware)

	e.GET("/admin/score-emulate", s.scoreEmulateHandler(), s.adminMiddleware)
	e.GET("/admin/get-config", s.getConfigHandler(), s.adminMiddleware)
	e.POST("/admin/set-config", s.ctfConfigHandler(), s.adminMiddleware)
	e.POST("/admin/open-challenge", s.openChallengeHandler(), s.adminMiddleware)
	e.POST("/admin/close-challenge", s.closeChallengeHandler(), s.adminMiddleware)
	e.POST("/admin/update-challenge", s.updateChallengeHandler(), s.adminMiddleware)
	e.POST("/admin/new-challenge", s.newChallengeHandler(), s.adminMiddleware)
	e.GET("/admin/list-challenges", s.listChallengesHandler(), s.adminMiddleware)
	e.GET("/admin/tasks.md", s.tasksMDHandler(), s.adminMiddleware)
	e.GET("/admin/team", s.adminTeamHandler(), s.adminMiddleware)
	e.GET("/admin/teams", s.listTeamHandler(), s.adminMiddleware)
	e.POST("/admin/update-email", s.updateTeamEmail(), s.adminMiddleware)
	e.POST("/admin/recalc-series", s.recalcSeries(), s.adminMiddleware)
	e.GET("/admin/all-team-series", s.allTeamSeries(), s.adminMiddleware)
	e.POST("/admin/get-presigned-url", s.getPresignedURLHandler(), s.adminMiddleware)
	e.POST("/admin/sql", s.sqlHandler(), s.adminMiddleware)

	// prometheus exporter
	e.GET("/admin/metrics", s.metricsHandler(), s.adminMiddleware)

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

func errorMessageHandle(c echo.Context, status int, msg string) error {
	return c.JSON(status, map[string]interface{}{
		"message": msg,
	})
}
func messageHandle(c echo.Context, msg string) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": msg,
	})
}

func messageHandleWithStatus(c echo.Context, status int, msg string) error {
	return c.JSON(status, map[string]interface{}{
		"message": msg,
	})
}

func (s *server) tokenCookie(token *model.LoginToken) *http.Cookie {
	return &http.Cookie{
		Name:     s.SessionKey,
		Value:    token.Token,
		Expires:  time.Unix(token.ExpiresAt, 0),
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
	Team *model.Team
}

func (s *server) getLoginTeam(c echo.Context) (*model.Team, error) {
	cookie, err := c.Cookie(s.SessionKey)
	if err != nil {
		return nil, err
	}

	team, err := s.app.GetLoginTeam(cookie.Value)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (s *server) getLoginToken(c echo.Context) (string, error) {
	cookie, err := c.Cookie(s.SessionKey)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
