package miniwechat

import (
	checkImageService "wechat-tools/internal/service/miniwechat"

	"github.com/gin-gonic/gin"
)

func NewCheckImageHandler(engine *gin.Engine) *CheckImageHandler {
	return &CheckImageHandler{
		engine:  engine,
		service: checkImageService.NewCheckImageService(),
	}
}

type CheckImageHandler struct {
	engine  *gin.Engine
	service checkImageService.CheckImageServiceIFace
}

func (h *CheckImageHandler) RegisterRoutes(routerGroup *gin.RouterGroup) {
	g := routerGroup.Group("/check-image")

	g.POST("/check", h.Check)
}
