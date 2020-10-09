package repository

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"golang.org/x/xerrors"
)

type SubmissionRepository interface {
	ListValidSubmissions() ([]*model.ValidSubmission, error)
	InsertSubmission(s *model.Submission) error
	InsertValidableSubmission(s *model.Submission) (bool, error)

	GetWrongCount(teamID uint, duration time.Duration) (int, error)
	LockSubmission(teamID uint, duration time.Duration) error
	CheckSubmittable(teamID uint) (bool, error)
}

func (r *repository) ListValidSubmissions() ([]*model.ValidSubmission, error) {
	var submissions []*model.ValidSubmission
	if err := r.db.Find(&submissions).Error; err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return submissions, nil
}

func (r *repository) InsertSubmission(s *model.Submission) error {
	if err := r.db.Create(s).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (r *repository) InsertValidableSubmission(s *model.Submission) (bool, error) {
	valid := false

	// ややこしいのでtransactionをとって行う
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// とりあえずsubmissionは保存しておく
		if err := r.db.Create(s).Error; err != nil {
			return xerrors.Errorf(": %w", err)
		}

		// 既存の提出を読んでvalidityを決定する
		var count int
		if err := r.db.Model(&model.ValidSubmission{}).Where("team_id = ? AND challenge_id = ?", s.TeamId, s.ChallengeId).Count(&count).Error; err != nil {
			return xerrors.Errorf(": %w", err)
		}
		if count == 0 {
			valid = true
		}

		if valid {
			// validならvalid submissionを作成する
			vs := model.ValidSubmission{
				SubmissionId: s.ID,
				ChallengeId:  *s.ChallengeId,
				TeamId:       s.TeamId,
			}
			if err := r.db.Create(&vs).Error; err != nil {
				// ただしConstraint Errorが起きたらやはりValidではなかった
				if isDuplicatedError(err) {
					valid = false
					return nil
				} else {
					return xerrors.Errorf(": %w", err)
				}
			}
		}
		return nil
	})
	if err != nil {
		return false, xerrors.Errorf(": %w", err)
	}
	return valid, nil
}

func (r *repository) GetWrongCount(teamID uint, duration time.Duration) (int, error) {
	t := time.Now().Add(-duration)
	var count int
	if err := r.db.Model(&model.Submission{}).Where("team_id = ? AND is_correct = ? AND created_at > ?", teamID, false, t).Count(&count).Error; err != nil {
		return 0, xerrors.Errorf(": %w", err)
	}
	return count, nil
}

func (r *repository) LockSubmission(teamID uint, duration time.Duration) error {
	if err := r.db.Create(&model.SubmissionLock{
		TeamId: teamID,
		Until:  time.Now().Add(duration),
	}).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (r *repository) CheckSubmittable(teamID uint) (bool, error) {
	var count uint
	if err := r.db.Model(&model.SubmissionLock{}).Where("team_id = ? AND until >= ?", teamID, time.Now()).Count(&count).Error; err != nil {
		return false, xerrors.Errorf(": %w", err)
	}
	return count == 0, nil
}
