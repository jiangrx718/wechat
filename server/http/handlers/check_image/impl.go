package check_image

import (
	checkImageService "wechat-tools/internal/service/check_image"

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
	service checkImageService.ServiceIFace
}

func (h *CheckImageHandler) RegisterRoutes(routerGroup *gin.RouterGroup) {
	g := routerGroup.Group("/check-image")

	g.POST("/check", h.Check)
}
