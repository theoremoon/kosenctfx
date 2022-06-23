package server

import (
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/theoremoon/kosenctfx/scoreserver/mailer"
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
	testDB       = "testuser:testpassword@tcp(db-test:3306)/testtable"
	testRedis    = "testuser:testpassword@tcp(db-test:3306)/testtable"
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

func newServer(t *testing.T) *server {
	t.Helper()

	app := newApp(t)
	redis, _ := redismock.NewClientMock()

	return &server{
		app:             app,
		SessionKey:      "kosenctfx",
		Token:           TOKEN,
		FrontendURL:     FRONTEND_URL,
		redis:           redis,
		AdminWebhook:    webhook.Dummy("ADMIN"),
		TaskOpenWebhook: webhook.Dummy("TASK OPEN"),
		SolveLogWebhook: webhook.Dummy("SOLVE"),
	}
}
