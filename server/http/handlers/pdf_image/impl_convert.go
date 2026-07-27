package pdf_image

import (
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
	"strings"

	pdfImageService "wechat-tools/internal/service/pdf_image"
	"wechat-tools/server/http/response"
	"wechat-tools/utils"

	"github.com/gin-gonic/gin"
)

// Convert PDF 转图片接口。
//
// 全程不在服务端持久化图片，转换结果以 base64 data URI 直接返回，前端（含微信小程序 <image>）可直接用于预览。
//
// 请求: multipart/form-data
//   - file:             PDF 文件（必填）
//   - format:           (可选) 图片格式 png/jpeg，默认 png
//   - dpi:              (可选) 输出分辨率 72-600，默认 150
//   - page:             (可选) 指定转换的页码（从 1 开始），0 或不传表示转换整个 PDF 的全部页
//   - quality:          (可选) 清晰度档位 low/medium/high/original，设置后覆盖 dpi
//   - remove_watermark: (可选) 渲染前是否尝试移除水印，传 1/true 启用；默认 false。
//                      基于结构化移除，识别不到时原样保留，不会误伤正文。
//
// 响应: {"code":0,"msg":"操作成功","data":{"total_pages":N,"images":[{"page":1,"data":"data:image/png;base64,...","mime":"image/png"},...]}}
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
	// quality: 清晰度档位 low/medium/high/original，设置后覆盖 dpi
	quality := ctx.DefaultPostForm("quality", "")
	// remove_watermark: 1/true 启用渲染前去水印（基于结构化移除，识别不到时原样保留）
	rawRW := ctx.DefaultPostForm("remove_watermark", "")
	removeWatermark := rawRW == "1" || strings.EqualFold(rawRW, "true")

	result, err := h.service.Convert(ctx, pdfData, fileHeader.Filename,
		pdfImageService.WithFormat(format),
		pdfImageService.WithDPI(dpi),
		pdfImageService.WithPage(page),
		pdfImageService.WithQuality(quality),
		pdfImageService.WithRemoveWatermark(removeWatermark),
	)
	if err != nil {
		logger.Errorw("Handler PdfImage Convert service.Convert error", "error", err)
		response.InternalError(ctx)
		return
	}

	if result.GetCode() != 0 {
		response.Failed(ctx, result.GetCode(), result.GetMessage(), nil)
		return
	}

	// 构造响应：将 PageImage.Data 编码为 base64 data URI，前端可直接用于预览
	cr := result.GetData().(*pdfImageService.ConvertResult)
	resp := h.buildConvertResponse(cr)

	response.Successful(ctx, resp)
}

// buildConvertResponse 将服务层结果转为带 data URI 的 HTTP 响应体
func (h *PdfImageHandler) buildConvertResponse(cr *pdfImageService.ConvertResult) *convertResponse {
	images := make([]pageImageResponse, 0, len(cr.Images))
	for _, img := range cr.Images {
		dataURI := fmt.Sprintf("data:%s;base64,%s", img.Mime, base64.StdEncoding.EncodeToString(img.Data))
		images = append(images, pageImageResponse{
			Page: img.Page,
			Data: dataURI,
			Mime: img.Mime,
		})
	}

	return &convertResponse{
		TotalPages: cr.TotalPages,
		Images:     images,
	}
}

// pageImageResponse 图片条目响应结构
type pageImageResponse struct {
	Page int    `json:"page"`
	Data string `json:"data"` // base64 data URI，可直接作为 <image src> 预览
	Mime string `json:"mime"`
}

// convertResponse 转换接口响应结构
type convertResponse struct {
	TotalPages int                 `json:"total_pages"`
	Images     []pageImageResponse `json:"images"`
}
