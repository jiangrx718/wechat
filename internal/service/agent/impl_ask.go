package agent

import (
	"context"
	"errors"
	"wechat-tools/internal/common"
	"wechat-tools/utils"

	"github.com/sashabaranov/go-openai"
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

	// 校验 client 是否初始化成功
	if s.client == nil || s.client.Client == nil {
		return result, errors.New("agent client not initialized")
	}

	// 通过 s.client 调用 openai
	resp, err := s.client.Client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.client.Model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: question},
		},
	})
	if err != nil {
		return result, err
	}

	result.Data = AgentAskResp{
		Answer: resp.Choices[0].Message.Content,
	}

	return result, nil
}
