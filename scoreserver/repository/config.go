package repository

import (
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

type ConfigRepository interface {
	SetConfig(conf *model.Config) error
	GetConfig() (*model.Config, error)
}

func (r *repository) SetConfig(conf *model.Config) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
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

func (r *repository) GetConfig() (*model.Config, error) {
	var conf model.Config
	if err := r.db.First(&conf).Error; err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return &conf, nil
}
