package router

import (
	"github.com/jiangrx816/wechat/core/log"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

func All() func(r *gin.Engine) {
	return func(r *gin.Engine) {

		// panic日志
		r.Use(ginzap.RecoveryWithZap(log.Sugar().Desugar(), true))
		r.MaxMultipartMemory = 30 << 20 // 30MB

		prefixRouter := r.Group("/api/")

		// 默认的Api路由
		Api(prefixRouter)

		// 中文绘本路由
		WechatApi(prefixRouter)
	}
}
