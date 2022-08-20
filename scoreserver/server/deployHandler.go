package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
)

type livingDeployment struct {
	ID uint32 `json:"id"`

	Challenge *service.Challenge `json:"challenge"`
	Agent     *model.Agent       `json:"agent"`
	Port      int64              `json:"port"`
	Status    string             `json:"status"` // waiting, deploying, available, retired, error

	RequestedAt int64 `json:"requested_at"`
	RetiresAt   int64 `json:"retires_at"`
}

func (s *server) adminListLivingDeploymentsHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ds, err := s.app.ListLivingDeployments()
		if err != nil {
			return errorHandle(c, err)
		}

		challengeIDs := make([]uint32, 0)
		agentIDs := make([]string, 0)
		// teamIDs := make([]uint32, 0)
		for _, d := range ds {
			challengeIDs = append(challengeIDs, d.ChallengeId)
			agentIDs = append(agentIDs, d.AgentId)
			// tid := d.TeamId
			// if tid != nil {
			// 	teamIDs = append(teamIDs, *tid)
			// }
		}

		challenges, err := s.app.ListChallengeByIDs(challengeIDs)
		challengeMap := make(map[uint32]*service.Challenge)
		if err != nil {
			return errorHandle(c, err)
		}
		for _, chal := range challenges {
			challengeMap[chal.ID] = chal
		}

		agents, err := s.app.ListAgentsByIDs(agentIDs)
		agentMap := make(map[string]*model.Agent)
		if err != nil {
			return errorHandle(c, err)
		}
		for _, a := range agents {
			agentMap[a.AgentID] = a
		}

		deployments := make([]*livingDeployment, 0, len(ds))
		for _, d := range ds {
			var chal *service.Challenge = nil
			if ch, exists := challengeMap[d.ChallengeId]; exists {
				chal = ch
			}
			var agent *model.Agent = nil
			if a, exists := agentMap[d.AgentId]; exists {
				agent = a
			}

			deployments = append(deployments, &livingDeployment{
				ID:          d.ID,
				Agent:       agent,
				Challenge:   chal,
				Port:        d.Port,
				Status:      d.Status,
				RequestedAt: d.RetiresAt,
				RetiresAt:   d.RetiresAt,
			})
		}
		return c.JSON(http.StatusOK, deployments)
	}
}

// agentMiddleware
func (s *server) adminRequestDeployHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		// receive agent's heartbeat and save last activity time and public address
		req := new(struct {
			TaskID  uint32 `json:"task_id"`
			AgentID string `json:"agent_id"`
		})
		if err := c.Bind(&req); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		agent, err := s.app.GetAgentByID(req.AgentID)
		if err != nil {
			return errorHandle(c, err)
		}

		task, err := s.app.GetRawChallengeByID(req.TaskID)
		if err != nil {
			return errorHandle(c, err)
		}

		// team id の指定はいまnilということに
		_, err = s.app.RequestDeploy(agent, task, nil)
		if err != nil {
			return errorHandle(c, err)
		}

		return c.JSON(http.StatusOK, nil)
	}
}
