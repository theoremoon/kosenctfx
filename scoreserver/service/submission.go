package service

import (
	"time"

	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"golang.org/x/xerrors"
)

type SubmissionApp interface {
	ListSubmissions(offset, limit int64) ([]*model.Submission, error)
	ListSubmissionByIDs(ids []uint32) ([]*model.Submission, error)
	ListValidSubmissions() ([]*model.Submission, error)
	ListTeamSubmissions(teamID uint32) ([]*model.Submission, error)
	CountSubmissions() (int64, error)
	CountValidSubmissions() (int64, error)

	GetWrongCount(teamID uint32, duration time.Duration) (int64, error)
	LockSubmission(teamID uint32, duration time.Duration) error
	CheckSubmittable(teamID uint32) (bool, error)
}

func (app *app) ListSubmissions(offset, limit int64) ([]*model.Submission, error) {
	var submissions []*model.Submission
	if err := app.db.Order("submissions.created_at desc").Offset(int(offset)).Limit(int(offset)).Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
}

func (app *app) ListSubmissionByIDs(ids []uint32) ([]*model.Submission, error) {
	var submissions []*model.Submission
	if err := app.db.Order("submissions.created_at desc").Where("submissions.id IN ?", ids).Find(&submissions).Error; err != nil {
		return nil, err
	}

	return submissions, nil
}

func (app *app) listValidSubmissions() ([]*model.ValidSubmission, error) {
	var submissions []*model.ValidSubmission
	if err := app.db.Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
}

func (app *app) ListValidSubmissions() ([]*model.Submission, error) {
	// model.Submissionのis_validカラムよりはこちらが信用できる
	valid_submissions, err := app.listValidSubmissions()
	if err != nil {
		return nil, err
	}
	submission_ids := make([]uint32, len(valid_submissions))
	for i, s := range valid_submissions {
		submission_ids[i] = s.SubmissionId
	}
	submissions, err := app.ListSubmissionByIDs(submission_ids)
	if err != nil {
		return nil, err
	}
	return submissions, nil
}

func (app *app) ListTeamSubmissions(teamID uint32) ([]*model.Submission, error) {
	var submissions []*model.Submission
	if err := app.db.Where("team_id = ?", teamID).Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
}

func (app *app) CountSubmissions() (int64, error) {
	var count int64
	if err := app.db.Model(&model.Submission{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (app *app) CountValidSubmissions() (int64, error) {
	var count int64
	if err := app.db.Model(&model.ValidSubmission{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (app *app) GetWrongCount(teamID uint32, duration time.Duration) (int64, error) {
	t := time.Now().Add(-duration).Unix()
	var count int64
	if err := app.db.Model(&model.Submission{}).Where("team_id = ? AND is_correct = ? AND created_at > ?", teamID, false, t).Count(&count).Error; err != nil {
		return 0, xerrors.Errorf(": %w", err)
	}
	return count, nil
}

func (app *app) LockSubmission(teamID uint32, duration time.Duration) error {
	if err := app.db.Create(&model.SubmissionLock{
		TeamId: teamID,
		Until:  time.Now().Add(duration).Unix(),
	}).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (app *app) CheckSubmittable(teamID uint32) (bool, error) {
	var count int64
	if err := app.db.Model(&model.SubmissionLock{}).Where("team_id = ? AND until >= ?", teamID, time.Now().Unix()).Count(&count).Error; err != nil {
		return false, xerrors.Errorf(": %w", err)
	}
	return count == 0, nil
}
