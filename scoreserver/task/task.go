package task

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Attachment struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type TaskDefinition struct {
	ID          string
	Name        string
	Description string
	Flag        string
	Author      string
	Category    string
	Tags        []string
	Attachments []Attachment
	IsSurvey    bool `yaml:"is_survey" json:"is_survey"`
	Compose     string
	Deployment  string
	Lifespan    int // インスタンスが起動してから自動でシャットダウンされるまでの時間（秒） / 0なら無限
}

func LoadTaskDefinition(path string) (*TaskDefinition, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var taskdef TaskDefinition
	if err := yaml.Unmarshal(buf, &taskdef); err != nil {
		return nil, err
	}

	return &taskdef, nil
}
