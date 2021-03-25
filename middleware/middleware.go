package middleware

import (
	"github.com/gin-gonic/gin"
)

// 初始化中间件服务，即将所有中间件注册到 gin 服务器上
func InitMiddleware(engine *gin.Engine) {
	// 中间件服务列表，按顺序
	middlewares := []gin.HandlerFunc{
		LogsMiddleware(),    // 先打印请求日志
		GatewayMiddleware(), // 再过网关，判断服务是否可用，是否重入
		LimiterMiddleware(), // 再过限流器，判断是否超负荷
	}

	// 中间件服务注册到 Gin 服务器上
	engine.Use(middlewares...)
}
