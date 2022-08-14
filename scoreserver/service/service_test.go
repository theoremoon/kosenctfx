package service_test

import (
	"testing"

	"github.com/theoremoon/kosenctfx/scoreserver/mailer"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	testDB = "testuser:testpassword@tcp(db-test:3306)/testtable"
)

func newApp(t *testing.T) service.App {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}
	model.Migrate(db)
	mailer := mailer.NewFakeMailer()
	return service.New(db, mailer)
}
