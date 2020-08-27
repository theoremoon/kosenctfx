package main

import (
	"log"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/xerrors"

	"github.com/theoremoon/kosenctfx/scoreserver/config"
	"github.com/theoremoon/kosenctfx/scoreserver/mailer"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/repository"
	"github.com/theoremoon/kosenctfx/scoreserver/server"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
	"github.com/theoremoon/kosenctfx/scoreserver/webhook"
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

	redis := redis.NewClient(&redis.Options{
		Addr: conf.RedisAddr,
	})

	repo := repository.New(db)
	repo.Migrate()
	mailer, err := mailer.New(conf.MailServer, conf.Email, conf.MailPassword)
	if err != nil {
		return err
	}

	app := service.New(repo, mailer)

	// admin ユーザを自動生成して適当なCTF情報を入れる
	if _, err := app.GetAdminUser(); err != nil {
		password := uuid.New().String()
		if _, err := app.CreateAdminUser(conf.Email, password); err != nil {
			return err
		}
		token := password
		log.Printf("---[ADMIN]---\n")
		log.Printf(" username: admin\n")
		log.Printf(" email: %s\n", conf.Email)
		log.Printf(" password: %s", password)
		log.Printf(" token: %s", token)

		err = app.SetCTFConfig(&model.Config{
			CTFName:      "KosenCTF X",
			Token:        token,
			StartAt:      time.Now(),
			EndAt:        time.Now(),
			RegisterOpen: false,
			CTFOpen:      false,
			LockCount:    5,
			LockDuration: 60,
			LockSecond:   300,
		})
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
	}

	ctfConf, err := app.GetCTFConfig()
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	srv := server.New(app, redis, conf.Front, ctfConf.Token)
	if conf.AdminWebhookURL != "" {
		srv.AdminWebhook = webhook.NewDiscord(conf.AdminWebhookURL)
	}
	if conf.SystemWebhookURL != "" {
		srv.SystemWebhook = webhook.NewDiscord(conf.SystemWebhookURL)
	}

	return srv.Start(conf.Addr)
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v\n", err)
	}
}
