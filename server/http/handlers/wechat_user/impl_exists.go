package wechat_user

import (
	"wechat-tools/server/http/response"
	"wechat-tools/utils"

	"github.com/gin-gonic/gin"
)

type ExistReq struct {
	DeviceId string `json:"device_id" binding:"required"`
}

func (h *WechatUserHandler) Exist(ctx *gin.Context) {
	var reqBody ExistReq
	var logger = utils.SugarContext(ctx)
	if err := ctx.Bind(&reqBody); err != nil {
		logger.Infow("Handler WechatUser Exist ctx.Bind err", "error", err)
		response.ParameterError(ctx, err)
		return
	}

	result, err := h.service.Exist(ctx, reqBody.DeviceId)
	if err != nil {
		logger.Errorw("Handler WechatUser Exist service.Create error", "error", err)
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
