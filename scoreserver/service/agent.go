package service

import (
	"time"

	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"gorm.io/gorm/clause"
)

type AgentApp interface {
	AgentHeartbeat(agentID, publicIP string) error
}

func (app *app) GetAgentByID(agentID string) (*model.Agent, error) {
	var agent model.Agent
	if err := app.db.Where("agent_id = ?", agentID).First(&agent).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

func (app *app) AgentHeartbeat(agentID, publicIP string) error {
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
