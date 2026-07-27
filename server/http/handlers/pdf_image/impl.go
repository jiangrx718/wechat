package pdf_image

import (
	pdfImageService "wechat-tools/internal/service/pdf_image"

	"github.com/gin-gonic/gin"
)

func NewPdfImageHandler(engine *gin.Engine) *PdfImageHandler {
	return &PdfImageHandler{
		engine:  engine,
		service: pdfImageService.NewPdfImageService(),
	}
}

type PdfImageHandler struct {
	engine  *gin.Engine
	service pdfImageService.ServiceIFace
}

func (h *PdfImageHandler) RegisterRoutes(routerGroup *gin.RouterGroup) {
	g := routerGroup.Group("/pdf-image")

	g.POST("/convert", h.Convert)
}
