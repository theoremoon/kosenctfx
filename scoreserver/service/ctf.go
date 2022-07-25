package service

import (
	"time"

	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
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
	now := time.Now().Unix()

	if !conf.CTFOpen {
		return CTFNotStarted
	} else if now < conf.StartAt {
		return CTFNotStarted
	} else if conf.StartAt <= now && now < conf.EndAt {
		return CTFRunning
	} else {
		return CTFEnded
	}
}

func (app *app) SetCTFConfig(conf *model.Config) error {
	err := app.db.Transaction(func(tx *gorm.DB) error {
		var c model.Config
		if err := tx.First(&c).Error; err != nil {
			if !xerrors.Is(err, gorm.ErrRecordNotFound) {
				return xerrors.Errorf(": %w", err)
			}

			// 存在しない時：つくる
			if err := tx.Create(conf).Error; err != nil {
				return xerrors.Errorf(": %w", err)
			}
			return nil
		}

		// 存在するとき： update
		conf.Model = c.Model
		if err := tx.Save(conf).Error; err != nil {
			return xerrors.Errorf(": %w", err)
		}
		return nil
	})
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (app *app) GetCTFConfig() (*model.Config, error) {
	var conf model.Config
	if err := app.db.First(&conf).Error; err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return &conf, nil
}
