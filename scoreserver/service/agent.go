package service

import (
	"errors"
	"time"

	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"gorm.io/gorm/clause"
)

type AgentApp interface {
	AgentHeartbeat(agentID, publicIP string) error
	ListAvailableAgents() ([]*model.Agent, error)
}

func (app *app) GetAgentByID(agentID string) (*model.Agent, error) {
	var agent model.Agent
	if err := app.db.Where("agent_id = ?", agentID).First(&agent).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

func (app *app) AgentHeartbeat(agentID, publicIP string) error {
	if agentID == "" {
		return errors.New("agent id must be specified")
	}

	// build agent model
	agent := &model.Agent{
		AgentID:        agentID,
		PublicIP:       publicIP,
		LastActivityAt: time.Now().Unix(),
	}

	// upsert
	err := app.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&agent).Error
	if err != nil {
		return err
	}
	return nil
}

func (app *app) ListAvailableAgents() ([]*model.Agent, error) {
	var agents []*model.Agent
	t := time.Now().Add(-10 * time.Second).Unix()
	if err := app.db.Where("last_activity_at >= ?", t).Find(&agents).Error; err != nil {
		return nil, err
	}
	return agents, nil
}
