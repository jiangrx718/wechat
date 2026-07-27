package pdf_image

import (
	pdfImageService "wechat-tools/internal/service/pdf_image"

	"github.com/gin-gonic/gin"
)

func NewPdfImageHandler() *PdfImageHandler {
	return &PdfImageHandler{
		service: pdfImageService.NewPdfImageService(),
	}
}

type PdfImageHandler struct {
	service pdfImageService.ServiceIFace
}

func (h *PdfImageHandler) RegisterRoutes(routerGroup *gin.RouterGroup) {
	g := routerGroup.Group("/pdf-image")
	g.POST("/convert", h.Convert)
}
