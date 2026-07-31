package agent

import (
	"context"
	"wechat-tools/internal/common"
)

type AgentServiceIFace interface {
	Ask(ctx context.Context, question string) (common.ServiceResult, error)
}
