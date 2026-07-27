package pdf_image

import (
	"io"
	"strconv"

	pdfImageService "wechat-tools/internal/service/pdf_image"
	"wechat-tools/server/http/response"
	"wechat-tools/utils"

	"github.com/gin-gonic/gin"
)

// Convert PDF 转图片接口
//
// 请求: multipart/form-data
//   - file: PDF 文件
//   - format: (可选) 图片格式，png 或 jpeg，默认 png
//   - dpi: (可选) 输出分辨率，默认 150
//
// 响应: {"code":0,"msg":"操作成功","data":{"total_pages":N,"images":[{"page":1,"data":"base64...","mime":"image/png"},...]}}
func (h *PdfImageHandler) Convert(ctx *gin.Context) {
	var logger = utils.SugarContext(ctx)

	// 接收上传的 PDF 文件（字段名 file）
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

	// 调用服务层转换
	result, err := h.service.Convert(ctx, pdfData, fileHeader.Filename,
		pdfImageService.WithFormat(format),
		pdfImageService.WithDPI(dpi),
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

	response.Successful(ctx, result.GetData())
}
