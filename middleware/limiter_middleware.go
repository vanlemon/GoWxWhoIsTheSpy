package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	limiter "lmf.mortal.com/GoLimiter"
	"lmf.mortal.com/GoWxWhoIsTheSpy/cconst"
	"lmf.mortal.com/GoWxWhoIsTheSpy/util"
	"net/http"

	logs "lmf.mortal.com/GoLogs"
)

// 限流器中间件
func LimiterMiddleware() gin.HandlerFunc {
	return func(gctx *gin.Context) {
		// step1: 获取请求 ctx
		ctxValue, exists := gctx.Get(cconst.CtxKey)
		ctx, ok := ctxValue.(context.Context)
		if !exists || !ok {
			panic("[Middleware Limiter] ctx not exists in gin ctx")
		}

		// step2: 打印请求日志
		methodKey := util.CtxGetString(ctx, util.MethodKey)// 获取请求方法唯一标识
		pass := limiter.CanPass(ctx, methodKey)
		if pass {
			logs.CtxInfo(ctx, "[Middleware Limiter] %s CanPass? \t %#v", methodKey, pass)
			// step3: 执行服务
			gctx.Next()
		} else {
			logs.CtxWarn(ctx, "[Middleware Limiter] %s CanPass? \t %#v", methodKey, pass)
			// step3: 服务器超载，直接返回
			// 包装返回信息
			gctx.JSON(http.StatusOK, util.Overload(ctx))
			gctx.Abort()
		}
	}
}
