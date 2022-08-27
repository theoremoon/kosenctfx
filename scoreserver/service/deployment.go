package service

import (
	"time"

	"github.com/theoremoon/kosenctfx/scoreserver/deployment"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
)

type DeploymentApp interface {
	GetDeploymentByID(id uint32) (*model.Deployment, error)
	GetHostAndPortByTeamAndTaskID(teamID, taskID uint32) (*HostAndPort, error)
	ListLivingDeployments() ([]*model.Deployment, error)
	ListLivingDeploymentsByTeamID(teamID uint32) ([]*model.Deployment, error)
	ListDeploymentRequestForAgent(agentID string) ([]*model.Deployment, error)
	ListRetireRequestForAgent(agentID string) ([]*model.Deployment, error)
	StartDeployment(deploymentID uint32, port int) error
	UpdateDeploymentStatus(deploymentID uint32, status string) error
	RequestDeploy(agent *model.Agent, task *model.Challenge, teamID *uint32) (*model.Deployment, error)
	RequestRetire(d *model.Deployment) error
}

func (app *app) GetDeploymentByID(id uint32) (*model.Deployment, error) {
	var deployment model.Deployment
	err := app.db.Where("id = ?", id).First(&deployment).Error
	if err != nil {
		return nil, err
	}
	return &deployment, nil
}

type HostAndPort struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (app *app) GetHostAndPortByTeamAndTaskID(teamID, taskID uint32) (*HostAndPort, error) {
	var res HostAndPort
	err := app.db.Table("deployments").
		Select("agents.public_ip as host, deployments.port as port").
		Joins("left join agents on deployments.agent_id = agents.agent_id").
		Where("(deployments.team_id = ? OR deployments.team_id = NULL) AND deployments.challenge_id = ? AND deployments.status = ?", teamID, taskID, deployment.STATUS_AVAILABLE).
		First(&res).
		Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (app *app) ListLivingDeploymentsByTeamID(teamID uint32) ([]*model.Deployment, error) {
	var deployments []*model.Deployment
	err := app.db.
		Where("status IN ? AND team_id = ?", deployment.LivingStatuses, teamID).
		Order("requested_at desc").
		Find(&deployments).Error
	if err != nil {
		return nil, err
	}
	return deployments, nil
}

func (app *app) ListLivingDeployments() ([]*model.Deployment, error) {
	var deployments []*model.Deployment
	err := app.db.
		Where("status IN ?", deployment.LivingStatuses).
		Order("requested_at desc").
		Find(&deployments).Error
	if err != nil {
		return nil, err
	}
	return deployments, nil
}

func (app *app) ListDeploymentRequestForAgent(agentID string) ([]*model.Deployment, error) {
	var deployments []*model.Deployment
	err := app.db.
		Where("status = ? AND agent_id = ?", deployment.STATUS_WAITING, agentID).
		Order("requested_at desc").
		Find(&deployments).Error
	if err != nil {
		return nil, err
	}
	return deployments, nil
}

func (app *app) ListRetireRequestForAgent(agentID string) ([]*model.Deployment, error) {
	var deployments []*model.Deployment
	err := app.db.
		Where("status = ? AND agent_id = ?", deployment.STATUS_RETIRE_REQUESTING, agentID).
		Order("requested_at desc").
		Find(&deployments).Error
	if err != nil {
		return nil, err
	}
	return deployments, nil
}

func (app *app) StartDeployment(deploymentID uint32, port int) error {
	err := app.db.
		Model(&model.Deployment{}).
		Where("id = ?", deploymentID).
		Updates(map[string]interface{}{
			"status": deployment.STATUS_DEPLOYING,
			"port":   port,
		}).Error
	if err != nil {
		return err
	}
	return nil
}

func (app *app) UpdateDeploymentStatus(deploymentID uint32, status string) error {
	err := deployment.ValidateStatus(status)
	if err != nil {
		return err
	}
	err = app.db.
		Model(&model.Deployment{}).
		Where("id = ?", deploymentID).
		Updates(map[string]interface{}{
			"status": status,
		}).Error
	if err != nil {
		return err
	}
	return nil
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

func (app *app) RequestRetire(d *model.Deployment) error {
	err := app.UpdateDeploymentStatus(d.ID, deployment.STATUS_RETIRE_REQUESTING)
	if err != nil {
		return err
	}
	return nil
}
