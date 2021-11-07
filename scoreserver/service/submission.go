package service

import (
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
}

func (app *app) ListSubmissions(offset, limit int64) ([]*model.Submission, error) {
	submissions, err := app.repo.ListSubmissions(offset, limit)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return submissions, nil
}

func (app *app) ListSubmissionByIDs(ids []uint32) ([]*model.Submission, error) {
	submissions, err := app.repo.ListSubmissionByIDs(ids)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return submissions, nil
}

func (app *app) ListValidSubmissions() ([]*model.Submission, error) {
	// model.Submissionのis_validカラムよりはこちらが信用できる
	valid_submissions, err := app.repo.ListValidSubmissions()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	submission_ids := make([]uint32, len(valid_submissions))
	for i, s := range valid_submissions {
		submission_ids[i] = s.SubmissionId
	}
	submissions, err := app.repo.ListSubmissionByIDs(submission_ids)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return submissions, nil
}
func (app *app) ListTeamSubmissions(teamID uint32) ([]*model.Submission, error) {
	submissions, err := app.repo.ListTeamSubmissions(teamID)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return submissions, nil
}

func (app *app) CountSubmissions() (int64, error) {
	count, err := app.repo.CountSubmissions()
	if err != nil {
		return 0, xerrors.Errorf(": %w", err)
	}
	return count, nil
}

func (app *app) CountValidSubmissions() (int64, error) {
	count, err := app.repo.CountValidSubmissions()
	if err != nil {
		return 0, xerrors.Errorf(": %w", err)
	}
	return count, nil
}
