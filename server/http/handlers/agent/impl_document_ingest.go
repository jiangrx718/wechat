package agent

import (
	"io"

	"wechat-tools/server/http/response"
	"wechat-tools/utils"

	"github.com/gin-gonic/gin"
)

// 入库文件大小上限：20MB
const ingestMaxSize = 20 << 20

// Ingest 单文件入库接口
//
// 请求: multipart/form-data, 字段 file = 文件（支持 txt/markdown/pdf/docx）
// 响应: {"code":0,"msg":"操作成功","data":{"document_count":1,"chunk_count":12}}
func (h *AgentHandler) Ingest(ctx *gin.Context) {
	var logger = utils.SugarContext(ctx)

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		logger.Infow("Handler Agent Ingest ctx.FormFile err", "error", err)
		response.ParameterError(ctx, err)
		return
	}

	if fileHeader.Size > ingestMaxSize {
		response.Failed(ctx, response.CodeParameterErr, "文件过大，请选择小于 20MB 的文件", nil)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		logger.Errorw("Handler Agent Ingest open file err", "error", err)
		response.InternalError(ctx)
		return
	}
	defer file.Close()

	media, err := io.ReadAll(file)
	if err != nil {
		logger.Errorw("Handler Agent Ingest read file err", "error", err)
		response.InternalError(ctx)
		return
	}

	result, err := h.service.Ingest(ctx, fileHeader.Filename, media)
	if err != nil {
		logger.Errorw("Handler Agent Ingest service.Ingest error", "error", err)
		response.InternalError(ctx)
		return
	}

	if result.GetCode() != 0 {
		response.Failed(ctx, result.GetCode(), result.GetMessage(), nil)
		return
	}

	response.Successful(ctx, result.GetData())
}
