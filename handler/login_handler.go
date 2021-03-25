package handler

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	logs "lmf.mortal.com/GoLogs"
	"lmf.mortal.com/GoWxWhoIsTheSpy/cconst"
	"lmf.mortal.com/GoWxWhoIsTheSpy/model"
	"lmf.mortal.com/GoWxWhoIsTheSpy/service"
	"lmf.mortal.com/GoWxWhoIsTheSpy/util"
	"net/http"
)

/**
登录服务：
1. wx_login：微信登录
2. update_userInfo：更新用户信息
*/
type LoginHandler struct {
}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{}
}

func (h *LoginHandler) Register(engine *gin.Engine) {
	user := engine.Group("user")
	{
		user.POST("wx_login", WxLogin)
		user.POST("update_userInfo", UpdateUserInfo)
	}
}

func WxLogin(gctx *gin.Context) {
	// step1: 获取请求 ctx，请求数据，请求唯一标识
	// 1.1 请求 ctx
	ctxValue, exists := gctx.Get(cconst.CtxKey)
	ctx, ok := ctxValue.(context.Context)
	if !exists || !ok {
		logs.CtxFatal(ctx, "[Handler User WxLogin] ctx not exists in gin ctx")
		panic("[Handler User WxLogin] ctx not exists in gin ctx")
	}
	// 1.2 请求数据
	reqDataValue, exists := gctx.Get(cconst.ReqDataKey)
	reqData, ok := reqDataValue.(map[string]string)
	if !exists || !ok {
		logs.CtxFatal(ctx, "[Handler User WxLogin] req_data not exists in gin ctx")
		panic("[Handler User WxLogin] req_data not exists in gin ctx")
	}
	// 1.3 请求唯一标识
	requestKey := util.CtxGetString(ctx, util.RequestKey)

	// step2: 定义返回数据，并打印返回日志，删除重入锁
	var resp map[string]interface{}
	defer func() {
		// 打印返回日志
		logs.CtxInfo(ctx, "[Handler User WxLogin] resp: %#v", resp)
		// 删除重入锁
		// Redis DEL 命令用于删除已存在的键。不存在的 key 会被忽略。
		err := model.GatewayRedis.Del(requestKey).Err()
		if err != nil {
			logs.CtxError(ctx, "[Handler User WxLogin] [Redis] redis lock del error: %#v, for requestKey: %#v", err, requestKey)
		}
		// 返回值写入 gctx
		gctx.JSON(http.StatusOK, resp)
	}()
	/********************************** 上述为固定流程 *********************************************/

	// step3: 解析请求数据并校检参数
	// 参数：用户登录码
	code, ok := reqData["code"]
	if !ok || util.IsNil(code) {
		logs.CtxError(ctx, "[Handler User WxLogin] reqData ( %#v ) no param `code`", reqData)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, "[Handler User WxLogin] reqData ( %#v ) no param `code`", reqData)
		return
	}
	logs.CtxInfo(ctx, "[Handler User WxLogin] reqData param `code`: %#v", code)

	// step4: 核心服务层执行请求
	tempId, err := service.CoreUserService.WxLogin(ctx, code)
	if err != nil {
		logs.CtxError(ctx, "[Handler User WxLogin] CoreUserService.WxLogin err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler User WxLogin] CoreUserService.WxLogin resp Success with tempId: %#v", tempId)

	/********************************** 下述为固定流程 *********************************************/
	// step5: 包装返回值
	resp = util.SuccessWithData(ctx, tempId)
}

