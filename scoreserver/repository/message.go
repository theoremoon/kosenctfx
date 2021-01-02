package repository

import (
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"golang.org/x/xerrors"
)

type MessageRepository interface {
	GetMessage(key string) (*model.Message, error)
	SetMessage(key, value string) error
}

func (r *repository) GetMessage(key string) (*model.Message, error) {

	var m model.Message
	if err := r.db.Where("key = ?", key).First(&m).Error; err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return &m, nil
}

func (r *repository) SetMessage(key, value string) error {
	if err := r.db.Create(&model.Message{
		Key:   key,
		Value: value,
	}).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}
