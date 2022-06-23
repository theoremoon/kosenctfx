package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-redis/redismock/v8"
	"github.com/labstack/echo/v4"
	"github.com/steinfletcher/apitest"
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
