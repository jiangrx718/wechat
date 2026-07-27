package pdf_image

import (
	"context"
	"strings"

	"wechat-tools/internal/common"
)

// ConvertParams 转换参数
type ConvertParams struct {
	Format          string // 图片格式: png / jpeg，默认 png
	DPI             int    // 输出分辨率，默认 150
	Page            int    // 指定转换的页码（从 1 开始），0 表示转换整个 PDF 的全部页
	Quality         string // 清晰度档位: low / medium / high / original，为空时使用 DPI
	RemoveWatermark bool   // 渲染前是否尝试移除 PDF 中的水印（基于结构化移除，识别不到则原样保留）
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

// WithPage 设置只转换指定页码（从 1 开始）。传入 0 表示转换整个 PDF 的全部页。
func WithPage(page int) ConvertOption {
	return func(p *ConvertParams) {
		p.Page = page
	}
}

// WithQuality 设置清晰度档位: low / medium / high / original。
// 设置后会覆盖 DPI，分别映射为 96 / 150 / 300 / 600。
func WithQuality(quality string) ConvertOption {
	return func(p *ConvertParams) {
		p.Quality = quality
	}
}

// WithRemoveWatermark 设置是否在渲染前尝试移除 PDF 中的水印。
// 采用基于 PDF 结构的精准移除（pdfcpu），识别不到水印时原样保留 PDF，不会误伤正文。
func WithRemoveWatermark(remove bool) ConvertOption {
	return func(p *ConvertParams) {
		p.RemoveWatermark = remove
	}
}

// dpiFromQuality 将清晰度档位映射为 DPI，未识别的档位返回 0（表示沿用 DPI 参数）
func dpiFromQuality(quality string) int {
	switch strings.ToLower(quality) {
	case "low":
		return 96
	case "medium":
		return 150
	case "high":
		return 300
	case "original":
		return 600
	default:
		return 0
	}
}

// PageImage 单页图片结果
type PageImage struct {
	Page int    `json:"page"`
	Data []byte `json:"-"`    // 图片二进制数据，不参与 JSON 序列化（由 handler 拼成 data URI 返回）
	Mime string `json:"mime"` // 图片 MIME 类型
}

// ConvertResult 转换结果
type ConvertResult struct {
	TotalPages int         `json:"total_pages"`
	Images     []PageImage `json:"images"`
}

// ServiceIFace PDF 转图片服务接口
type ServiceIFace interface {
	// Convert 将 PDF 字节数据转为图片。全程使用临时目录，转换完成后读入内存并立即清理，
	// 不在服务端持久化任何文件。返回结果中 PageImage.Data 为每张图片的二进制数据。
	// pdfData 为 PDF 原始字节，filename 为原始文件名（仅用于日志）。
	// opts 可指定 Format、DPI、Page、Quality、RemoveWatermark 等参数。
	Convert(ctx context.Context, pdfData []byte, filename string, opts ...ConvertOption) (common.ServiceResult, error)
}
