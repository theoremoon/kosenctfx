package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/theoremoon/kosenctfx/scoreserver/agent"
	"github.com/theoremoon/kosenctfx/scoreserver/client"
)

func mainLoop(client *client.Client) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	hostname, err := agent.GetHostname()
	if err != nil {
		log.Printf("%+v\n", err)
		return
	}

	t := time.NewTicker(5 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			_, err := client.Beat(ctx, hostname)
			if err != nil {
				log.Printf("%+v\n", err)
				break
			}

			// 結局ここで今管理しているtaskを管理する必要がありそう

		case <-ctx.Done():
			return
		}
	}
}
