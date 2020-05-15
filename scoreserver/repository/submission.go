package repository

import (
	"github.com/theoremoon/kosenctfx/scoreserver/model"
)

type SubmissionRepository interface {
	ListValidSubmissions() ([]*model.Submission, error)
}

func (r *repository) ListValidSubmissions() ([]*model.Submission, error) {
	var submissions []*model.Submission
	if err := r.db.Where("is_valid = ?", true).Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
}
