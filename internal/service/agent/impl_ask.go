package agent

import (
	"context"
	"wechat-tools/internal/common"
	"wechat-tools/utils"
)

type AgentAskResp struct {
	Answer string `json:"answer"`
}

func (s *AgentService) Ask(ctx context.Context, question string) (common.ServiceResult, error) {
	var (
		logger = utils.SugarContext(ctx)
		result = common.NewServiceResult()
	)
	logger.Info("问答开始")

	result.Data = AgentAskResp{}

	return result, nil
}
