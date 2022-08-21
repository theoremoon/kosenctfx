package service

import (
	"time"

	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/task/registry"
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
	SetRegistryConfig(*registry.RegistryConfig) error
	GetRegistryConfig() (*registry.RegistryConfig, error)
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

// registryConfigが正しいかどうかはこのメソッドの責務の範囲外
func (app *app) SetRegistryConfig(registry *registry.RegistryConfig) error {
	conf, err := app.GetCTFConfig()
	if err != nil {
		return err
	}

	conf.RegistryURL = registry.URL
	conf.RegistryUser = registry.Username
	conf.RegistryPassword = registry.Password

	err = app.SetCTFConfig(conf)
	if err != nil {
		return err
	}
	return nil
}

func (app *app) GetRegistryConfig() (*registry.RegistryConfig, error) {
	conf, err := app.GetCTFConfig()
	if err != nil {
		return nil, err
	}

	reg := registry.RegistryConfig{
		URL:      conf.RegistryURL,
		Username: conf.RegistryUser,
		Password: conf.RegistryPassword,
	}

	return &reg, nil
}
