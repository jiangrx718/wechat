package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jiangrx816/wechat/http/handler/wechat_handler"
)

func WechatApi(r *gin.RouterGroup) {

	prefixRouter := r.Group("v1").Use()
	apiHandler := wechat_handler.NewWechatHandler()
	{
		// 绘本栏目
		prefixRouter.GET("/init/getCategoryList", apiHandler.ApiGetCategoryList)
		// 绘本列表
		prefixRouter.GET("/chinese/getBookList", apiHandler.ApiChineseBookList)
		// 绘本详情
		prefixRouter.GET("/chinese/getBookInfo", apiHandler.ApiChineseBookInfo)
		// 绘本搜索
		prefixRouter.GET("/chinese/getBookSearch", apiHandler.ApiChineseBookSearch)
		// 英文绘本列表
		prefixRouter.GET("/english/getBookList", apiHandler.ApiEnglishBookList)
		// 英文绘本详情
		prefixRouter.GET("/english/getBookInfo", apiHandler.ApiEnglishBookInfo)
		// 英文绘本搜索
		prefixRouter.GET("/english/getBookSearch", apiHandler.ApiEnglishBookSearch)
		// 处理数据
		prefixRouter.POST("/english/handleData", apiHandler.ApiEnglishHandleData)

		//古诗绘本列表
		prefixRouter.GET("/poetry/getBookList", apiHandler.ApiPoetryBookList)
		//古诗绘本详情
		prefixRouter.GET("/poetry/getBookInfo", apiHandler.ApiPoetryBookInfo)

		// 数学计算相关
		prefixRouter.GET("/math/getCalculationList", apiHandler.ApiMathCalculationList)

		//微信key
		prefixRouter.GET("/wechat/signature", apiHandler.ApiWechatSignature)
	}
}
