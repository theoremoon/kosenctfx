package repository

import (
	"github.com/theoremoon/kosenctfx/scoreserver/model"
)

type TeamRepository interface {
	RegisterTeam(t *model.Team) error
	ListAllTeams() ([]*model.Team, error)
	GetTeamByID(teamId uint) (*model.Team, error)
	GetTeamByName(teamname string) (*model.Team, error)
	GetTeamByToken(token string) (*model.Team, error)
	UpdateTeamToken(t *model.Team, token string) error
}

func (r *repository) RegisterTeam(t *model.Team) error {
	if err := r.db.Create(t).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) ListAllTeams() ([]*model.Team, error) {
	var teams []*model.Team
	if err := r.db.Find(&teams).Error; err != nil {
		return nil, err
	}
	return teams, nil
}

func (r *repository) GetTeamByID(teamID uint) (*model.Team, error) {
	var t model.Team
	if err := r.db.Where("id = ?", teamID).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *repository) GetTeamByName(teamname string) (*model.Team, error) {
	var t model.Team
	if err := r.db.Where("teamname = ?", teamname).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *repository) GetTeamByToken(token string) (*model.Team, error) {
	var t model.Team
	if err := r.db.Where("token = ?", token).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *repository) UpdateTeamToken(t *model.Team, token string) error {
	if err := r.db.Model(t).Update("token", token).Error; err != nil {
		return err
	}
	return nil
}
