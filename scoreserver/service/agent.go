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
	ListAgentsByIDs(ids []string) ([]*model.Agent, error)
	AgentHeartbeat(agentID, publicIP string) error
	ListAvailableAgents() ([]*model.Agent, error)
	RequestDeploy(agent *model.Agent, task *model.Challenge, teamID *uint32) (*model.Deployment, error)
}

func (app *app) GetAgentByID(agentID string) (*model.Agent, error) {
	var agent model.Agent
	if err := app.db.Where("agent_id = ?", agentID).First(&agent).Error; err != nil {
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
	t := time.Now().Add(-10 * time.Second).Unix()
	if err := app.db.Where("last_activity_at >= ?", t).Find(&agents).Error; err != nil {
		return nil, err
	}
	return agents, nil
}

// rate limit とかはこのメソッドでは気にしてない
func (app *app) RequestDeploy(agent *model.Agent, task *model.Challenge, teamID *uint32) (*model.Deployment, error) {
	now := time.Now()
	d := model.Deployment{
		ChallengeId: task.ID,
		AgentId:     agent.AgentID,
		Port:        -1,
		Status:      deployment.STATUS_WAITING,
		RequestedAt: now.Unix(),
		RetiresAt:   0,
		TeamId:      teamID,
	}
	// lifespanがあればretires atを定義する
	if task.Lifespan > 0 {
		d.RetiresAt = now.Add(time.Duration(task.Lifespan) * time.Second).Unix()
	}

	// gorm2ではこのタイミングでIDが埋められる
	err := app.db.Create(&d).Error
	if err != nil {
		return nil, err
	}
	return &d, nil
}
