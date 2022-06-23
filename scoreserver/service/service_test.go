package service_test

import (
	"testing"

	"github.com/theoremoon/kosenctfx/scoreserver/mailer"
	"github.com/theoremoon/kosenctfx/scoreserver/repository"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	testDB = "testuser:testpassword@tcp(db-test:3306)/testtable"
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
