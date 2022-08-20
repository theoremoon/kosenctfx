package service

import (
	"github.com/theoremoon/kosenctfx/scoreserver/deployment"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
)

type DeploymentApp interface {
	ListLivingDeployments() ([]*model.Deployment, error)
}

func (app *app) ListLivingDeployments() ([]*model.Deployment, error) {
	var deployments []*model.Deployment
	livingStatuses := []string{deployment.STATUS_WAITING, deployment.STATUS_DEPLOYING, deployment.STATUS_AVAILABLE}
	if err := app.db.Where("status IN ?", livingStatuses).Find(&deployments).Error; err != nil {
		return nil, err
	}
	return deployments, nil
}
