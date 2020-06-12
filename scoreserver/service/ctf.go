package service

import "time"

type CTFStatus int

const (
	CTFNotStarted    CTFStatus = 2
	CTFRunning       CTFStatus = 1
	CTFEnded         CTFStatus = 4
	InvalidCTFSTatus CTFStatus = 0
)

type CTFApp interface {
	CurrentCTFStatus() (CTFStatus, error)
}

func (app *app) CurrentCTFStatus() (CTFStatus, error) {
	now := time.Now()
	conf, err := app.repo.GetConfig()
	if err != nil {
		return InvalidCTFSTatus, err
	}

	if now.Before(conf.StartAt) {
		return CTFNotStarted, nil
	} else if now.After(conf.StartAt) && now.Before(conf.EndAt) {
		return CTFRunning, nil
	} else {
		return CTFEnded, nil
	}
}
