package repository

import (
	"time"

	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

type SubmissionRepository interface {
	ListSubmissions(offset, limit int64) ([]*model.Submission, error)
	ListSubmissionByIDs(ids []uint32) ([]*model.Submission, error)
	ListValidSubmissions() ([]*model.ValidSubmission, error)
	ListTeamSubmissions(teamID uint32) ([]*model.Submission, error)
	InsertSubmission(s *model.Submission) error
	InsertValidableSubmission(s *model.Submission) (bool, error)
	MarkSubmissionValid(id uint32) error
	CountSubmissions() (int64, error)
	CountValidSubmissions() (int64, error)

	GetWrongCount(teamID uint32, duration time.Duration) (int64, error)
	LockSubmission(teamID uint32, duration time.Duration) error
	CheckSubmittable(teamID uint32) (bool, error)
}

func (r *repository) ListSubmissions(offset, limit int64) ([]*model.Submission, error) {
	var submissions []*model.Submission
	if err := r.db.Order("submissions.created_at desc").Offset(int(offset)).Limit(int(offset)).Find(&submissions).Error; err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return submissions, nil
}

func (r *repository) ListSubmissionByIDs(ids []uint32) ([]*model.Submission, error) {
	var submissions []*model.Submission
	if err := r.db.Order("submissions.created_at desc").Where("submissions.id IN ?", ids).Find(&submissions).Error; err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return submissions, nil
}

func (r *repository) ListValidSubmissions() ([]*model.ValidSubmission, error) {
	var submissions []*model.ValidSubmission
	if err := r.db.Find(&submissions).Error; err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return submissions, nil
}

func (r *repository) ListTeamSubmissions(teamID uint32) ([]*model.Submission, error) {
	var submissions []*model.Submission
	if err := r.db.Where("team_id = ?", teamID).Find(&submissions).Error; err != nil {
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
		var count int64
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

func (r *repository) MarkSubmissionValid(id uint32) error {
	if err := r.db.Model(&model.Submission{}).Where("id = ?", id).Update("is_valid", true).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (r *repository) CountSubmissions() (int64, error) {
	var count int64
	if err := r.db.Model(&model.Submission{}).Count(&count).Error; err != nil {
		return 0, xerrors.Errorf(": %w", err)
	}
	return count, nil
}

func (r *repository) CountValidSubmissions() (int64, error) {
	var count int64
	if err := r.db.Model(&model.ValidSubmission{}).Count(&count).Error; err != nil {
		return 0, xerrors.Errorf(": %w", err)
	}
	return count, nil
}

func (r *repository) GetWrongCount(teamID uint32, duration time.Duration) (int64, error) {
	t := time.Now().Add(-duration).Unix()
	var count int64
	if err := r.db.Model(&model.Submission{}).Where("team_id = ? AND is_correct = ? AND created_at > ?", teamID, false, t).Count(&count).Error; err != nil {
		return 0, xerrors.Errorf(": %w", err)
	}
	return count, nil
}

func (r *repository) LockSubmission(teamID uint32, duration time.Duration) error {
	if err := r.db.Create(&model.SubmissionLock{
		TeamId: teamID,
		Until:  time.Now().Add(duration).Unix(),
	}).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (r *repository) CheckSubmittable(teamID uint32) (bool, error) {
	var count int64
	if err := r.db.Model(&model.SubmissionLock{}).Where("team_id = ? AND until >= ?", teamID, time.Now().Unix()).Count(&count).Error; err != nil {
		return false, xerrors.Errorf(": %w", err)
	}
	return count == 0, nil
}
