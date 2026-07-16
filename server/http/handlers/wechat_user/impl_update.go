package wechat_user

import (
	"wechat-tools/server/http/response"
	"wechat-tools/utils"

	"github.com/gin-gonic/gin"
)

type UpdateReq struct {
	UserName string `json:"user_name" binding:"required"`
	Score    int    `json:"score"`
}

// Update 更新绘本接口
func (h *WechatUserHandler) Update(ctx *gin.Context) {
	var reqBody UpdateReq
	var logger = utils.SugarContext(ctx)
	if err := ctx.Bind(&reqBody); err != nil {
		logger.Infow("Handler WechatUser Update ctx.Bind err", "error", err)
		response.ParameterError(ctx, err)
		return
	}

	result, err := h.service.Update(ctx, reqBody.UserName, reqBody.Score)
	if err != nil {
		logger.Errorw("Handler WechatUser Update service.Update error", "error", err)
		response.InternalError(ctx)
		return
	}

	if result.GetCode() != 0 {
		response.Failed(ctx, result.GetCode(), result.GetMessage(), result.GetData())
		return
	}

	response.Successful(ctx, result.GetData())
}
