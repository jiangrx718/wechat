package wechat_user

import (
	wechatUserService "wechat-tools/internal/service/wechat_user"
	"wechat-tools/server/http/response"
	"wechat-tools/utils"

	"github.com/gin-gonic/gin"
)

func (h *WechatUserHandler) List(ctx *gin.Context) {
	var logger = utils.SugarContext(ctx)
	result, err := h.service.List(ctx)
	if err != nil {
		logger.Errorw("Handler WechatUser List service.List error", "error", err)
		response.InternalError(ctx)
		return
	}

	if result.GetCode() != 0 {
		response.Failed(ctx, result.GetCode(), result.GetMessage(), result.GetData())
		return
	}

	data, ok := result.GetData().(wechatUserService.ListResponseData)
	if !ok {
		response.InternalError(ctx)
		return
	}

	offset, limit, total := 1, 1, 1
	response.SuccessfulWithPagination(
		ctx,
		data.List,
		&offset,
		&limit,
		&total,
	)
}
