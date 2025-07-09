package api_service

import (
	"github.com/gin-gonic/gin"
	"github.com/jiangrx816/wechat/common/request"
	"github.com/jiangrx816/wechat/common/response"
	rxLog "github.com/jiangrx816/wechat/core/log"
	"github.com/jiangrx816/wechat/core/server/api"
)

func (ps *ApiService) ApiServiceIndexGet(ctx *gin.Context, page int, size int) (resp response.IndexResponse, apiErr api.Error) {
	rxLog.Sugar().Infof("page: %+v size: %+v", page, size)

	return
}

func (ps *ApiService) ApiServiceIndexPost(ctx *gin.Context, params request.IndexPosttRequest) (resp response.IndexResponse, apiErr api.Error) {
	rxLog.Sugar().Infof("params: %+v", params)

	return
}
