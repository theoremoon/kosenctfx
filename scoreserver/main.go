package main

import (
	"log"
	"time"

	"github.com/go-redis/redis"
	"golang.org/x/xerrors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/theoremoon/kosenctfx/scoreserver/config"
	"github.com/theoremoon/kosenctfx/scoreserver/mailer"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/repository"
	"github.com/theoremoon/kosenctfx/scoreserver/server"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
	"github.com/theoremoon/kosenctfx/scoreserver/webhook"
)

const scoreFunc = `
func calc(count) {
    if (count == 0) {
        return 500
    }
    a = 2.079
    b = -1.5
    c = math.Log10(toFloat(count))
    return toInt(toFloat(500) * math.Pow(1.0 + c * c / a, b))
}
`

func run() error {
	conf, err := config.Load()
	if err != nil {
		return err
	}
	db, err := gorm.Open(mysql.Open(conf.Dbdsn), &gorm.Config{})
	if err != nil {
		return err
	}
	rawdb, err := db.DB()
	if err != nil {
		return err
	}
	defer rawdb.Close()
	rawdb.SetMaxOpenConns(100)

	opt, err := redis.ParseURL(conf.RedisAddr)
	if err != nil {
		return err
	}
	redis := redis.NewClient(opt)

	repo := repository.New(db)
	repo.Migrate()

	var mailSender mailer.Mailer
	if conf.MailFake {
		mailSender = mailer.NewFakeMailer()
	} else {
		mailSender, err = mailer.New(conf.MailServer, conf.Email, conf.MailPassword)
		if err != nil {
			return err
		}
	}

	app := service.New(repo, mailSender)

	// admin ユーザを自動生成して適当なCTF情報を入れる
	if _, err := app.GetAdminTeam(); err != nil {
		t, err := app.RegisterTeam("admin", conf.AdminToken, conf.Email, "")
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
		if err := app.MakeTeamAdmin(t); err != nil {
			return xerrors.Errorf(": %w", err)
		}

		err = app.SetCTFConfig(&model.Config{
			CTFName:      "KosenCTF X",
			Token:        conf.AdminToken,
			StartAt:      time.Now().Unix(),
			EndAt:        time.Now().Unix(),
			RegisterOpen: false,
			CTFOpen:      false,
			LockCount:    5,
			LockDuration: 60,
			LockSecond:   300,
			ScoreExpr:    scoreFunc,
		})
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
	}

	ctfConf, err := app.GetCTFConfig()
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	srv := server.New(app, db, redis, conf.Front, ctfConf.Token)
	if conf.AdminWebhookURL != "" {
		srv.AdminWebhook = webhook.NewDiscord(conf.AdminWebhookURL, 1*time.Second)
	}
	if conf.TaskOpenWebhookURL != "" {
		srv.TaskOpenWebhook = webhook.NewDiscord(conf.TaskOpenWebhookURL, 1*time.Second)
	}
	if conf.SolveLogWebhookURL != "" {
		srv.SolveLogWebhook = webhook.NewDiscord(conf.SolveLogWebhookURL, 1*time.Second)
	}

	return srv.Start(conf.Addr)
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v\n", err)
	}
}
