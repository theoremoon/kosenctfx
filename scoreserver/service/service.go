package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/theoremoon/kosenctfx/scoreserver/repository"
)

type App interface{}
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
