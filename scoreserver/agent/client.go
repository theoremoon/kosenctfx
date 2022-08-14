package agent

import (
	"context"

	"github.com/theoremoon/kosenctfx/scoreserver/client"
)

type Client struct {
	*client.Client
	Hostname string
	Hostaddr string
}

type BeatResult struct {
}

func (c *Client) Beat(ctx context.Context) (*BeatResult, error) {
	res, err := c.Post(ctx, "/agent/beat", map[string]interface{}{
		"agent_id": c.Hostname,
		"hostaddr": c.Hostaddr,
	})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result BeatResult
	if err := client.DecodeBody(res, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
