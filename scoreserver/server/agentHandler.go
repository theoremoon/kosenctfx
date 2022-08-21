package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/theoremoon/kosenctfx/scoreserver/agent/order"
)

// agentMiddleware
func (s *server) agentBeatHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		// receive agent's heartbeat and save last activity time and public address
		req := new(struct {
			AgentID string `json:"agent_id"`
		})
		if err := c.Bind(req); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		// get public ip from request
		// ここは構成によって変わるところ
		ip := c.RealIP()
		err := s.app.AgentHeartbeat(req.AgentID, ip)
		if err != nil {
			return errorHandle(c, err)
		}

		// agentが実行するべきorderを取得してくる
		deployments, err := s.app.ListDeploymentRequestForAgent(req.AgentID)
		if err != nil {
			return errorHandle(c, err)
		}

		return c.JSON(http.StatusOK, order.Order{
			Deployments: deployments,
		})
	}
}

// agentMiddleware
func (s *server) agentStartDeploymentHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(struct {
			DeploymentID uint32 `json:"deployment_id"`
			TaskID       uint32 `json:"task_id"`
			Port         int    `json:"port"`
		})
		if err := c.Bind(req); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		err := s.app.StartDeployment(req.DeploymentID, req.Port)
		if err != nil {
			return errorHandle(c, err)
		}

		task, err := s.app.GetRawChallengeByID(req.TaskID)
		if err != nil {
			return errorHandle(c, err)
		}

		registry, err := s.app.GetRegistryConfig()
		if err != nil {
			return errorHandle(c, err)
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"compose":  task.Compose,
			"registry": registry,
		})
	}
}

// adminMiddleware
func (s *server) listAgentsHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		agents, err := s.app.ListAvailableAgents()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, agents)
	}
}

// agentMiddleware
func (s *server) agentUpdateDeploymentStatusHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(struct {
			DeploymentID uint32 `json:"deployment_id"`
			Status       string `json:"status"`
		})
		if err := c.Bind(req); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		err := s.app.UpdateDeploymentStatus(req.DeploymentID, req.Status)
		if err != nil {
			return errorHandle(c, err)
		}

		return c.NoContent(http.StatusOK)
	}
}
