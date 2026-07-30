package miniwechat

import (
	wechatUserService "wechat-tools/internal/service/miniwechat"

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
	service wechatUserService.WechatUserServiceIFace
}

func (h *WechatUserHandler) RegisterRoutes(routerGroup *gin.RouterGroup) {
	g := routerGroup.Group("/miniwechat-user")

	g.POST("/exists", h.Exist)
	g.POST("/create", h.Create)
	g.POST("/update", h.Update)
	g.GET("/list", h.List)
}
