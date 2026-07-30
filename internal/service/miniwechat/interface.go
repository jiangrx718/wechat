package miniwechat

import (
	"context"
	"wechat-tools/internal/common"
)

type WechatUserServiceIFace interface {
	// Exist 是否存在
	Exist(ctx context.Context, deviceId string) (common.ServiceResult, error)
	// Create 创建
	Create(ctx context.Context, userName, deviceId string, score int) (common.ServiceResult, error)
	// Update 更新
	Update(ctx context.Context, userName string, score int) (common.ServiceResult, error)
	// List 列表
	List(ctx context.Context) (common.ServiceResult, error)
}

type CheckImageServiceIFace interface {
	// Check 对图片做内容安全检测（调用微信 security.imgSecCheck）。
	// media 为图片的原始字节，filename 为文件名（含扩展名，供微信识别类型）。
	// 返回 ServiceResult：code=0 时 data 为 *CheckResult{Pass bool}。
	Check(ctx context.Context, media []byte, filename string) (common.ServiceResult, error)
}

// CheckResult 检测结果，序列化后形如 {"pass":true}
type CheckResult struct {
	Pass bool `json:"pass"`
}
