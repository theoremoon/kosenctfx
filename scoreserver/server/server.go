package server

import (
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	graphqlHandler "github.com/99designs/gqlgen/handler"
	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/theoremoon/kosenctfx/scoreserver/bucket"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/resolver"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
	"github.com/theoremoon/kosenctfx/scoreserver/webhook"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
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
	ProfileUpdateMessage          = "team profile is successfully updated"

	CTFAlreadyStartedMessage = "CTF has already started"
	CTFNotStartedMessage     = "CTF has not started yet"
	CTFNotRunningMessage     = "CTF not running"

	CTFClosedMessage          = "Competition is closed now"
	RegistrationClosedMessage = "Registraction is closed now"

	ChallengeOpenTemplate          = "`%s` is opened"
	ChallengeAlreadyOpenedTemplate = "`%s` is already opened"
	ChallengeCloseTemplate         = "`%s` is closed"
	ChallengeAlreadyClosedTemplate = "`%s` is already closed"
	ChallengeAddTemplate           = "Add challenge: `%s`"
	ChallengeUpdateTemplate        = "Updated the challenge: `%s`"

	ConfigUpdateMessage = "Config is Updated"

	NotImplementedMessage = "Not Implemented"
)

const (
	GRAPHQL_COMPLEXITY_LIMIT = 200
)

type Server interface {
	Start(addr string) error
}
type server struct {
	app           service.App
	db            *gorm.DB
	Token         string
	SessionKey    string
	FrontendURL   string
	redis         *redis.Client
	AdminWebhook  webhook.Webhook
	SolveWebhook  webhook.Webhook
	SystemWebhook webhook.Webhook
	Bucket        bucket.Bucket
}

func New(app service.App, db *gorm.DB, redis *redis.Client, frontendURL, token string) *server {
	return &server{
		app:           app,
		db:            db,
		SessionKey:    "kosenctfx",
		Token:         token,
		FrontendURL:   frontendURL,
		redis:         redis,
		AdminWebhook:  webhook.Dummy("ADMIN"),
		SystemWebhook: webhook.Dummy("SYSTEM"),
		SolveWebhook:  webhook.Dummy("SOLVE"),
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
	e.GET("/info", s.infoHandler())
	e.GET("/info-update", s.infoUpdateHandler())

	e.POST("/passwordreset-request", s.passwordresetRequestHandler(), s.notLoginMiddleware)
	e.POST("/passwordreset", s.passwordresetHandler(), s.notLoginMiddleware)
	e.POST("/update-profile", s.profileUpdateHandler(), s.loginMiddleware)

	e.GET("/team/:id", s.teamHandler())

	e.POST("/submit", s.submitHandler(), s.loginMiddleware, s.ctfStartedMiddleware, s.ctfPlayableMiddleware)

	e.GET("/admin/score-emulate", s.scoreEmulateHandler(), s.adminMiddleware)
	e.GET("/admin/get-config", s.getConfigHandler(), s.adminMiddleware)
	e.POST("/admin/set-config", s.ctfConfigHandler(), s.adminMiddleware)
	e.POST("/admin/open-challenge", s.openChallengeHandler(), s.adminMiddleware)
	e.POST("/admin/close-challenge", s.closeChallengeHandler(), s.adminMiddleware)
	e.POST("/admin/update-challenge", s.updateChallengeHandler(), s.adminMiddleware)
	e.POST("/admin/new-challenge", s.newChallengeHandler(), s.adminMiddleware)
	// e.POST("/admin/new-notification", s.newNotificationHandler(), s.adminMiddleware)
	e.GET("/admin/list-challenges", s.listChallengesHandler(), s.adminMiddleware)
	e.POST("/admin/set-challenge-status", s.setChallengeStatusHandler(), s.adminMiddleware)
	e.POST("/admin/get-presigned-url", s.getPresignedURLHandler(), s.adminMiddleware)

	// GraphQL
	// e.POST("/query", s.graphQLHandler(), s.resolveLoginMiddleware, s.attachLoaderMiddleware)
	e.GET("/playground", s.playgroundHandler())

	e.POST("/admin/sql", s.sqlHandler(), s.adminMiddleware)

	return e.Start(addr)
}

func (s *server) graphQLHandler() echo.HandlerFunc {
	h := handler.New(resolver.NewExecutableSchema(
		resolver.Config{
			Resolvers: resolver.NewResolver(s.app),
		},
	))
	h.AddTransport(transport.POST{})
	h.Use(extension.FixedComplexityLimit(GRAPHQL_COMPLEXITY_LIMIT))
	h.Use(extension.Introspection{})
	return echo.WrapHandler(h)
}

func (s *server) playgroundHandler() echo.HandlerFunc {
	h := graphqlHandler.Playground("PlayGround", "/query")
	return echo.WrapHandler(h)
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
