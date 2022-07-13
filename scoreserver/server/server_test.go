package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-redis/redismock/v8"
	"github.com/labstack/echo/v4"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/theoremoon/kosenctfx/scoreserver/mailer"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/repository"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
	"github.com/theoremoon/kosenctfx/scoreserver/webhook"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	TOKEN        = "token"
	FRONTEND_URL = ""
)

func newRepository(t *testing.T) repository.Repository {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}

	repo := repository.New(db)
	repo.Migrate()
	return repo
}

func newApp(t *testing.T) service.App {
	t.Helper()

	repo := newRepository(t)
	mailer := mailer.NewFakeMailer()
	return service.New(repo, mailer)
}

func newServer(t *testing.T, app service.App) *echo.Echo {
	t.Helper()

	redis, _ := redismock.NewClientMock()

	admin, err := app.RegisterTeam("admin", "admin", "admin@example.com", "")
	if err != nil {
		panic(err)
	}
	if err := app.MakeTeamAdmin(admin); err != nil {
		panic(err)
	}

	err = app.SetCTFConfig(&model.Config{
		CTFName:      "KosenCTF X",
		Token:        "admin@example.com",
		StartAt:      time.Now().Unix(),
		EndAt:        time.Now().Add(1 * time.Hour).Unix(),
		RegisterOpen: true,
		CTFOpen:      true,
		LockCount:    5,
		LockDuration: 60,
		LockSecond:   300,
		ScoreExpr:    "func calc(count) { return count; }",
	})
	if err != nil {
		panic(err)
	}

	srv := &server{
		app:             app,
		SessionKey:      "kosenctfx",
		Token:           TOKEN,
		FrontendURL:     FRONTEND_URL,
		redis:           redis,
		AdminWebhook:    webhook.Dummy("ADMIN"),
		TaskOpenWebhook: webhook.Dummy("TASK OPEN"),
		SolveLogWebhook: webhook.Dummy("SOLVE"),
	}
	return srv.build(true)
}

func TestRegister(t *testing.T) {
	t.Parallel()
	app := newApp(t)
	s := newServer(t, app)

	tests := []struct {
		Teamname string
		Password string
		Email    string
		Country  string
		message  string
		expected int
	}{
		{
			Teamname: "team1",
			Password: "password",
			Email:    "team1@example.com",
			Country:  "jp",
			message:  "registerできる",
			expected: http.StatusOK,
		},
		{
			Teamname: "team1",
			Password: gofakeit.UUID(),
			Email:    gofakeit.Email(),
			Country:  "",
			message:  "teamname unique",
			expected: http.StatusBadRequest,
		},
		{
			Teamname: gofakeit.Name(),
			Password: gofakeit.UUID(),
			Email:    "team1@example.com",
			Country:  "",
			message:  "email unique",
			expected: http.StatusBadRequest,
		},
		{
			Teamname: "",
			Password: gofakeit.UUID(),
			Email:    gofakeit.Email(),
			Country:  "",
			message:  "teamnameいる",
			expected: http.StatusBadRequest,
		},
		{
			Teamname: gofakeit.Name(),
			Password: "",
			Email:    gofakeit.Email(),
			Country:  "",
			message:  "passwordいる",
			expected: http.StatusBadRequest,
		},
		{
			Teamname: gofakeit.Name(),
			Password: gofakeit.UUID(),
			Email:    "",
			Country:  "",
			message:  "emailいる",
			expected: http.StatusBadRequest,
		},
		{
			Teamname: gofakeit.Name(),
			Password: gofakeit.UUID(),
			Email:    gofakeit.Email(),
			Country:  gofakeit.UUID(),
			message:  "countryはcountry",
			expected: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.message, func(t *testing.T) {
			apitest.New().Handler(s).
				Post("/register").
				JSON(tt).
				Expect(t).Status(tt.expected).End()
		})
	}
}

func createTeam(t *testing.T, app service.App, s *echo.Echo) (*model.Team, []*apitest.Cookie) {
	t.Helper()

	// register
	user := struct {
		Teamname string
		Password string
		Email    string
	}{
		Teamname: gofakeit.Name(),
		Password: gofakeit.UUID(),
		Email:    gofakeit.Email(),
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	registerReq := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(userJSON))
	registerReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	registerRec := httptest.NewRecorder()
	s.ServeHTTP(registerRec, registerReq)

	// login
	loginReq := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(userJSON))
	loginReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	loginRec := httptest.NewRecorder()
	s.ServeHTTP(loginRec, loginReq)

	cookies := make([]*apitest.Cookie, 0)
	for _, c := range loginRec.Result().Cookies() {
		cookies = append(cookies, apitest.FromHTTPCookie(c))
	}

	// get uesr
	team, err := app.GetTeamByName(user.Teamname)
	if err != nil {
		panic(err)
	}

	return team, cookies
}

func TestLogin(t *testing.T) {
	t.Parallel()
	app := newApp(t)
	s := newServer(t, app)

	// register
	user := struct {
		Teamname string
		Password string
		Email    string
	}{
		Teamname: gofakeit.Name(),
		Password: gofakeit.UUID(),
		Email:    gofakeit.Email(),
	}
	apitest.New().Handler(s).
		Post("/register").
		JSON(user).
		Expect(t).Status(http.StatusOK).End()

	// login
	tests := []struct {
		Teamname string
		Password string
		message  string
		expected int
	}{
		{
			Teamname: user.Teamname,
			Password: user.Password,
			message:  "loginできる",
			expected: http.StatusOK,
		},
		{
			Teamname: strings.ToUpper(user.Teamname),
			Password: user.Password,
			message:  "username case sensitive",
			expected: http.StatusBadRequest,
		},
		{
			Teamname: user.Teamname,
			Password: strings.ToUpper(user.Password),
			message:  "password case sensitive",
			expected: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.message, func(t *testing.T) {
			apitest.New().Handler(s).
				Post("/login").
				JSON(tt).
				Expect(t).Status(tt.expected).End()
		})
	}
}

func TestAccount(t *testing.T) {
	t.Parallel()
	app := newApp(t)
	s := newServer(t, app)

	t.Run("アカウント取得できる", func(t *testing.T) {
		t.Parallel()

		team, cookies := createTeam(t, app, s)
		apitest.New().Handler(s).
			Get("/account").Cookies(cookies...).
			Expect(t).Assert(
			jsonpath.Chain().
				Equal("$.teamname", team.Teamname).
				Equal("$.country", team.CountryCode).
				Equal("$.is_admin", team.IsAdmin).
				End(),
		).Status(http.StatusOK).End()
	})

	t.Run("ログインしてないときも一旦200", func(t *testing.T) {
		t.Parallel()

		apitest.New().Handler(s).
			Get("/account").
			Expect(t).
			Body("").
			Status(http.StatusOK).End()
	})
}
