package service

import (
	"time"

	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"golang.org/x/xerrors"
)

type CTFStatus int

const (
	CTFNotStarted    CTFStatus = 2
	CTFRunning       CTFStatus = 1
	CTFEnded         CTFStatus = 4
	InvalidCTFSTatus CTFStatus = 0
)

type CTFApp interface {
	GetCTFConfig() (*model.Config, error)
	SetCTFConfig(config *model.Config) error
}

func CalcCTFStatus(conf *model.Config) CTFStatus {
	now := time.Now()

	if !conf.CTFOpen {
		return CTFNotStarted
	} else if now.Before(conf.StartAt) {
		return CTFNotStarted
	} else if now.After(conf.StartAt) && now.Before(conf.EndAt) {
		return CTFRunning
	} else {
		return CTFEnded
	}
}

func (app *app) GetCTFConfig() (*model.Config, error) {
	conf, err := app.repo.GetConfig()
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func (app *app) SetCTFConfig(config *model.Config) error {
	err := app.repo.SetConfig(config)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}
