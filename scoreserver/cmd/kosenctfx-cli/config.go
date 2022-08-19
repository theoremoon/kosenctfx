package main

import (
	"context"

	"github.com/theoremoon/kosenctfx/scoreserver/config"
)

func configMain(path string) error {
	conf, err := config.LoadConfigFile(path)
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = client.SetRegistryConf(ctx, &conf.Registry)
	if err != nil {
		return err
	}

	return nil
}
