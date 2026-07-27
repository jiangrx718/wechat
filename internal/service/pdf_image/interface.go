package pdf_image

import (
	"context"

	"wechat-tools/internal/common"
)

// ServiceIFace PDF 转图片服务接口
type ServiceIFace interface {
	// Convert 将 PDF 字节数据转为图片，返回 base64 编码的图片列表。
	// pdfData 为 PDF 原始字节，filename 为原始文件名（仅用于日志）。
	// opts 可指定 Format 和 DPI 等参数。
	Convert(ctx context.Context, pdfData []byte, filename string, opts ...ConvertOption) (common.ServiceResult, error)
}

// ConvertParams 转换参数
type ConvertParams struct {
	Format string // 图片格式: png / jpeg，默认 png
	DPI    int    // 输出分辨率，默认 150
}

// ConvertOption 转换参数选项
type ConvertOption func(*ConvertParams)

// WithFormat 设置输出图片格式
func WithFormat(format string) ConvertOption {
	return func(p *ConvertParams) {
		p.Format = format
	}
}

// WithDPI 设置输出分辨率（DPI）
func WithDPI(dpi int) ConvertOption {
	return func(p *ConvertParams) {
		p.DPI = dpi
	}
}

// PageImage 单页图片结果
type PageImage struct {
	Page int    `json:"page"`
	Data string `json:"data"` // base64 编码的图片数据
	Mime string `json:"mime"` // 图片 MIME 类型
}

// ConvertResult 转换结果
type ConvertResult struct {
	TotalPages int         `json:"total_pages"`
	Images     []PageImage `json:"images"`
}
