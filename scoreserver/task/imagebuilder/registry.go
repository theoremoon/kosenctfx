package imagebuilder

import (
	"encoding/base64"
	"encoding/json"

	configtypes "github.com/docker/cli/cli/config/types"
	mobytypes "github.com/moby/moby/api/types"
)

type RegistryConfig struct {
	URL      string `yaml:"url" json:"url"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

func (c *RegistryConfig) GetRegistry() string {
	return c.URL
}

func (c *RegistryConfig) GetRegistryAuth() string {
	aconf := mobytypes.AuthConfig{
		Username: c.Username,
		Password: c.Password,
	}

	blob, err := json.Marshal(aconf)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(blob)
}

func (c *RegistryConfig) GetAuthConfig() configtypes.AuthConfig {
	return configtypes.AuthConfig{
		Username: c.Username,
		Password: c.Password,
	}
}
