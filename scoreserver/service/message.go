package service

import (
	"golang.org/x/xerrors"
)

type MessageApp interface {
	GetMessage(key string) (string, error)
	SetMessage(key, value string) error
}

func (app *app) GetMessage(key string) (string, error) {
	m, err := app.repo.GetMessage(key)
	if err != nil {
		return "", xerrors.Errorf(": %w", err)
	}
	return m.Value, nil
}

func (app *app) SetMessage(key, value string) error {
	err := app.repo.SetMessage(key, value)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}
