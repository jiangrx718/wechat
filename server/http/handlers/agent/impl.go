package agent

import (
	agentService "wechat-tools/internal/service/agent"

	"github.com/gin-gonic/gin"
)

// NewAgentHandler 创建知识库处理器
func NewAgentHandler() *AgentHandler {
	return &AgentHandler{
		service: agentService.NewAgentService(),
	}
}

// AgentHandler 知识库处理器：文件入库 + 语义检索
type AgentHandler struct {
	service agentService.AgentServiceIFace
}

// RegisterRoutes 注册路由
func (h *AgentHandler) RegisterRoutes(routerGroup *gin.RouterGroup) {
	g := routerGroup.Group("/agent")

	// POST /api/agent/documents/ingest 单文件入库
	g.POST("/documents/ingest", h.Ingest)
	// POST /api/agent/search 语义检索
	g.POST("/search", h.Search)
}
