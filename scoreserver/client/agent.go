package client

import (
	"context"
)

type BeatResult struct {
}

func (c *Client) Beat(ctx context.Context, Hostname string) (*BeatResult, error) {
	res, err := c.Post(ctx, "/agent/beat", map[string]interface{}{
		"agent_id": Hostname,
	})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result BeatResult
	if err := DecodeBody(res, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
