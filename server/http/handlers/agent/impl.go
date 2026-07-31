package agent

import (
	agentService "wechat-tools/internal/service/agent"

	"github.com/gin-gonic/gin"
)

func NewAgentHandler(engine *gin.Engine) *AgentHandler {
	return &AgentHandler{
		engine:  engine,
		service: agentService.NewAgentService(),
	}
}

type AgentHandler struct {
	engine  *gin.Engine
	service agentService.AgentServiceIFace
}

func (h *AgentHandler) RegisterRoutes(routerGroup *gin.RouterGroup) {
	g := routerGroup.Group("/agent")

	g.POST("/Ask", h.Ask)
}
