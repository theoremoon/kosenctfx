package client

import (
	"context"

	"github.com/theoremoon/kosenctfx/scoreserver/task/registry"
)

type GetRegistryConfResult struct {
	URL      string `json:"url"`
	Username string `json:"user"`
	Password string `json:"password"`
}

func (c *Client) GetRegistryConf(ctx context.Context) (*GetRegistryConfResult, error) {
	res, err := c.Get(ctx, "/admin/registry-conf", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result GetRegistryConfResult
	if err := DecodeBody(res, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SetRegistryConf(ctx context.Context, registry *registry.RegistryConfig) error {
	res, err := c.Post(ctx, "/admin/set-registry-conf", map[string]interface{}{
		"url":      registry.URL,
		"user":     registry.Username,
		"password": registry.Password,
	})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
