package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strconv"

	"github.com/theoremoon/kosenctfx/scoreserver/service"
)

func (r *challengeResolver) ID(ctx context.Context, obj *service.Challenge) (string, error) {
	return strconv.FormatUint(uint64(obj.ID), 10), nil
}

func (r *challengeResolver) Score(ctx context.Context, obj *service.Challenge) (int, error) {
	return int(obj.Score), nil
}

func (r *solvedByResolver) TeamID(ctx context.Context, obj *service.SolvedBy) (string, error) {
	return strconv.FormatUint(uint64(obj.TeamID), 10), nil
}

// Challenge returns ChallengeResolver implementation.
func (r *Resolver) Challenge() ChallengeResolver { return &challengeResolver{r} }

// SolvedBy returns SolvedByResolver implementation.
func (r *Resolver) SolvedBy() SolvedByResolver { return &solvedByResolver{r} }

type challengeResolver struct{ *Resolver }
type solvedByResolver struct{ *Resolver }
