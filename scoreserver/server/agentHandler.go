package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *server) agentHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		// receive agent's heartbeat and save last activity time and public address
		req := new(struct {
			AgentID string `json:"agent_id"`
		})
		if err := c.Bind(&req); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		// get public ip from request
		// ここは構成によって変わるところ
		ip := c.RealIP()

		if err := s.app.AgentHeartbeat(req.AgentID, ip); err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return nil
	}
}
