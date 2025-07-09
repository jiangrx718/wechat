package api_handler

import (
	"strconv"

	"github.com/jiangrx816/wechat/common/request"
	"github.com/jiangrx816/wechat/utils/errs"

	"github.com/gin-gonic/gin"
)

/**
 * @Description 默认Demo GET
 **/
func (ph *ApiHandler) ApiIndex(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	response, err := ph.service.ApiServiceIndexGet(ctx, page, size)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}
	ctx.JSON(errs.SucResp(response))
}

/**
 * @Description 默认Demo POST
 **/
func (ph *ApiHandler) ApiIndexPost(ctx *gin.Context) {
	var json request.IndexPosttRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(errs.SucErrResp("", "参数错误"))
		return
	}
	result, err := ph.service.ApiServiceIndexPost(ctx, json)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}
	ctx.JSON(errs.SucResp(result))
}
