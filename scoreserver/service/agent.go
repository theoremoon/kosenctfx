package service

import (
	"errors"
	"time"

	"github.com/theoremoon/kosenctfx/scoreserver/deployment"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"gorm.io/gorm/clause"
)

type AgentApp interface {
	GetAgentByID(agentID string) (*model.Agent, error)
	GetAgentForRequestDeployment() (*model.Agent, error)
	ListAgentsByIDs(ids []string) ([]*model.Agent, error)
	AgentHeartbeat(agentID, publicIP string) error
	ListAvailableAgents() ([]*model.Agent, error)
}

// この時間以内に最終アクセスがあったエージェントは活きているとみなす
const (
	LIVING_THRESHOLD = 15 * time.Second
)

func (app *app) GetAgentByID(agentID string) (*model.Agent, error) {
	var agent model.Agent
	if err := app.db.Where("agent_id = ?", agentID).First(&agent).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

// このメソッドで引いてきたエージェントにタスクをデプロイすることになる
// ある程度負荷がばらけるように引いてくる
func (app *app) GetAgentForRequestDeployment() (*model.Agent, error) {
	var agent model.Agent
	// 生きてるエージェントの中で
	// 一番動いてるタスクが少ないエージェントを引いてくる
	err := app.db.Joins("JOIN deployments ON deployments.agent_id = agents.agent_id").
		Where("deployments.status in ? AND agents.last_activity_at >= ?", deployment.LivingStatuses, time.Now().Add(-LIVING_THRESHOLD)).
		Order("COUNT(deployments) DESC").
		First(&agent).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (app *app) ListAgentsByIDs(ids []string) ([]*model.Agent, error) {
	var agents []*model.Agent
	if err := app.db.Where("agent_id IN ?", ids).Find(&agents).Error; err != nil {
		return nil, err
	}
	return agents, nil
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
	t := time.Now().Add(-LIVING_THRESHOLD).Unix()
	if err := app.db.Where("last_activity_at >= ?", t).Find(&agents).Error; err != nil {
		return nil, err
	}
	return agents, nil
}
