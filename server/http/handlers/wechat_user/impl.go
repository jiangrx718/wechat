package wechat_user

import (
	wechatUserService "wechat-tools/internal/service/wechat_user"

	"github.com/gin-gonic/gin"
)

func NewWechatUserHandler(engine *gin.Engine) *WechatUserHandler {
	return &WechatUserHandler{
		engine:  engine,
		service: wechatUserService.NewWechatUserService(),
	}
}

type WechatUserHandler struct {
	engine  *gin.Engine
	service wechatUserService.ServiceIFace
}

func (h *WechatUserHandler) RegisterRoutes(routerGroup *gin.RouterGroup) {
	g := routerGroup.Group("/wechat-user")

	g.POST("/exists", h.Exist)
	g.POST("/create", h.Create)
	g.POST("/update", h.Update)
	g.GET("/list", h.List)
}
