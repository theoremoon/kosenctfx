package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/theoremoon/kosenctfx/scoreserver/deployment"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
)

const (
	SIMULTANEOUS_DEPLOYMENT_LIMIT = 4
)

func (s *server) getDeploymentHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		lc := c.(*loginContext)
		taskID, err := strconv.ParseUint(c.Param("task_id"), 10, 32)
		if err != nil {
			return errorHandle(c, err)
		}
		hostAndPort, err := s.app.GetHostAndPortByTeamAndTaskID(lc.Team.ID, uint32(taskID))
		if err != nil {
			return errorHandle(c, err)
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"host": hostAndPort.Host,
			"port": hostAndPort.Port,
		})
	}
}

func (s *server) requestDeployHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		lc := c.(*loginContext)
		req := new(struct {
			TaskID uint32 `json:"task_id"`
		})
		if err := c.Bind(req); err != nil {
			return errorHandle(c, err)
		}

		// check the task is opened and deployable
		chal, err := s.app.GetRawChallengeByID(req.TaskID)
		if err != nil || !chal.IsOpen {
			return errorMessageHandle(c, http.StatusBadRequest, NoSuchTaskMessage)
		}
		if chal.Deployment != deployment.TYPE_EACH {
			// message さぼってる
			return errorMessageHandle(c, http.StatusBadRequest, InvalidRequestMessage)
		}

		deployments, err := s.app.ListLivingDeploymentsByTeamID(lc.Team.ID)
		if err != nil {
			return errorHandle(c, err)
		}
		// rate limit check
		if len(deployments) >= SIMULTANEOUS_DEPLOYMENT_LIMIT {
			return errorMessageHandle(c, http.StatusTooManyRequests, TooManyDeploymentsMessage)
		}
		now := time.Now()
		for _, d := range deployments {
			if now.Sub(time.Unix(d.UpdatedAt, 0)) < 1*time.Second {
				return errorMessageHandle(c, http.StatusTooManyRequests, TooFrequentlyDeploymentMessage)
			}
		}

		// duplicate check（同じタスクを2個デプロイさせない）
		for _, d := range deployments {
			if d.ChallengeId == chal.ID {
				// message さぼってる
				return errorMessageHandle(c, http.StatusBadRequest, InvalidRequestMessage)
			}
		}

		// なんかエージェントに割り当てるアルゴリズム
		agent, err := s.app.GetAgentForRequestDeployment()
		if err != nil {
			return errorHandle(c, err)
		}

		// do request
		tid := lc.Team.ID
		_, err = s.app.RequestDeploy(agent, chal, &tid)
		if err != nil {
			return errorHandle(c, err)
		}

		// message ID振るのサボってる
		return messageHandle(c, "Your request is queued. Please wait...")
	}
}

func (s *server) requestRetireHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		lc := c.(*loginContext)
		req := new(struct {
			TaskID uint32 `json:"task_id"`
		})
		if err := c.Bind(req); err != nil {
			return errorHandle(c, err)
		}

		// check the task is opened and deployable
		chal, err := s.app.GetRawChallengeByID(req.TaskID)
		if err != nil || !chal.IsOpen {
			return errorMessageHandle(c, http.StatusBadRequest, NoSuchTaskMessage)
		}
		if chal.Deployment != deployment.TYPE_EACH {
			// message さぼってる
			return errorMessageHandle(c, http.StatusBadRequest, InvalidRequestMessage)
		}

		deployments, err := s.app.ListLivingDeploymentsByTeamID(lc.Team.ID)
		if err != nil {
			return errorHandle(c, err)
		}
		// rate limit check
		now := time.Now()
		for _, d := range deployments {
			if now.Sub(time.Unix(d.UpdatedAt, 0)) < 1*time.Second {
				return errorMessageHandle(c, http.StatusTooManyRequests, TooFrequentlyDeploymentMessage)
			}
		}

		// タスクが存在することを確認したらretire requestを詰む
		var targetDeployment *model.Deployment = nil
		for _, d := range deployments {
			if d.ChallengeId == chal.ID {
				targetDeployment = d
				break
			}
		}
		if targetDeployment == nil {
			return errorMessageHandle(c, http.StatusBadRequest, InvalidRequestMessage)
		}

		// do request
		err = s.app.RequestRetire(targetDeployment)
		if err != nil {
			return errorHandle(c, err)
		}

		// message ID振るのサボってる
		return messageHandle(c, "Your request is queued. Please wait...")
	}
}

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
				RequestedAt: d.RequestedAt,
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
		if err := c.Bind(req); err != nil {
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
