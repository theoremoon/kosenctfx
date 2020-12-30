package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strconv"

	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/repository"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
	"golang.org/x/xerrors"
)

func (r *queryResolver) GetChallenge(ctx context.Context, id string) (*service.Challenge, error) {
	if err := checkPermission(ctx, PERM_CONTESTANT); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	cID, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	c, err := r.app.GetChallengeByID(uint32(cID))
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	if err := checkPermission(ctx, PERM_ADMIN); err == nil {
		return c, nil
	} else {
		if !c.IsOpen {
			return nil, xerrors.New("no such challenge")
		}
		c.Flag = ""
		return c, nil
	}
}

func (r *queryResolver) ListChallenges(ctx context.Context) ([]*service.Challenge, error) {
	// TODO: contestant側もこちらを使うように改造……するかもしれない
	if err := checkPermission(ctx, PERM_ADMIN); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	chals, err := r.app.ListAllRawChallenges()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	teams, err := r.app.ListTeams()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	// TODO: scorefeedを毎回計算すると重たいのでcacheにのせたりする……かもしれない
	cs, _, err := r.app.ScoreFeed(chals, teams)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return cs, nil
}

func (r *queryResolver) GetTeam(ctx context.Context, id string) (*model.Team, error) {
	tID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	t, err := r.app.GetTeamByID(uint32(tID))
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return t, nil
}

func (r *queryResolver) ListSubmissions(ctx context.Context, page PaginationInput) ([]*repository.Submission, error) {
	if err := checkPermission(ctx, PERM_ADMIN); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	submissions, err := r.app.ListSubmissions(int64(page.Offset), int64(page.Limit))
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return submissions, nil
}

func (r *queryResolver) GetNumberOfSubmissions(ctx context.Context) (int, error) {
	if err := checkPermission(ctx, PERM_ADMIN); err != nil {
		return 0, xerrors.Errorf(": %w", err)
	}

	count, err := r.app.CountSubmissions()
	if err != nil {
		return 0, xerrors.Errorf(": %w", err)
	}
	return int(count), nil
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
