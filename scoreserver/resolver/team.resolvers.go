package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strconv"

	"github.com/theoremoon/kosenctfx/scoreserver/model"
)

func (r *teamResolver) ID(ctx context.Context, obj *model.Team) (string, error) {
	return strconv.FormatUint(uint64(obj.ID), 10), nil
}

func (r *teamResolver) Name(ctx context.Context, obj *model.Team) (string, error) {
	return obj.Teamname, nil
}

// Team returns TeamResolver implementation.
func (r *Resolver) Team() TeamResolver { return &teamResolver{r} }

type teamResolver struct{ *Resolver }
