package service

import (
	"github.com/theoremoon/kosenctfx/scoreserver/model"
)

type ClarificationResponseType uint

const (
	ClarificationAnswerYes ClarificationResponseType = iota
	ClarificationAnswerNo
	ClarificationAnswerNoComment
	ClarificationAnswerNotRelated
	ClarificationAnswerSeeTheChallenge
	ClarificationAnswerSupportIndivisually
)

type ClarificationApp interface {
	NewClarification(user *model.User, content string) (*model.Clarification, error)
	ListOpenClarifications() ([]*model.Clarification, error)
	UpdateClarification(id uint, response ClarificationResponseType, isCompleted, isPublic bool) error
}

func (app *app) NewClarification(user *model.User, content string) (*model.Clarification, error) {
	return nil, ErrorMessage("not implemented")
}
func (app *app) ListOpenClarifications() ([]*model.Clarification, error) {
	return nil, ErrorMessage("not implemented")
}
func (app *app) UpdateClarification(id uint, response ClarificationResponseType, isCompleted, isPublic bool) error {
	return ErrorMessage("not implemented")
}
