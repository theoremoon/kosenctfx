package config

import (
	"os"

	yaml "github.com/goccy/go-yaml"
	"github.com/theoremoon/kosenctfx/scoreserver/task/registry"
)

type ConfigDefinition struct {
	Registry registry.RegistryConfig `yaml:"registry"`
}

func LoadConfigFile(path string) (*ConfigDefinition, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf ConfigDefinition
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
