package agent

import (
	"wechat-tools/utils"
	"wechat-tools/utils/tool/response"

	"github.com/gin-gonic/gin"
)

type AskReq struct {
	Question string `json:"question" binding:"required"`
}

func (h *AgentHandler) Ask(ctx *gin.Context) {
	var reqBody AskReq
	var logger = utils.SugarContext(ctx)
	if err := ctx.Bind(&reqBody); err != nil {
		logger.Infow("Handler Agent Ask ctx.Bind err", "error", err)
		response.ParameterError(ctx, err)
		return
	}

	result, err := h.service.Ask(ctx, reqBody.Question)
	if err != nil {
		logger.Errorw("Handler Agent Ask service.Ask error", "error", err)
		response.InternalError(ctx)
		return
	}

	if result.GetCode() != 0 {
		response.Failed(ctx, result.GetCode(), result.GetMessage(), result.GetData())
		return
	}

	response.Successful(ctx, result.GetData())
}
