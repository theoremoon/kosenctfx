package main

import (
	"log"
	"net/url"
	"os"

	kosenctfxClient "github.com/theoremoon/kosenctfx/scoreserver/client"
	cli "github.com/urfave/cli/v2"
)

var (
	client *kosenctfxClient.Client
)

func main() {
	app := &cli.App{
		Name: "kosenctfx cli",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "url",
				Aliases: []string{"u"},
			},
			&cli.StringFlag{
				Name: "api-key",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "challenge",
				Usage: "build challenge image and apply to scoreserver",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "dir",
						Aliases: []string{"d"},
					},
				},
				Action: func(c *cli.Context) error {
					err := challengeMain(c.String("dir"))
					if err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:  "config",
				Usage: "set the server configuration",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config-file",
						Aliases: []string{"f"},
						Usage:   "kosenctfx-conf.yml",
					},
				},
				Action: func(c *cli.Context) error {
					err := configMain(c.String("config-file"))
					if err != nil {
						return err
					}
					return nil
				},
			},
		},
		Before: func(c *cli.Context) error {
			u, err := url.Parse(c.String("url"))
			if err != nil {
				return err
			}
			client = kosenctfxClient.NewClient(u, c.String("api-key"))
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
