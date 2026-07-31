package agent

import (
	"wechat-tools/internal/dao"
	"wechat-tools/utils"

	"gorm.io/gorm"
)

type AgentService struct {
	db     *gorm.DB
	client *AgentClient
}

func NewAgentService() *AgentService {
	s := &AgentService{db: utils.DB()}
	dao.SetDefault(utils.DB())

	// 初始化 openai client
	if agentClient, err := NewAgentClient(); err == nil {
		s.client = agentClient
	}
	return s
}
