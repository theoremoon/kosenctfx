package loader

import (
	"context"
	"strconv"

	dataloader "github.com/graph-gophers/dataloader/v6"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
	"golang.org/x/xerrors"
)

type idKey struct {
	ID uint32
}

func (k idKey) String() string {
	return strconv.FormatUint(uint64(k.ID), 10)
}

func (k idKey) Raw() interface{} {
	return k.ID
}

func Attach(app service.App, ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, submissionLoader, dataloader.NewBatchedLoader(newSubmissionLoader(app)))
	ctx = context.WithValue(ctx, challengeLoader, dataloader.NewBatchedLoader(newChallengeLoader(app)))
	return ctx
}

func getLoader(ctx context.Context, key string) (*dataloader.Loader, error) {
	ldr, ok := ctx.Value(key).(*dataloader.Loader)
	if !ok {
		return nil, xerrors.Errorf("no such loader: %s", key)
	}
	return ldr, nil
}
