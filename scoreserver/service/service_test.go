package service

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"github.com/theoremoon/kosenctfx/scoreserver/config"
	"github.com/theoremoon/kosenctfx/scoreserver/repository"
)

type TestApp interface {
	App
	Close()
}

type testApp struct {
	app
	db *gorm.DB
}

func (app *testApp) Close() {
	app.db.Close()
}

func newTestApp(t *testing.T) TestApp {
	t.Helper()

	conf, err := config.Load()
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open("mysql", conf.Dbdsn)
	if err != nil {
		panic(err)
	}
	db.BlockGlobalUpdate(true)

	repo := repository.New(db)
	repo.Migrate()

	return &testApp{
		app: New(repo),
		db:  db,
	}
}
