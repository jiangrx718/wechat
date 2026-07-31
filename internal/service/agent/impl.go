package agent

import (
	"wechat-tools/internal/dao"
	"wechat-tools/utils"

	"gorm.io/gorm"
)

type AgentService struct {
	db *gorm.DB
}

func NewAgentService() *AgentService {
	s := &AgentService{db: utils.DB()}
	dao.SetDefault(utils.DB())
	return s
}
