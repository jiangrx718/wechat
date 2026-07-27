package pdf_image

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	pdfImageService "wechat-tools/internal/service/pdf_image"
	"wechat-tools/server/http/response"
	"wechat-tools/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Convert PDF 转图片接口。
//
// 请求: multipart/form-data
//   - file:   PDF 文件（必填）
//   - format: (可选) 图片格式 png/jpeg，默认 png
//   - dpi:    (可选) 输出分辨率 72-600，默认 150
//   - page:   (可选) 指定转换的页码（从 1 开始），0 或不传表示转换整个 PDF 的全部页
//
// 响应: {"code":0,"msg":"操作成功","data":{"total_pages":N,"images":[{"page":1,"url":"https://...","mime":"image/png"},...]}}
func (h *PdfImageHandler) Convert(ctx *gin.Context) {
	var logger = utils.SugarContext(ctx)

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		logger.Infow("Handler PdfImage Convert ctx.FormFile err", "error", err)
		response.ParameterError(ctx, err)
		return
	}

	// 限制 PDF 文件大小（50MB）
	const maxSize = 50 << 20 // 50MB
	if fileHeader.Size > maxSize {
		response.Failed(ctx, response.CodeParameterErr, "PDF 文件过大，请选择小于 50MB 的文件", nil)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		logger.Errorw("Handler PdfImage Convert open file err", "error", err)
		response.InternalError(ctx)
		return
	}
	defer file.Close()

	pdfData, err := io.ReadAll(file)
	if err != nil {
		logger.Errorw("Handler PdfImage Convert read file err", "error", err)
		response.InternalError(ctx)
		return
	}

	// 解析可选参数
	format := ctx.DefaultPostForm("format", "png")
	dpiStr := ctx.DefaultPostForm("dpi", "150")
	dpi, err := strconv.Atoi(dpiStr)
	if err != nil {
		dpi = 150
	}
	// page: 指定转换的页码（从 1 开始），0 或不传表示转换整个 PDF 的全部页
	pageStr := ctx.DefaultPostForm("page", "0")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		page = 0
	}

	// 创建唯一会话目录，转换后的图片文件会直接存入此目录
	sessionID := uuid.New().String()
	sessionDir := filepath.Join(h.staticDir, sessionID)

	result, err := h.service.Convert(ctx, pdfData, fileHeader.Filename,
		pdfImageService.WithFormat(format),
		pdfImageService.WithDPI(dpi),
		pdfImageService.WithPage(page),
		pdfImageService.WithOutputDir(sessionDir),
	)
	if err != nil {
		logger.Errorw("Handler PdfImage Convert service.Convert error", "error", err)
		h.cleanupSession(sessionDir)
		response.InternalError(ctx)
		return
	}

	if result.GetCode() != 0 {
		h.cleanupSession(sessionDir)
		response.Failed(ctx, result.GetCode(), result.GetMessage(), nil)
		return
	}

	// 构造响应：将 PageImage.Path 转为小程序可访问的完整 URL
	cr := result.GetData().(*pdfImageService.ConvertResult)
	resp := h.buildConvertResponse(ctx, sessionID, cr)

	response.Successful(ctx, resp)
}

// buildConvertResponse 将服务层结果转为带 URL 的 HTTP 响应体
func (h *PdfImageHandler) buildConvertResponse(ctx *gin.Context, sessionID string, cr *pdfImageService.ConvertResult) *convertResponse {
	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, ctx.Request.Host)

	images := make([]pageImageResponse, 0, len(cr.Images))
	for _, img := range cr.Images {
		images = append(images, pageImageResponse{
			Page: img.Page,
			URL:  fmt.Sprintf("%s/api/pdf-image/session/%s/%s", baseURL, sessionID, filepath.Base(img.Path)),
			Mime: img.Mime,
		})
	}

	return &convertResponse{
		TotalPages: cr.TotalPages,
		Images:     images,
	}
}

// cleanupSession 删除失败的会话目录
func (h *PdfImageHandler) cleanupSession(sessionDir string) {
	_ = os.RemoveAll(sessionDir)
}

// pageImageResponse 图片条目响应结构
type pageImageResponse struct {
	Page int    `json:"page"`
	URL  string `json:"url"`
	Mime string `json:"mime"`
}

// convertResponse 转换接口响应结构
type convertResponse struct {
	TotalPages int                 `json:"total_pages"`
	Images     []pageImageResponse `json:"images"`
}
