package loader

import (
	"context"

	dataloader "github.com/graph-gophers/dataloader/v6"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
	"golang.org/x/xerrors"
)

const challengeLoader = "challengeLoader"

func LoadChallenge(ctx context.Context, id uint32) (*service.Challenge, error) {
	loader, err := getLoader(ctx, challengeLoader)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	thunk := loader.Load(ctx, idKey{id})
	data, err := thunk()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return data.(*service.Challenge), nil
}

func newChallengeLoader(app service.ChallengeApp) dataloader.BatchFunc {
	return func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		ids := make([]uint32, len(keys))
		for i, key := range keys {
			ids[i] = key.(idKey).ID
		}

		results := make([]*dataloader.Result, len(keys))
		cs, _ := app.ListChallengeByIDs(ids)
		for i, data := range cs {
			results[i] = &dataloader.Result{Data: data, Error: nil}
		}
		return results
	}
}
