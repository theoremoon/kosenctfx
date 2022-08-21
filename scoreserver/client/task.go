package client

import (
	"context"

	"github.com/theoremoon/kosenctfx/scoreserver/task"
)

type GetPresignedURLResult struct {
	PresignedURL string `json:"presignedURL"`
	DownloadURL  string `json:"downloadURL"`
}

func (c *Client) GetPresignedURL(ctx context.Context, filename string) (*GetPresignedURLResult, error) {
	res, err := c.Post(ctx, "/admin/get-presigned-url", map[string]interface{}{
		"key": filename,
	})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result GetPresignedURLResult
	if err := DecodeBody(res, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) NewChallenge(ctx context.Context, taskdef *task.TaskDefinition) error {
	res, err := c.Post(ctx, "/admin/new-challenge", map[string]interface{}{
		"name":        taskdef.Name,
		"flag":        taskdef.Flag,
		"category":    taskdef.Category,
		"description": taskdef.Description,
		"author":      taskdef.Author,
		"is_survey":   taskdef.IsSurvey,
		"tags":        taskdef.Tags,
		"attachments": taskdef.Attachments,
		"compose":     taskdef.Compose,
		"deployment":  taskdef.Deployment,
		"lifespan":    taskdef.Lifespan,
	})
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
