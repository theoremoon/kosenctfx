package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/theoremoon/kosenctfx/scoreserver/mailer"
	"github.com/theoremoon/kosenctfx/scoreserver/repository"
)

type App interface {
	UserApp
	TeamApp
	ChallengeApp
	CTFApp
	NotificationApp
}

type app struct {
	repo   repository.Repository
	mailer mailer.Mailer
}

func New(repo repository.Repository, mailer mailer.Mailer) App {
	return &app{
		mailer: mailer,
		repo:   repo,
	}
}

var LoginTokenLifeSpan = 7 * 24 * time.Hour // default is 1week

func tokenExpiredTime() time.Time {
	return time.Now().Add(LoginTokenLifeSpan)
}

func newToken() string {
	return uuid.New().String()
}

func (app *app) ValidateRunning(t time.Time) error {
	conf, err := app.repo.GetConfig()
	if err != nil {
		return err
	}
	if !t.After(conf.StartAt) {
		return NewErrorMessage("CTF has not started yet")
	}
	if !t.Before(conf.EndAt) {
		return NewErrorMessage("CTF has alredy finished")
	}
	return nil
}
