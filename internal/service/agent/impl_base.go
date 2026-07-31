package agent

import (
	"errors"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type AgentClient struct {
	Client *openai.Client
	Model  string
}

func NewAgentClient() (*AgentClient, error) {
	// 获取环境变量的KEY/BASE_URL/MODEL
	agentApiKey := strings.TrimSpace(os.Getenv("AGENT_API_KEY"))
	agentBaseURL := strings.TrimSpace(os.Getenv("AGENT_BASE_URL"))
	agentModel := strings.TrimSpace(os.Getenv("AGENT_MODEL"))
	if agentApiKey == "" || agentBaseURL == "" || agentModel == "" {
		return nil, errors.New("AGENT_API_KEY,AGENT_BASE_URL,AGENT_MODEL is require")
	}

	// 创建对应的默认配置
	cfg := openai.DefaultConfig(agentApiKey)
	cfg.BaseURL = agentBaseURL

	// 初始对应的AgentClient结构体
	var agentClient AgentClient
	agentClient.Client = openai.NewClientWithConfig(cfg)
	agentClient.Model = agentModel

	return &agentClient, nil
}
