package wechat_user

import (
	"wechat-tools/server/http/response"
	"wechat-tools/utils"

	"github.com/gin-gonic/gin"
)

type CreateReq struct {
	UserName string `json:"user_name" binding:"required"`
	DeviceId string `json:"device_id" binding:"required"`
	Score    int    `json:"score"`
}

func (h *WechatUserHandler) Create(ctx *gin.Context) {
	var reqBody CreateReq
	var logger = utils.SugarContext(ctx)
	if err := ctx.Bind(&reqBody); err != nil {
		logger.Infow("Handler WechatUser Create ctx.Bind err", "error", err)
		response.ParameterError(ctx, err)
		return
	}

	result, err := h.service.Create(ctx, reqBody.UserName, reqBody.DeviceId, reqBody.Score)
	if err != nil {
		logger.Errorw("Handler WechatUser Create service.Create error", "error", err)
		response.InternalError(ctx)
		return
	}

	if result.GetCode() != 0 {
		response.Failed(ctx, result.GetCode(), result.GetMessage(), result.GetData())
		return
	}

	ctx.JSON(200, result)
	return
}
