package service

import (
	"github.com/theoremoon/kosenctfx/scoreserver/repository"
	"golang.org/x/xerrors"
)

type SubmissionApp interface {
	ListSubmissions(offset, limit int64) ([]*repository.Submission, error)
	CountSubmissions() (int64, error)
}

func (app *app) ListSubmissions(offset, limit int64) ([]*repository.Submission, error) {
	submissions, err := app.repo.ListSubmissions(offset, limit)
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
