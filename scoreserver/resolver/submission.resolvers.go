package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strconv"

	"github.com/theoremoon/kosenctfx/scoreserver/loader"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/repository"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
	"golang.org/x/xerrors"
)

func (r *submissionResolver) ID(ctx context.Context, obj *repository.Submission) (string, error) {
	if err := checkPermission(ctx, PERM_ADMIN); err != nil {
		return "", xerrors.Errorf(": %w", err)
	}
	return strconv.FormatUint(uint64(obj.ID), 10), nil
}

func (r *submissionResolver) Team(ctx context.Context, obj *repository.Submission) (*model.Team, error) {
	if err := checkPermission(ctx, PERM_ADMIN); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	t, err := r.app.GetTeamByID(obj.TeamId)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return t, nil
}

func (r *submissionResolver) Challenge(ctx context.Context, obj *repository.Submission) (*service.Challenge, error) {
	if err := checkPermission(ctx, PERM_ADMIN); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	if obj.ChallengeId == nil {
		return nil, nil
	}

	c, err := loader.LoadChallenge(ctx, *obj.ChallengeId)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return c, nil
}

// Submission returns SubmissionResolver implementation.
func (r *Resolver) Submission() SubmissionResolver { return &submissionResolver{r} }

type submissionResolver struct{ *Resolver }
