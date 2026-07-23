package check_image

import (
	"context"

	"wechat-tools/internal/common"
)

type ServiceIFace interface {
	// Check 对图片做内容安全检测（调用微信 security.imgSecCheck）。
	// media 为图片的原始字节，filename 为文件名（含扩展名，供微信识别类型）。
	// 返回 ServiceResult：code=0 时 data 为 *CheckResult{Pass bool}。
	Check(ctx context.Context, media []byte, filename string) (common.ServiceResult, error)
}

// CheckResult 检测结果，序列化后形如 {"pass":true}
type CheckResult struct {
	Pass bool `json:"pass"`
}
