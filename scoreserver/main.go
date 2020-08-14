package main

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"github.com/theoremoon/kosenctfx/scoreserver/config"
	"github.com/theoremoon/kosenctfx/scoreserver/mailer"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/repository"
	"github.com/theoremoon/kosenctfx/scoreserver/server"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
)

func run() error {
	conf, err := config.Load()
	if err != nil {
		return err
	}
	db, err := gorm.Open("mysql", conf.Dbdsn)
	if err != nil {
		return err
	}
	defer db.Close()
	db.BlockGlobalUpdate(true)

	repo := repository.New(db)
	repo.Migrate()
	mailer, err := mailer.New(conf.MailServer, conf.Email, conf.MailPassword)
	if err != nil {
		return err
	}

	app := service.New(repo, mailer)

	// admin ユーザを自動生成して適当なCTF情報を入れる
	if _, err := app.GetAdminUser(); err != nil {
		// password := uuid.New().String()
		password := "password"
		if _, err := app.CreateAdminUser(conf.Email, password); err != nil {
			return err
		}

		log.Printf("---[ADMIN]---\n")
		log.Printf(" username: admin\n")
		log.Printf(" email: %s\n", conf.Email)
		log.Printf(" password: %s", password)

		app.SetCTFConfig(&model.Config{
			CTFName:       "KosenCTF X",
			StartAt:       time.Now(),
			EndAt:         time.Now(),
			RegisterOpen:  true, // FIXME: for production, this value should be false
			CTFOpen:       false,
			LockCount:     5,
			LockFrequency: 10,
			LockDuration:  1200,
			MinScore:      100,
			MaxScore:      500,
			CountToMin:    60,
		})
	}

	srv := server.New(app, conf.Front)
	return srv.Start(conf.Addr)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
