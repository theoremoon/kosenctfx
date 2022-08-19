package main

import (
	"log"
	"net/url"
	"os"

	"github.com/theoremoon/kosenctfx/scoreserver/client"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "kosenctfx-agent",
		Usage: "kosenctfx agent",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "url",
				Aliases: []string{"u"},
			},
			&cli.StringFlag{
				Name: "api-key",
			},
		},
		Action: func(c *cli.Context) error {
			u, err := url.Parse(c.String("url"))
			if err != nil {
				return err
			}

			client := client.NewClient(u, c.String("api-key"))

			// registryConf, err := client.GetRegistryConfig(context.Background())
			// if err != nil {
			// 	return nil
			// }

			mainLoop(client)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
