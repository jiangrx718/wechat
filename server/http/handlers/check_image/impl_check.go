package check_image

import (
	"io"

	"wechat-tools/server/http/response"
	"wechat-tools/utils"

	"github.com/gin-gonic/gin"
)

// Check 图片内容安全检测接口
//
// 请求: multipart/form-data, 字段 media = 图片文件
// 响应: {"code":0,"msg":"操作成功","data":{"pass":true}}
//
//	pass=true 合规; pass=false 命中违规内容
func (h *CheckImageHandler) Check(ctx *gin.Context) {
	var logger = utils.SugarContext(ctx)

	// 接收前端 wx.uploadFile 上传的图片（字段名 media）
	fileHeader, err := ctx.FormFile("media")
	if err != nil {
		logger.Infow("Handler CheckImage ctx.FormFile err", "error", err)
		response.ParameterError(ctx, err)
		return
	}

	// 限制文件大小（img_sec_check 单图上限约 1MB）
	if fileHeader.Size > 1<<20 { // 1MB
		response.Failed(ctx, response.CodeParameterErr, "图片过大，请选择小于 1MB 的图片", nil)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		logger.Errorw("Handler CheckImage open file err", "error", err)
		response.InternalError(ctx)
		return
	}
	defer file.Close()

	media, err := io.ReadAll(file)
	if err != nil {
		logger.Errorw("Handler CheckImage read file err", "error", err)
		response.InternalError(ctx)
		return
	}

	result, err := h.service.Check(ctx, media, fileHeader.Filename)
	if err != nil {
		logger.Errorw("Handler CheckImage service.Check error", "error", err)
		response.InternalError(ctx)
		return
	}

	if result.GetCode() != 0 {
		response.Failed(ctx, result.GetCode(), result.GetMessage(), nil)
		return
	}

	response.Successful(ctx, result.GetData())
}
