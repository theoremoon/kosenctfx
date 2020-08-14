package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"golang.org/x/xerrors"
)

const (
	WORKDIR = "challenge_manager"
)

func pathes(name string) (string, string) {
	return filepath.Join(WORKDIR, fmt.Sprintf("docker_compose_%s.yml", name)), filepath.Join(WORKDIR, fmt.Sprintf("docker_compose_%s_solve.yml", name))
}

type Challenge struct {
	ComposePath      string
	SolveComposePath string
	Flag             string
	Host             string
	Port             string
	IsRunning        bool
}

type checker struct {
	Challenges map[string]*Challenge
	Token      string // Scoreserverにアクセスするためのトークン
	ServerURL  string // ScoreserverのURL

	sync.RWMutex
}

/// 問題チェックの結果を送るわよ
/// 送るだけ
func (c *checker) SendResult(name string, result bool) {
	j, _ := json.Marshal(struct {
		Name   string `json:"name"`
		Result bool   `json:"result"`
	}{
		Name:   name,
		Result: result,
	})
	_, _ = http.Post(c.ServerURL+"/check-result", "application/json", bytes.NewBuffer(j))
}

func (c *checker) Check(duration time.Duration) {
	for {
		c.RLock()
		// 各問題についてsolve scriptを回す
		for name, chal := range c.Challenges {
			cmd := exec.Command(
				"docker-compose",
				"--compose-file",
				chal.SolveComposePath,
				"run",
				"--rm",
				"-e",
				"HOST="+chal.Host,
				"-e",
				"PORT="+chal.Port,
			)

			output, err := cmd.Output()
			if err != nil {
				log.Printf("ERROR: %s\n", err.Error())
				// 一応解けないで送っておく
				c.SendResult(name, false)
			} else {
				if strings.Contains(string(output), chal.Flag) {
					// 解けた
					c.SendResult(name, true)

				} else {
					// 解けない
					c.SendResult(name, false)
				}
			}
		}
		c.RUnlock()

		// 一定期間休憩
		time.Sleep(duration)
	}
}

func run() error {
	// ChallengeManagerにアクセスするためのToken
	token := uuid.New().String()
	fmt.Sprintf("TOKEN: %s", token)

	// checkerを走らせる
	checker := &checker{
		Challenges: make(map[string]*Challenge),
	}
	go func() {
		checker.Check(5 * time.Minute)
	}()

	// コマンドはHTTP経由で受け付ける（楽なので）
	server := echo.New()
	server.Use(middleware.Logger())
	server.POST("/register", func(c echo.Context) error {
		req := new(struct {
			Name         string `json:"name"`
			Compose      string `json:"compose"`
			SolveCompose string `json:"solve_compose"`
			Flag         string `json:"flag"`
			Host         string `json:"host"`
			Port         string `json:"port"`
			Token        string `json:"token"`
		})
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if token != req.Token {
			return c.JSON(http.StatusBadRequest, "invalid token")
		}

		composePath, solveComposePath := pathes(req.Name)
		if err := ioutil.WriteFile(composePath, []byte(req.Compose), 0644); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		if err := ioutil.WriteFile(solveComposePath, []byte(req.SolveCompose), 0644); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		checker.Lock()
		checker.Challenges[req.Name] = &Challenge{
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
			Name  string `json:"name"`
			Token string `json:"token"`
		})
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if token != req.Token {
			return c.JSON(http.StatusBadRequest, "invalid token")
		}

		if _, ok := checker.Challenges[req.Name]; !ok {
			return c.JSON(http.StatusBadRequest, "challenge not registered: "+req.Name)
		}

		// swarmに接続してChallengeの起動を依頼する
		// この辺はシェルを呼んだほうが早い
		cmd := exec.Command(
			"docker",
			"stack",
			"deploy",
			"--compose-file",
			checker.Challenges[req.Name].ComposePath,
			req.Name,
		)
		if err := cmd.Run(); err != nil {
			return c.JSON(http.StatusBadRequest, "failed to deploy with error: "+err.Error())
		}

		// チェックするようにする
		checker.Lock()
		checker.Challenges[req.Name].IsRunning = false
		checker.Unlock()
		return c.NoContent(http.StatusOK)
	})

	server.POST("/stop", func(c echo.Context) error {
		req := new(struct {
			Name  string `json:"name"`
			Token string `json:"token"`
		})
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if token != req.Token {
			return c.JSON(http.StatusBadRequest, "invalid token")
		}

		checker.Lock()
		if _, ok := checker.Challenges[req.Name]; ok {
			checker.Challenges[req.Name].IsRunning = false
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
