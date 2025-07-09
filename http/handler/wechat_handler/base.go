package wechat_handler

import (
	"github.com/jiangrx816/wechat/service/wechat_service"
)

func NewWechatHandler() *WechatHandler {
	return &WechatHandler{
		service: wechat_service.NewWechatService(),
	}
}

type WechatHandler struct {
	service *wechat_service.WechatService
}
