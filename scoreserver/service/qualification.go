package service

import (
	"github.com/theoremoon/kosenctfx/scoreserver/model"
)

type QualificationApp interface {
	NewQualification(user *model.User, content string) (*model.Qualification, error)
}

func (app *app) NewQualification(user *model.User, content string) (*model.Qualification, error) {
	return nil, ErrorMessage("not implemented")
}
