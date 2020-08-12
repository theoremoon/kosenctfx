package service

import (
	"github.com/theoremoon/kosenctfx/scoreserver/model"
)

type NotificationApp interface {
	ListNotifications() ([]*model.Notification, error)
	AddNotification(content string) (*model.Notification, error)
}

func (app *app) ListNotifications() ([]*model.Notification, error) {
	return nil, NewErrorMessage("not implemented")
}
func (app *app) AddNotification(content string) (*model.Notification, error) {
	return nil, NewErrorMessage("not implemented")
}
