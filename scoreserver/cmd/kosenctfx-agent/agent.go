package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/theoremoon/kosenctfx/scoreserver/agent"
)

func mainLoop(agent agent.Agent) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	t := time.NewTicker(5 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			_, err := agent.Client().Beat(ctx)
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
