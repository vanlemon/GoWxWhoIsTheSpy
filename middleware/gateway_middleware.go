package middleware

import (
	"context"
	"encoding/json"
	"lmf.mortal.com/GoWxWhoIsTheSpy/config"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"lmf.mortal.com/GoWxWhoIsTheSpy/cconst"
	"lmf.mortal.com/GoWxWhoIsTheSpy/model"
	"lmf.mortal.com/GoWxWhoIsTheSpy/util"

	logs "lmf.mortal.com/GoLogs"
)

// 网关中间件（防重入），接口筛选
func GatewayMiddleware() gin.HandlerFunc {
	return func(gctx *gin.Context) {
		// step1: 获取请求 ctx
		ctxValue, exists := gctx.Get(cconst.CtxKey)
		ctx, ok := ctxValue.(context.Context)
		if !exists || !ok {
			panic("[Middleware Gateway] ctx not exists in gin ctx")
		}

		// step2：获取请求方法唯一标识，判断该方法是否提供服务
		methodKey := gctx.Request.Method + ":" + gctx.Request.URL.String() // 获取请求方法唯一标识
		methodKey = strings.Split(methodKey, "?")[0]                       // 去掉 HTTP 请求参数，即？之后的内容
		ctx = util.CtxWithMethodKey(ctx, methodKey)                        // 设置请求方法唯一标识
		gctx.Set(cconst.CtxKey, ctx)                                       // 设置请求方法唯一标识
		if !config.GatewayExists(methodKey) {
			logs.CtxError(ctx, "[Middleware Gateway] gateway not exists methodKey: %#v", methodKey)
			// 包装返回信息
			gctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		// step3: 获取请求签名，以请求签名该请求作为唯一标识，进行防重入操作
		sign := gctx.PostForm("sign")
		if util.IsNil(sign) {
			// step3: 请求无签名，直接返回
			logs.CtxError(ctx, "[Middleware Gateway] no sign")
			// 包装返回信息
			gctx.JSON(http.StatusOK, util.NoSign(ctx))
			gctx.Abort()
			return
		}

		// step4：请求加 redis 锁，在验证签名前加锁，防止无效请求重入
		var requestKey = methodKey + "-" + sign                                         // 请求唯一标识为请求方法唯一标识 + sign
		logs.CtxInfo(ctx, "[Middleware Gateway] requestKey: %s", requestKey)            // 请求唯一标识
		ctx = util.CtxWithRequestKey(ctx, requestKey)                                   // 设置请求参数唯一标识
		gctx.Set(cconst.CtxKey, ctx)                                                    // 设置请求参数唯一标识
		success, err := model.GatewayRedis.SetNX(requestKey, nil, time.Second).Result() // setNX 加 redis 锁

		// step5.1: 加锁失败
		if err != nil {
			logs.CtxError(ctx, "[Middleware Gateway] %s redis lock setNX error: %#v", requestKey, err)
			// 包装返回信息
			gctx.JSON(http.StatusOK, util.ErrorWithMessage(ctx, err.Error()))
			gctx.Abort()
			return
		}

		// step5.2: 加锁重入
		if !success {
			logs.CtxWarn(ctx, "[Middleware Gateway] %s duplicate", requestKey)
			// 包装返回信息
			gctx.JSON(http.StatusOK, util.Duplicate(ctx))
			gctx.Abort()
			return
		}

		// step6: 获取请求参数
		reqDataJsonString := gctx.PostForm("data")
		if util.IsNil(reqDataJsonString) {
			// step6: 请求无参数
			logs.CtxError(ctx, "[Middleware Gateway] no data")
			// 包装返回信息
			gctx.JSON(http.StatusOK, util.NoData(ctx))
			gctx.Abort()
			return
		}
		logs.CtxInfo(ctx, "[Middleware Gateway] request: %#v with sign: %#v", reqDataJsonString, sign)

		// step7: 解析请求参数
		reqDataMap := make(map[string]string)
		err = json.Unmarshal([]byte(reqDataJsonString), &reqDataMap)
		if err != nil {
			logs.CtxError(ctx, "[Middleware Gateway] request Unmarshal error: %#v", err)
		}
		logs.CtxInfo(ctx, "[Middleware Gateway] request: %#v with sign: %#v", reqDataMap, sign)
		gctx.Set(cconst.ReqDataKey, reqDataMap) // 请求参数写入 gin 上下文

		// step8: 检验请求签名
		reqSign := util.GetSign(reqDataJsonString)
		logs.CtxInfo(ctx, "[Middleware Gateway] request: %#v with sign: %#v, expected sign: %#v", reqDataMap, sign, reqSign)
		if sign != reqSign {
			logs.CtxError(ctx, "[Middleware Gateway] sign error")
			// 包装返回信息
			gctx.JSON(http.StatusOK, util.SignError(ctx))
			gctx.Abort()
			return
		}

		// step9: 执行服务
		gctx.Next()
	}
}
