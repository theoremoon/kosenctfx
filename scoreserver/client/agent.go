package client

import (
	"context"

	"github.com/theoremoon/kosenctfx/scoreserver/agent/order"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/task/registry"
)

func (c *Client) Beat(ctx context.Context, Hostname string) (*order.Order, error) {
	res, err := c.Post(ctx, "/agent/beat", map[string]interface{}{
		"agent_id": Hostname,
	})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result order.Order
	if err := DecodeBody(res, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

type StartDeploymentResult struct {
	Compose  string                   `json:"compose"`
	Registry *registry.RegistryConfig `json:"registry"`
}

func (c *Client) StartDeployment(ctx context.Context, deployment *model.Deployment, port int) (*StartDeploymentResult, error) {
	res, err := c.Post(ctx, "/agent/start-deployment", map[string]interface{}{
		"deployment_id": deployment.ID,
		"task_id":       deployment.ChallengeId,
		"port":          port,
	})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result StartDeploymentResult
	if err := DecodeBody(res, &result); err != nil {
		return nil, err
	}
	return &result, nil

}

func (c *Client) UpdateDeploymentStatus(ctx context.Context, deployment *model.Deployment, status string) error {
	res, err := c.Post(ctx, "/agent/update-deployment-status", map[string]interface{}{
		"deployment_id": deployment.ID,
		"status":        status,
	})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
