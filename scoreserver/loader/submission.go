package loader

import (
	"context"

	dataloader "github.com/graph-gophers/dataloader/v6"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
	"golang.org/x/xerrors"
)

const submissionLoader = "submissionLoader"

func LoadSubmission(ctx context.Context, id uint32) (*model.Submission, error) {
	loader, err := getLoader(ctx, submissionLoader)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	thunk := loader.Load(ctx, idKey{id})
	data, err := thunk()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return data.(*model.Submission), nil
}

func newSubmissionLoader(app service.SubmissionApp) dataloader.BatchFunc {
	return func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		ids := make([]uint32, len(keys))
		for i, key := range keys {
			ids[i] = key.(idKey).ID
		}

		results := make([]*dataloader.Result, len(keys))
		submissions, _ := app.ListSubmissionByIDs(ids)
		for i, data := range submissions {
			results[i] = &dataloader.Result{Data: data, Error: nil}
		}
		return results
	}
}
