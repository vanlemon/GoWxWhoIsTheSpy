package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"lmf.mortal.com/GoWxWhoIsTheSpy/cconst"
	"lmf.mortal.com/GoWxWhoIsTheSpy/config"

	logs "lmf.mortal.com/GoLogs"
)

// 日志服务中间件
func LogsMiddleware() gin.HandlerFunc {
	return func(gctx *gin.Context) {
		// step1: 根据运行环境初始化 ctx，每次请求会分配唯一的 logid
		ctx := logs.Ctx(config.ConfigInstance.Env)
		gctx.Set(cconst.CtxKey, ctx) // 设置 ctx

		// step2: 打印请求日志
		host := gctx.Request.Host     // 请求主机
		url := gctx.Request.URL       // 请求 url
		method := gctx.Request.Method // 请求接口
		reqTime := time.Now()            // 请求时间
		logs.CtxInfo(ctx, "[Middleware Access] %s \t %s \t %s \t %s", reqTime.Format("2006-01-02 15:04:05"), host, url, method)

		// step3: 执行服务
		gctx.Next()

		// step4: 打印返回日志
		respTime := time.Now()            // 返回时间
		costTime := respTime.Sub(reqTime) // 耗时
		logs.CtxInfo(ctx, "[Middleware Access] %s \t %s \t %s \t %s \t %#v \t %#v", respTime.Format("2006-01-02 15:04:05"), host, url, method, gctx.Writer.Status(), costTime)
	}
}