func UpdateUserInfo(gctx *gin.Context) {
	// step1: 获取请求 ctx，请求数据，请求唯一标识
	// 1.1 请求 ctx
	ctxValue, exists := gctx.Get(cconst.CtxKey)
	ctx, ok := ctxValue.(context.Context)
	if !exists || !ok {
		logs.CtxFatal(ctx, "[Handler User UpdateUserInfo] ctx not exists in gin ctx")
		panic("[Handler User UpdateUserInfo] ctx not exists in gin ctx")
	}
	// 1.2 请求数据
	reqDataValue, exists := gctx.Get(cconst.ReqDataKey)
	reqData, ok := reqDataValue.(map[string]string)
	if !exists || !ok {
		logs.CtxFatal(ctx, "[Handler User UpdateUserInfo] req_data not exists in gin ctx")
		panic("[Handler User UpdateUserInfo] req_data not exists in gin ctx")
	}
	// 1.3 请求唯一标识
	requestKey := util.CtxGetString(ctx, util.RequestKey)

	// step2: 定义返回数据，并打印返回日志，删除重入锁
	var resp map[string]interface{}
	defer func() {
		// 打印返回日志
		logs.CtxInfo(ctx, "[Handler User UpdateUserInfo] resp: %#v", resp)
		// 删除重入锁
		// Redis DEL 命令用于删除已存在的键。不存在的 key 会被忽略。
		err := model.GatewayRedis.Del(requestKey).Err()
		if err != nil {
			logs.CtxError(ctx, "[Handler User UpdateUserInfo] [Redis] redis lock del error: %#v, for requestKey: %#v", err, requestKey)
		}
		// 返回值写入 gctx
		gctx.JSON(http.StatusOK, resp)
	}()
	/********************************** 上述为固定流程 *********************************************/

	// step3: 解析请求数据并校检参数
	// 参数：用户临时 ID
	tempId, ok := reqData["tempId"]
	if !ok || util.IsNil(tempId) {
		logs.CtxError(ctx, "[Handler User UpdateUserInfo] reqData ( %#v ) no param `tempId`", reqData)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, "[Handler User UpdateUserInfo] reqData ( %#v ) no param `tempId`", reqData)
		return
	}
	logs.CtxInfo(ctx, "[Handler User UpdateUserInfo] reqData param `tempId`: %#v", tempId)
	// 参数：用户信息
	userInfoValue, ok := reqData["userInfo"]
	if !ok || util.IsNil(userInfoValue) {
		logs.CtxError(ctx, "[Handler User UpdateUserInfo] reqData ( %#v ) no param `userInfo`", reqData)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, "[Handler User UpdateUserInfo] reqData ( %#v ) no param `userInfo`", reqData)
		return
	}
	logs.CtxInfo(ctx, "[Handler User UpdateUserInfo] reqData param `userInfo`: %#v", userInfoValue)
	// 解析参数：用户信息
	userInfo := &model.User{}
	err := json.Unmarshal([]byte(userInfoValue), userInfo)
	if err != nil {
		logs.CtxError(ctx, "[Handler User UpdateUserInfo] reqData ( %#v ) param `userInfo` err: %#v", reqData, err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, "[Handler User UpdateUserInfo] reqData ( %#v ) param `userInfo` err: %#v", reqData, err)
		return
	}
	logs.CtxInfo(ctx, "[Handler User UpdateUserInfo] reqData param `userInfo`: %#v", userInfo)

	// step4: 核心服务层执行请求
	openId,err := service.CoreUserService.GetUserOpenId(ctx, tempId)
	if err == cconst.RedisNil {
		logs.CtxError(ctx, "[Handler User UpdateUserInfo] CoreUserService.GetUserOpenId err: %#v", err)
		// 包装返回信息
		resp = util.UserNotLogin(ctx)
		return
	}
	if err != nil {
		logs.CtxError(ctx, "[Handler User UpdateUserInfo] CoreUserService.GetUserOpenId err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler User UpdateUserInfo] CoreUserService.GetUserOpenId resp Success")
	err = service.CoreUserService.UpdateUserInfo(ctx, openId, userInfo)
	if err != nil {
		logs.CtxError(ctx, "[Handler User UpdateUserInfo] CoreUserService.UpdateUserInfo err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler User UpdateUserInfo] CoreUserService.UpdateUserInfo resp Success")

	/********************************** 下述为固定流程 *********************************************/
	// step5: 包装返回值
	resp = util.Success(ctx)
}
