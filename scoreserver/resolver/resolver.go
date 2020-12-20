package resolver

//go:generate go run github.com/99designs/gqlgen

import (
	"context"

	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
	"golang.org/x/xerrors"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

const teamKey = "github.com/theoremoon/kosenctfx/scoreserver/resolver::resolver"
const (
	PERM_CONTESTANT = iota
	PERM_ADMIN
)

func AttachTeam(ctx context.Context, team *model.Team) context.Context {
	return context.WithValue(ctx, teamKey, team)
}

func getLoginTeam(ctx context.Context) (*model.Team, error) {
	team := ctx.Value(teamKey).(*model.Team)
	if team == nil {
		return nil, xerrors.New("unauthorized")
	}
	return team, nil
}

func checkPermission(ctx context.Context, permission int) error {
	t, err := getLoginTeam(ctx)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	if permission == PERM_ADMIN {
		if !t.IsAdmin {
			return xerrors.New("unauthorized")
		}
	}
	return nil
}

type Resolver struct {
	app service.App
}

func NewResolver(app service.App) ResolverRoot {
	return &Resolver{
		app: app,
	}
}
