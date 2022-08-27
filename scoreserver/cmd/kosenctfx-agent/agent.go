package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/theoremoon/kosenctfx/scoreserver/agent"
	"github.com/theoremoon/kosenctfx/scoreserver/client"
	"github.com/theoremoon/kosenctfx/scoreserver/deployment"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/task"
)

func mainLoop(client *client.Client, dir string) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	hostname, err := agent.GetHostname()
	if err != nil {
		log.Printf("%+v\n", err)
		return
	}

	t := time.NewTicker(5 * time.Second)
	defer t.Stop()

	deployments := make(map[uint32]context.CancelFunc)
	wg := sync.WaitGroup{}

	for {
		select {
		case <-t.C:
			order, err := client.Beat(ctx, hostname)
			if err != nil {
				log.Printf("%+v\n", err)
				break
			}

			// do deploy
			for _, d := range order.Deployments {
				if _, exists := deployments[d.ID]; exists {
					continue
				}
				ctx2, cancel := context.WithCancel(ctx)
				deployments[d.ID] = cancel
				wg.Add(1)
				go func(d *model.Deployment) {
					err := doDeploy(ctx2, client, dir, d)
					delete(deployments, d.ID)
					if err != nil {
						client.UpdateDeploymentStatus(ctx2, d, deployment.STATUS_ERROR)
						log.Printf("%+v\n", err)
					}
					wg.Done()
				}(d)
			}

			// do retire
			for _, d := range order.Retires {
				cancel, exists := deployments[d.ID]
				if exists {
					// retire
					cancel()
				} else {
					// make error
					client.UpdateDeploymentStatus(ctx, d, deployment.STATUS_ERROR)
				}
			}

			// 結局ここで今管理しているtaskを管理する必要がありそう

		case <-ctx.Done():
			wg.Wait()
			return
		}
	}
}

func doDeploy(ctx context.Context, client *client.Client, dir string, d *model.Deployment) error {
	// fix port
	var err error
	port, err := assignPort()
	if err != nil {
		return err
	}

	// get task and update status
	res, err := client.StartDeployment(ctx, d, port)
	if err != nil {
		return err
	}

	// docker-compose up
	taskID := strconv.FormatUint(uint64(d.ID), 10)
	composePath := filepath.Join(dir, taskID, "docker-compose.yml")
	composeConfig, err := task.ParseComposeConfig(composePath, []byte(res.Compose))
	composeConfig.Name = taskID
	if err != nil {
		return err
	}

	compose, err := task.NewCompose(res.Registry)
	if err != nil {
		return err
	}

	// timeout
	var ctx2 context.Context
	var cancel context.CancelFunc
	if d.RetiresAt > 0 {
		ctx2, cancel = context.WithDeadline(ctx, time.Unix(d.RetiresAt, 0))
	} else {
		ctx2, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	// tell status
	err = client.UpdateDeploymentStatus(ctx, d, deployment.STATUS_AVAILABLE)
	if err != nil {
		log.Println(err)
	}

	// ctx2 の cancelでautomaticallyにstopされる
	go compose.Up(ctx2, composeConfig, port)

	// wait
	select {
	case <-ctx2.Done():
	}
	log.Println("shutting down...")

	// tell staus
	err = client.UpdateDeploymentStatus(ctx, d, deployment.STATUS_RETIRED)
	if err != nil {
		log.Println(err)
	}

	// cleanup
	compose.Down(context.Background(), taskID)
	return nil
}

func assignPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
