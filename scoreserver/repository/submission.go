package repository

import (
	"time"

	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

type SubmissionRepository interface {
	ListSubmissions(offset, limit int64) ([]*Submission, error)
	ListValidSubmissions() ([]*model.ValidSubmission, error)
	InsertSubmission(s *model.Submission) error
	InsertValidableSubmission(s *model.Submission) (bool, error)
	CountSubmissions() (int64, error)

	GetWrongCount(teamID uint32, duration time.Duration) (int64, error)
	LockSubmission(teamID uint32, duration time.Duration) error
	CheckSubmittable(teamID uint32) (bool, error)
}

type Submission struct {
	ID          uint32  `gorm:"column:id"`
	TeamId      uint32  `gorm:"column:team_id"`
	ChallengeId *uint32 `gorm:"column:challenge_id"`
	Flag        string  `gorm:"column:flag"`
	IsCorrect   bool    `gorm:"column:is_correct"`
	IsValid     bool    `gorm:"column:is_valid"`
	SubmittedAt int64   `gorm:"column:submitted_at"`
}

func (r *repository) ListSubmissions(offset, limit int64) ([]*Submission, error) {
	rows, err := r.db.Model(&model.Submission{}).Order("created_at desc").Offset(int(offset)).Limit(int(offset)).Select("submissions.id AS id, submissions.team_id AS team_id, submissions.challenge_id AS challenge_id, submissions.flag AS flag, submissions.is_correct AS is_correct, (valid_submissions.is_valid IS NOT NULL) AS is_valid, submissions.created_at as submitted_at").Joins("left outer join").Rows()
	defer rows.Close()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	submissions := make([]*Submission, 0, 1000)
	for rows.Next() {
		var submission Submission
		r.db.ScanRows(rows, &submission)
		submissions = append(submissions, &submission)
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

func (r *repository) CountSubmissions() (int64, error) {
	var count int64
	if err := r.db.Model(&model.Submission{}).Count(&count).Error; err != nil {
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
