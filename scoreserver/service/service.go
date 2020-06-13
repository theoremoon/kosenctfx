package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/theoremoon/kosenctfx/scoreserver/repository"
)

type App interface {
	UserApp
	TeamApp
	ChallengeApp
	CTFApp
	ClarificationApp
	NotificationApp
}

type app struct {
	repo repository.Repository
}

func New(repo repository.Repository) App {
	return &app{
		repo: repo,
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
		return ErrorMessage("CTF has not started yet")
	}
	if !t.Before(conf.EndAt) {
		return ErrorMessage("CTF has alredy finished")
	}
	return nil
}
