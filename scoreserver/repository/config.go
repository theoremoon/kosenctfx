package repository

import (
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"golang.org/x/xerrors"
)

type ConfigRepository interface {
	SetConfig(conf *model.Config) error
	GetConfig() (*model.Config, error)
}

func (r *repository) SetConfig(conf *model.Config) error {
	if err := r.db.Create(conf).Error; err != nil {
		return err
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
