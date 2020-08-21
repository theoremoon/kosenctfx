package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"golang.org/x/xerrors"
)

const (
	WORKDIR = "/tmp"
)

func pathes(id uint) (string, string) {
	return filepath.Join(WORKDIR, fmt.Sprintf("docker-compose_%d.yml", id)), filepath.Join(WORKDIR, fmt.Sprintf("docker-compose_%d_solve.yml", id))
}

type Challenge struct {
	Name             string
	ComposePath      string
	SolveComposePath string
	Flag             string
	Host             string
	Port             string
	IsRunning        bool
}

type checker struct {
	Challenges map[uint]*Challenge
	Token      string // Scoreserverにアクセスするためのトークン
	ServerURL  string // ScoreserverのURL
	DockerHost string // remote dockerのURL (tcp://x.y.z.w)

	sync.RWMutex
}

/// 問題チェックの結果を送るわよ
/// 送るだけ
func (c *checker) SendResult(id uint, result bool) {
	j, _ := json.Marshal(struct {
		ID     uint `json:"id"`
		Result bool `json:"result"`
	}{
		ID:     id,
		Result: result,
	})
	req, err := http.NewRequest(http.MethodPost, c.ServerURL+"/admin/set-challenge-status", bytes.NewBuffer(j))
	if err != nil {
		log.Printf("%s\n", c.ServerURL+"/admin/set-challenge-status")
		log.Printf("ERROR %v\n", xerrors.Errorf(": %w", err))
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("ERROR %v\n", xerrors.Errorf(": %w", err))
		return
	}
}

func (c *checker) Check(duration time.Duration) {
	for {
		c.RLock()
		// 各問題について並列で solve scriptを回す
		wg := sync.WaitGroup{}
		for id, chal := range c.Challenges {
			// solve scriptがない場合は何もしない
			if chal.SolveComposePath == "" {
				continue
			}

			go func(id uint, chal *Challenge) {
				wg.Add(1)

				cmd := exec.Command(
					"docker-compose",
					"-H", c.DockerHost,
					"-f",
					chal.SolveComposePath,
					"run",
					"--rm",
					"-e",
					"HOST="+chal.Host,
					"-e",
					"PORT="+chal.Port,
					"solve",
				)

				output, err := cmd.Output()
				if err != nil {
					log.Printf("ERROR: %+v\n", xerrors.Errorf(": %w", err))
					log.Println(string(output))
					// 一応解けないで送っておく
					c.SendResult(id, false)
				} else {
					if strings.Contains(string(output), chal.Flag) {
						// 解けた
						c.SendResult(id, true)

					} else {
						// 解けない
						c.SendResult(id, false)
					}
				}
				wg.Done()
			}(id, chal)
		}

		wg.Wait()
		c.RUnlock()

		// 一定期間休憩
		time.Sleep(duration)
	}
}

func run() error {
	dockerHost := flag.String("H", "", "tcp://x.y.z.w")
	token_ := flag.String("token", "", "token for this challenge manager")

	flag.Parse()
	token := *token_

	// checkerを走らせる
	checker := &checker{
		Challenges: make(map[uint]*Challenge),
		DockerHost: *dockerHost,
	}
	go func() {
		checker.Check(5 * time.Minute)
	}()

	// コマンドはHTTP経由で受け付ける（楽なので）
	server := echo.New()
	server.Use(middleware.Logger())
	server.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == token, nil
		},
	}))
	server.POST("/init", func(c echo.Context) error {
		req := new(struct {
			ServerToken string `json:"server_token"`
			ServerURL   string `json:"server_url"`
		})
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		checker.ServerURL = req.ServerURL
		checker.Token = req.ServerToken
		return c.NoContent(http.StatusOK)
	})
	server.POST("/register", func(c echo.Context) error {
		req := new(struct {
			ID           uint   `json:"id"`
			Compose      string `json:"compose"`
			SolveCompose string `json:"solve_compose"`
			Flag         string `json:"flag"`
			Host         string `json:"host"`
			Port         string `json:"port"`
		})
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		composePath, solveComposePath := pathes(req.ID)
		if err := ioutil.WriteFile(composePath, []byte(req.Compose), 0644); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		if req.SolveCompose == "" {
			solveComposePath = ""
		} else {
			if err := ioutil.WriteFile(solveComposePath, []byte(req.SolveCompose), 0644); err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
		}

		checker.Lock()
		checker.Challenges[req.ID] = &Challenge{
			ComposePath:      composePath,
			SolveComposePath: solveComposePath,
			IsRunning:        false,
			Flag:             req.Flag,
			Host:             req.Host,
			Port:             req.Port,
		}
		checker.Unlock()

		return c.NoContent(http.StatusOK)
	})

	server.POST("/start", func(c echo.Context) error {
		req := new(struct {
			ID uint `json:"id"`
		})
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		if _, ok := checker.Challenges[req.ID]; !ok {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("challenge not registered: %d", req.ID))
		}

		// swarmに接続してChallengeの起動を依頼する
		// この辺はシェルを呼んだほうが早い
		cmd := exec.Command(
			"docker",
			"-H", checker.DockerHost,
			"stack",
			"deploy",
			"--compose-file",
			checker.Challenges[req.ID].ComposePath,
			"--with-registry-auth",
			"--resolve-image=always",
			fmt.Sprintf("%d", req.ID),
		)
		if output, err := cmd.Output(); err != nil {
			fmt.Println(string(output))
			return c.JSON(http.StatusBadRequest, "failed to deploy with error: "+err.Error())
		}

		// チェックするようにする
		checker.Lock()
		checker.Challenges[req.ID].IsRunning = false
		checker.Unlock()
		return c.NoContent(http.StatusOK)
	})

	server.POST("/stop", func(c echo.Context) error {
		req := new(struct {
			ID uint `json:"id"`
		})
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		checker.Lock()
		if _, ok := checker.Challenges[req.ID]; ok {
			checker.Challenges[req.ID].IsRunning = false
		}
		checker.Unlock()
		return c.NoContent(http.StatusOK)
	})

	if err := server.Start(":5000"); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v\n", err)
	}
}
