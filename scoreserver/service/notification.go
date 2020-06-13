package service

import (
	"github.com/theoremoon/kosenctfx/scoreserver/model"
)

type NotificationApp interface {
	ListNotifications() ([]*model.Notification, error)
}

func (app *app) ListNotifications() ([]*model.Notification, error) {
	return nil, ErrorMessage("not implemented")
}
