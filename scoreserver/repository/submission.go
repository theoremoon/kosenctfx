package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"golang.org/x/xerrors"
)

type SubmissionRepository interface {
	ListValidSubmissions() ([]*model.Submission, error)
	FindValidSubmission(userId uint, challengeId uint) (*model.Submission, error)
	InsertSubmission(s *model.Submission) error

	InsertSubmissionTx(s *model.Submission) (*model.Submission, bool, error)
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

/// トランザクションロックを掛けながら Submission を行う
/// 返り値のboolはtrueならvalid, falseならvaildではないということになる
func (r *repository) InsertSubmissionTx(s *model.Submission) (*model.Submission, bool, error) {
	valid := true
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var submission model.Submission
		if err := tx.Where("user_id = ? AND challenge_id = ? AND is_valid = ?", s.UserId, s.ChallengeId, true).First(&submission).Error; err != nil {
			if !gorm.IsRecordNotFoundError(err) {
				return xerrors.Errorf(": %w", err)
			}
			valid = false
		}
		s.IsValid = valid
		if err := tx.Create(s).Error; err != nil {
			return xerrors.Errorf(": %w", err)
		}
		return nil
	}); err != nil {
		return nil, false, xerrors.Errorf(": %w", err)
	}

	return s, valid, nil
}
