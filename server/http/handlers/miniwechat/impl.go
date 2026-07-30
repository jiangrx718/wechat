package miniwechat

import (
	checkImageService "wechat-tools/internal/service/miniwechat"
	wechatUserService "wechat-tools/internal/service/miniwechat"

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

// NewHomeHandler 创建首页卡片处理器
func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

// HomeHandler 首页卡片处理器
type HomeHandler struct{}

// RegisterRoutes 注册路由
func (h *HomeHandler) RegisterRoutes(routerGroup *gin.RouterGroup) {
	g := routerGroup.Group("/home")

	// GET /api/home/categories 获取首页玩法卡片配置
	g.GET("/categories", h.Categories)
}

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
