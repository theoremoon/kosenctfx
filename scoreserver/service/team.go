package service

import (
	"github.com/jinzhu/gorm"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
)

type Team struct {
}

type TeamApp interface {
	RegisterTeam(teamname string) (*model.Team, error)
	RegisterAdminTeam(teamname string) (*model.Team, error)
	GetTeamByID(teamID uint) (*model.Team, error)
	GetUserTeam(userID uint) (*model.Team, error)
	UpdateTeamname(teamID uint, newTeamname string) error
}

func (app *app) RegisterTeam(teamname string) (*model.Team, error) {
	if err := app.validateTeamname(teamname); err != nil {
		return nil, err
	}

	t := model.Team{
		Teamname: teamname,
		Token:    newToken(),
	}
	if err := app.repo.RegisterTeam(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (app *app) RegisterAdminTeam(teamname string) (*model.Team, error) {
	if err := app.validateTeamname(teamname); err != nil {
		return nil, err
	}

	t := model.Team{
		Teamname: teamname,
		Token:    newToken(),
		IsAdmin:  true,
	}
	if err := app.repo.RegisterTeam(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (app *app) GetTeamByID(teamID uint) (*model.Team, error) {
	team, err := app.repo.GetTeamByID(teamID)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (app *app) GetUserTeam(userID uint) (*model.Team, error) {
	return nil, ErrorMessage("not implemented")
}

func (app *app) UpdateTeamname(teamID uint, newTeamname string) error {
	return ErrorMessage("not implemented")
}

func (app *app) validateTeamname(teamname string) error {
	if teamname == "" {
		return ErrorMessage("teamname is required")
	}
	if len(teamname) >= 128 {
		return ErrorMessage("teamname should be shorter than 128")
	}

	if _, err := app.repo.GetTeamByName(teamname); err == nil {
		return ErrorMessage("teamname already used")
	} else if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}
	return nil
}
