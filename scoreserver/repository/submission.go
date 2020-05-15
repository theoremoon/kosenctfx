package repository

import (
	"github.com/theoremoon/kosenctfx/scoreserver/model"
)

type SubmissionRepository interface {
	ListValidSubmissions() ([]*model.Submission, error)
	FindValidSubmission(userId uint, challengeId uint) (*model.Submission, error)
	InsertSubmission(s *model.Submission) error
}

func (r *repository) ListValidSubmissions() ([]*model.Submission, error) {
	var submissions []*model.Submission
	if err := r.db.Where("is_valid = ?", true).Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
}

func (r *repository) FindValidSubmission(userId uint, challengeId uint) (*model.Submission, error) {
	var submission model.Submission
	if err := r.db.Where("user_id = ? AND challenge_id = ? AND is_valid = ?", userId, challengeId, true).First(&submission).Error; err != nil {
		return nil, err
	}
	return &submission, nil
}

func (r *repository) InsertSubmission(s *model.Submission) error {
	if err := r.db.Create(s).Error; err != nil {
		return err
	}
	return nil
}
