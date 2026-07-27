package pdf_image

import (
	"os"
	"path/filepath"
	"time"

	pdfImageService "wechat-tools/internal/service/pdf_image"

	"github.com/gin-gonic/gin"
)

// defaultStaticDir 转换后的图片文件在服务端的持久化根目录
const defaultStaticDir = "server/static/pdf_images"

func NewPdfImageHandler(engine *gin.Engine) *PdfImageHandler {
	return &PdfImageHandler{
		engine:    engine,
		service:   pdfImageService.NewPdfImageService(),
		staticDir: defaultStaticDir,
	}
}

type PdfImageHandler struct {
	engine    *gin.Engine
	service   pdfImageService.ServiceIFace
	staticDir string
}

func (h *PdfImageHandler) RegisterRoutes(routerGroup *gin.RouterGroup) {
	absDir, err := filepath.Abs(h.staticDir)
	if err != nil {
		panic("failed to resolve pdf image static dir: " + err.Error())
	}
	if err := os.MkdirAll(absDir, 0755); err != nil {
		panic("failed to create pdf image static dir: " + err.Error())
	}
	h.staticDir = absDir

	g := routerGroup.Group("/pdf-image")

	// 静态文件路由：/api/pdf-image/session/<sessionId>/<filename>
	g.Static("/session", absDir)

	g.POST("/convert", h.Convert)

	// 后台清理：每 10 分钟删除超过 30 分钟的旧会话目录
	go h.cleanupLoop(30*time.Minute, 10*time.Minute)
}

// cleanupLoop 定期清理过期的会话目录
func (h *PdfImageHandler) cleanupLoop(maxAge, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		entries, err := os.ReadDir(h.staticDir)
		if err != nil {
			continue
		}
		now := time.Now()
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			info, err := e.Info()
			if err != nil {
				continue
			}
			if now.Sub(info.ModTime()) > maxAge {
				os.RemoveAll(filepath.Join(h.staticDir, e.Name()))
			}
		}
	}
}
