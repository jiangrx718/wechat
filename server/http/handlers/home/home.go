package home

// home.go 首页玩法卡片接口
//
// 用途：小程序首页 pages/home/home 通过本接口拉取玩法卡片配置，
// 后端通过 visible 字段控制每张卡片 / 每个分类是否在首页展示。
//
// 路由：GET /api/home/categories
// 响应：{"code":0,"msg":"操作成功","data":{"categories":[...]}}

import (
	"sort"

	"github.com/gin-gonic/gin"

	"wechat-tools/server/http/response"
)

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

// categoriesResp 响应数据
type categoriesResp struct {
	DailyEntry dailyEntryDTO `json:"dailyEntry"` // 「今日推荐玩法」入口卡片配置
	Categories []categoryDTO `json:"categories"`
}

// Categories 获取首页玩法卡片配置
//
// 下发规则：
//  1. dailyEntry 直接下发（Visible 控制今日推荐入口卡片是否展示）；
//  2. 整个分类 Visible=false 时不下发该分类；
//  3. 单张卡片 Visible=false 时不下发该卡片（分类仍下发）；
//  4. 分类内卡片全部隐藏时，整个分类不下发。
// 返回数据已按 Sort 升序排好，前端可直接渲染。
func (h *HomeHandler) Categories(ctx *gin.Context) {
	out := make([]categoryDTO, 0, len(cardsData))

	for _, cat := range cardsData {
		if !cat.Visible {
			continue // 分类整体下架
		}

		// 过滤隐藏卡片
		visibleGames := make([]gameDTO, 0, len(cat.Games))
		for _, g := range cat.Games {
			if g.Visible {
				visibleGames = append(visibleGames, g)
			}
		}
		if len(visibleGames) == 0 {
			continue // 分类内无可见卡片，跳过整个分类
		}

		// 排序（防御性，data.go 已按顺序写好）
		sort.Slice(visibleGames, func(i, j int) bool {
			return visibleGames[i].Sort < visibleGames[j].Sort
		})

		cat.Games = visibleGames
		out = append(out, cat)
	}

	// 分类整体排序
	sort.Slice(out, func(i, j int) bool {
		return out[i].Sort < out[j].Sort
	})

	response.Successful(ctx, categoriesResp{
		DailyEntry: dailyEntryData,
		Categories: out,
	})
}
