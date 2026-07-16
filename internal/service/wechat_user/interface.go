package wechat_user

import (
	"context"
	"wechat-tools/internal/common"
)

type ServiceIFace interface {
	// Exist 是否存在
	Exist(ctx context.Context, deviceId string) (common.ServiceResult, error)
	// Create 创建
	Create(ctx context.Context, userName, deviceId string, score int) (common.ServiceResult, error)
	// Update 更新
	Update(ctx context.Context, userName string, score int) (common.ServiceResult, error)
	// List 列表
	List(ctx context.Context) (common.ServiceResult, error)
}
