package service

import (
	"github.com/theoremoon/kosenctfx/scoreserver/deployment"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
)

type DeploymentApp interface {
	ListLivingDeployments() ([]*model.Deployment, error)
	ListDeploymentRequestForAgent(agentID string) ([]*model.Deployment, error)
	StartDeployment(deploymentID uint32, port int) error
	UpdateDeploymentStatus(deploymentID uint32, status string) error
}

func (app *app) ListLivingDeployments() ([]*model.Deployment, error) {
	var deployments []*model.Deployment
	livingStatuses := []string{deployment.STATUS_WAITING, deployment.STATUS_DEPLOYING, deployment.STATUS_AVAILABLE}
	err := app.db.
		Where("status IN ?", livingStatuses).
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
		Where("status = ? AND agent_id = ?", "waiting", agentID).
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
