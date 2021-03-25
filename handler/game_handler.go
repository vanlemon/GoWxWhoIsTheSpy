package handler

import (
	"context"
	"lmf.mortal.com/GoWxWhoIsTheSpy/service"
	"net/http"

	logs "lmf.mortal.com/GoLogs"

	"lmf.mortal.com/GoWxWhoIsTheSpy/cconst"
	"lmf.mortal.com/GoWxWhoIsTheSpy/model"
	"lmf.mortal.com/GoWxWhoIsTheSpy/util"

	"github.com/gin-gonic/gin"
)

/**
游戏服务：
1. ready_game：（房客）准备游戏
2. start_game：（房主）开始游戏
3. end_game：（房主）结束游戏
4. restart_game：（房主）重开游戏（结束游戏 - 揭晓谜底 + 开始游戏 = 重开游戏）
*/
type GameHandler struct {
}

func NewGameHandler() *GameHandler {
	return &GameHandler{}
}

func (h *GameHandler) Register(engine *gin.Engine) {
	game := engine.Group("game")
	{
		game.POST("ready_game", ReadyGame)
		game.POST("start_game", StartGame)
		game.POST("end_game", EndGame)
		game.POST("restart_game", RestartGame)
	}
}

func ReadyGame(gctx *gin.Context) {
	// step1: 获取请求 ctx，请求数据，请求唯一标识
	// 1.1 请求 ctx
	ctxValue, exists := gctx.Get(cconst.CtxKey)
	ctx, ok := ctxValue.(context.Context)
	if !exists || !ok {
		logs.CtxFatal(ctx, "[Handler Game ReadyGame] ctx not exists in gin ctx")
		panic("[Handler Game ReadyGame] ctx not exists in gin ctx")
	}
	// 1.2 请求数据
	reqDataValue, exists := gctx.Get(cconst.ReqDataKey)
	reqData, ok := reqDataValue.(map[string]string)
	if !exists || !ok {
		logs.CtxFatal(ctx, "[Handler Game ReadyGame] req_data not exists in gin ctx")
		panic("[Handler Game ReadyGame] req_data not exists in gin ctx")
	}
	// 1.3 请求唯一标识
	requestKey := util.CtxGetString(ctx, util.RequestKey)
	// 1.4 接口唯一标识
	methodKey := util.CtxGetString(ctx, util.MethodKey)

	// step2: 定义返回数据，并打印返回日志，删除重入锁
	var resp map[string]interface{}
	defer func() {
		// 打印返回日志
		logs.CtxInfo(ctx, "[Handler Game ReadyGame] resp: %#v", resp)
		// 删除重入锁
		// Redis DEL 命令用于删除已存在的键。不存在的 key 会被忽略。
		err := model.GatewayRedis.Del(requestKey).Err()
		if err != nil {
			logs.CtxError(ctx, "[Handler Game ReadyGame] [Redis] redis lock del error: %#v, for requestKey: %#v", err, requestKey)
		}
		// 返回值写入 gctx
		gctx.JSON(http.StatusOK, resp)
	}()
	/********************************** 上述为固定流程 *********************************************/

	// step3: 解析请求数据并校检参数
	// 参数：用户临时 ID
	tempId, ok := reqData["tempId"]
	if !ok || util.IsNil(tempId) {
		logs.CtxError(ctx, "[Handler Game ReadyGame] reqData ( %#v ) no param `tempId`", reqData)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, "[Handler Game ReadyGame] reqData ( %#v ) no param `tempId`", reqData)
		return
	}
	logs.CtxInfo(ctx, "[Handler Game ReadyGame] reqData param `tempId`: %#v", tempId)

	// step4: 核心服务层执行请求
	// 用户服务：获取用户 openId
	openId, err := service.CoreUserService.GetUserOpenId(ctx, tempId)
	if err != nil {
		logs.CtxError(ctx, "[Handler Game ReadyGame] CoreUserService.GetUserOpenId err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Game ReadyGame] CoreUserService.GetUserOpenId resp Success")
	// 房间服务：获取房间 Id
	roomId, err := service.CoreRoomService.GetRoomId(ctx, openId)
	if err == cconst.RedisNil {
		logs.CtxError(ctx, "[Handler Game ReadyGame] CoreRoomService.GetRoomId RoomInvalid")
		// 包装返回信息
		resp = util.RoomInvalid(ctx)
		return
	}
	if err != nil {
		logs.CtxError(ctx, "[Handler Game ReadyGame] CoreRoomService.GetRoomId err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Game ReadyGame] CoreRoomService.GetRoomId resp Success")
	// 房间服务：获取房间信息和房间锁
	roomInfo, roomInfoUnlockFunc, err := service.CoreRoomService.GetRoomInfoAndLockFromRoomId(ctx, openId, roomId, methodKey) // TODO：cconst.RedisNil
	if err != nil {
		logs.CtxError(ctx, "[Handler Game ReadyGame] CoreRoomService.GetRoomInfoAndLockFromRoomId err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Game ReadyGame] CoreRoomService.GetRoomInfoAndLockFromRoomId resp Success")
	// 释放房间锁
	defer func() {
		unlockErr := roomInfoUnlockFunc(roomInfo, err)
		if unlockErr != nil {
			logs.CtxError(ctx, "[Handler Game ReadyGame] defer roomInfoUnlockFunc error: %#v", unlockErr)
			// 包装返回信息
			resp = util.ErrorWithMessage(ctx, unlockErr.Error())
			return
		}
		logs.CtxInfo(ctx, "[Handler Game ReadyGame] defer roomInfoUnlockFunc resp Success")
	}()
	// 游戏服务：执行逻辑
	roomInfo, err = service.CoreGameService.ReadyGame(ctx, openId, roomInfo)
	if err != nil {
		logs.CtxError(ctx, "[Handler Game ReadyGame] CoreGameService.ReadyGame err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Game ReadyGame] CoreGameService.ReadyGame resp Success")
	// TODO WebSocket

	/********************************** 下述为固定流程 *********************************************/
	// step5: 包装返回值
	resp = util.Success(ctx)
}

func StartGame(gctx *gin.Context) {
	// step1: 获取请求 ctx，请求数据，请求唯一标识
	// 1.1 请求 ctx
	ctxValue, exists := gctx.Get(cconst.CtxKey)
	ctx, ok := ctxValue.(context.Context)
	if !exists || !ok {
		logs.CtxFatal(ctx, "[Handler Game StartGame] ctx not exists in gin ctx")
		panic("[Handler Game StartGame] ctx not exists in gin ctx")
	}
	// 1.2 请求数据
	reqDataValue, exists := gctx.Get(cconst.ReqDataKey)
	reqData, ok := reqDataValue.(map[string]string)
	if !exists || !ok {
		logs.CtxFatal(ctx, "[Handler Game StartGame] req_data not exists in gin ctx")
		panic("[Handler Game StartGame] req_data not exists in gin ctx")
	}
	// 1.3 请求唯一标识
	requestKey := util.CtxGetString(ctx, util.RequestKey)
	// 1.4 接口唯一标识
	methodKey := util.CtxGetString(ctx, util.MethodKey)

	// step2: 定义返回数据，并打印返回日志，删除重入锁
	var resp map[string]interface{}
	defer func() {
		// 打印返回日志
		logs.CtxInfo(ctx, "[Handler Game StartGame] resp: %#v", resp)
		// 删除重入锁
		// Redis DEL 命令用于删除已存在的键。不存在的 key 会被忽略。
		err := model.GatewayRedis.Del(requestKey).Err()
		if err != nil {
			logs.CtxError(ctx, "[Handler Game StartGame] [Redis] redis lock del error: %#v, for requestKey: %#v", err, requestKey)
		}
		// 返回值写入 gctx
		gctx.JSON(http.StatusOK, resp)
	}()
	/********************************** 上述为固定流程 *********************************************/

	// step3: 解析请求数据并校检参数
	// 参数：用户临时 ID
	tempId, ok := reqData["tempId"]
	if !ok || util.IsNil(tempId) {
		logs.CtxError(ctx, "[Handler Game StartGame] reqData ( %#v ) no param `tempId`", reqData)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, "[Handler Game StartGame] reqData ( %#v ) no param `tempId`", reqData)
		return
	}
	logs.CtxInfo(ctx, "[Handler Game StartGame] reqData param `tempId`: %#v", tempId)

	// step4: 核心服务层执行请求
	// 用户服务：获取用户 openId
	openId, err := service.CoreUserService.GetUserOpenId(ctx, tempId)
	if err != nil {
		logs.CtxError(ctx, "[Handler Game StartGame] CoreUserService.GetUserOpenId err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Game StartGame] CoreUserService.GetUserOpenId resp Success")
	// 房间服务：获取房间 Id
	roomId, err := service.CoreRoomService.GetRoomId(ctx, openId)
	if err == cconst.RedisNil {
		logs.CtxError(ctx, "[Handler Game StartGame] CoreRoomService.GetRoomId RoomInvalid")
		// 包装返回信息
		resp = util.RoomInvalid(ctx)
		return
	}
	if err != nil {
		logs.CtxError(ctx, "[Handler Game StartGame] CoreRoomService.GetRoomId err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Game StartGame] CoreRoomService.GetRoomId resp Success")
	// 房间服务：获取房间信息和房间锁
	roomInfo, roomInfoUnlockFunc, err := service.CoreRoomService.GetRoomInfoAndLockFromRoomId(ctx, openId, roomId, methodKey) // TODO：cconst.RedisNil
	if err != nil {
		logs.CtxError(ctx, "[Handler Game StartGame] CoreRoomService.GetRoomInfoAndLockFromRoomId err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Game StartGame] CoreRoomService.GetRoomInfoAndLockFromRoomId resp Success")
	// 释放房间锁
	defer func() {
		unlockErr := roomInfoUnlockFunc(roomInfo, err)
		if unlockErr != nil {
			logs.CtxError(ctx, "[Handler Game StartGame] defer roomInfoUnlockFunc error: %#v", unlockErr)
			// 包装返回信息
			resp = util.ErrorWithMessage(ctx, unlockErr.Error())
			return
		}
		logs.CtxInfo(ctx, "[Handler Game StartGame] defer roomInfoUnlockFunc resp Success")
	}()
	// 游戏服务：执行逻辑
	roomInfo, err = service.CoreGameService.StartGame(ctx, openId, roomInfo)
	if err != nil {
		logs.CtxError(ctx, "[Handler Game StartGame] CoreGameService.StartGame err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Game StartGame] CoreGameService.StartGame resp Success")
	// TODO WebSocket

	/********************************** 下述为固定流程 *********************************************/
	// step5: 包装返回值
	resp = util.Success(ctx)
}

func EndGame(gctx *gin.Context) {
	// step1: 获取请求 ctx，请求数据，请求唯一标识
	// 1.1 请求 ctx
	ctxValue, exists := gctx.Get(cconst.CtxKey)
	ctx, ok := ctxValue.(context.Context)
	if !exists || !ok {
		logs.CtxFatal(ctx, "[Handler Game EndGame] ctx not exists in gin ctx")
		panic("[Handler Game EndGame] ctx not exists in gin ctx")
	}
	// 1.2 请求数据
	reqDataValue, exists := gctx.Get(cconst.ReqDataKey)
	reqData, ok := reqDataValue.(map[string]string)
	if !exists || !ok {
		logs.CtxFatal(ctx, "[Handler Game EndGame] req_data not exists in gin ctx")
		panic("[Handler Game EndGame] req_data not exists in gin ctx")
	}
	// 1.3 请求唯一标识
	requestKey := util.CtxGetString(ctx, util.RequestKey)
	// 1.4 接口唯一标识
	methodKey := util.CtxGetString(ctx, util.MethodKey)

	// step2: 定义返回数据，并打印返回日志，删除重入锁
	var resp map[string]interface{}
	defer func() {
		// 打印返回日志
		logs.CtxInfo(ctx, "[Handler Game EndGame] resp: %#v", resp)
		// 删除重入锁
		// Redis DEL 命令用于删除已存在的键。不存在的 key 会被忽略。
		err := model.GatewayRedis.Del(requestKey).Err()
		if err != nil {
			logs.CtxError(ctx, "[Handler Game EndGame] [Redis] redis lock del error: %#v, for requestKey: %#v", err, requestKey)
		}
		// 返回值写入 gctx
		gctx.JSON(http.StatusOK, resp)
	}()
	/********************************** 上述为固定流程 *********************************************/

	// step3: 解析请求数据并校检参数
	// 参数：用户临时 ID
	tempId, ok := reqData["tempId"]
	if !ok || util.IsNil(tempId) {
		logs.CtxError(ctx, "[Handler Game EndGame] reqData ( %#v ) no param `tempId`", reqData)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, "[Handler Game EndGame] reqData ( %#v ) no param `tempId`", reqData)
		return
	}
	logs.CtxInfo(ctx, "[Handler Game EndGame] reqData param `tempId`: %#v", tempId)

	// step4: 核心服务层执行请求
	// 用户服务：获取用户 openId
	openId, err := service.CoreUserService.GetUserOpenId(ctx, tempId)
	if err != nil {
		logs.CtxError(ctx, "[Handler Game EndGame] CoreUserService.GetUserOpenId err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Game EndGame] CoreUserService.GetUserOpenId resp Success")
	// 房间服务：获取房间 Id
	roomId, err := service.CoreRoomService.GetRoomId(ctx, openId)
	if err == cconst.RedisNil {
		logs.CtxError(ctx, "[Handler Game EndGame] CoreRoomService.GetRoomId RoomInvalid")
		// 包装返回信息
		resp = util.RoomInvalid(ctx)
		return
	}
	if err != nil {
		logs.CtxError(ctx, "[Handler Game EndGame] CoreRoomService.GetRoomId err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Game EndGame] CoreRoomService.GetRoomId resp Success")
	// 房间服务：获取房间信息和房间锁
	roomInfo, roomInfoUnlockFunc, err := service.CoreRoomService.GetRoomInfoAndLockFromRoomId(ctx, openId, roomId, methodKey) // TODO：cconst.RedisNil
	if err != nil {
		logs.CtxError(ctx, "[Handler Game EndGame] CoreRoomService.GetRoomInfoAndLockFromRoomId err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Game EndGame] CoreRoomService.GetRoomInfoAndLockFromRoomId resp Success")
	// 释放房间锁
	defer func() {
		unlockErr := roomInfoUnlockFunc(roomInfo, err)
		if unlockErr != nil {
			logs.CtxError(ctx, "[Handler Game EndGame] defer roomInfoUnlockFunc error: %#v", unlockErr)
			// 包装返回信息
			resp = util.ErrorWithMessage(ctx, unlockErr.Error())
			return
		}
		logs.CtxInfo(ctx, "[Handler Game EndGame] defer roomInfoUnlockFunc resp Success")
	}()
	// 游戏服务：执行逻辑
	roomInfo, err = service.CoreGameService.EndGame(ctx, openId, roomInfo)
	if err != nil {
		logs.CtxError(ctx, "[Handler Game EndGame] CoreGameService.EndGame err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Game EndGame] CoreGameService.EndGame resp Success")
	// TODO WebSocket

	/********************************** 下述为固定流程 *********************************************/
	// step5: 包装返回值
	resp = util.Success(ctx)
}

func RestartGame(gctx *gin.Context) {
	// step1: 获取请求 ctx，请求数据，请求唯一标识
	// 1.1 请求 ctx
	ctxValue, exists := gctx.Get(cconst.CtxKey)
	ctx, ok := ctxValue.(context.Context)
	if !exists || !ok {
		logs.CtxFatal(ctx, "[Handler Game RestartGame] ctx not exists in gin ctx")
		panic("[Handler Game RestartGame] ctx not exists in gin ctx")
	}
	// 1.2 请求数据
	reqDataValue, exists := gctx.Get(cconst.ReqDataKey)
	reqData, ok := reqDataValue.(map[string]string)
	if !exists || !ok {
		logs.CtxFatal(ctx, "[Handler Game RestartGame] req_data not exists in gin ctx")
		panic("[Handler Game RestartGame] req_data not exists in gin ctx")
	}
	// 1.3 请求唯一标识
	requestKey := util.CtxGetString(ctx, util.RequestKey)
	// 1.4 接口唯一标识
	methodKey := util.CtxGetString(ctx, util.MethodKey)

	// step2: 定义返回数据，并打印返回日志，删除重入锁
	var resp map[string]interface{}
	defer func() {
		// 打印返回日志
		logs.CtxInfo(ctx, "[Handler Game RestartGame] resp: %#v", resp)
		// 删除重入锁
		// Redis DEL 命令用于删除已存在的键。不存在的 key 会被忽略。
		err := model.GatewayRedis.Del(requestKey).Err()
		if err != nil {
			logs.CtxError(ctx, "[Handler Game RestartGame] [Redis] redis lock del error: %#v, for requestKey: %#v", err, requestKey)
		}
		// 返回值写入 gctx
		gctx.JSON(http.StatusOK, resp)
	}()
	/********************************** 上述为固定流程 *********************************************/

	// step3: 解析请求数据并校检参数
	// 参数：用户临时 ID
	tempId, ok := reqData["tempId"]
	if !ok || util.IsNil(tempId) {
		logs.CtxError(ctx, "[Handler Game RestartGame] reqData ( %#v ) no param `tempId`", reqData)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, "[Handler Game RestartGame] reqData ( %#v ) no param `tempId`", reqData)
		return
	}
	logs.CtxInfo(ctx, "[Handler Game RestartGame] reqData param `tempId`: %#v", tempId)

	// step4: 核心服务层执行请求
	// 用户服务：获取用户 openId
	openId, err := service.CoreUserService.GetUserOpenId(ctx, tempId)
	if err != nil {
		logs.CtxError(ctx, "[Handler Game RestartGame] CoreUserService.GetUserOpenId err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Game RestartGame] CoreUserService.GetUserOpenId resp Success")
	// 房间服务：获取房间 Id
	roomId, err := service.CoreRoomService.GetRoomId(ctx, openId)
	if err == cconst.RedisNil {
		logs.CtxError(ctx, "[Handler Game RestartGame] CoreRoomService.GetRoomId RoomInvalid")
		// 包装返回信息
		resp = util.RoomInvalid(ctx)
		return
	}
	if err != nil {
		logs.CtxError(ctx, "[Handler Game RestartGame] CoreRoomService.GetRoomId err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Game RestartGame] CoreRoomService.GetRoomId resp Success")
	// 房间服务：获取房间信息和房间锁
	roomInfo, roomInfoUnlockFunc, err := service.CoreRoomService.GetRoomInfoAndLockFromRoomId(ctx, openId, roomId, methodKey) // TODO：cconst.RedisNil
	if err != nil {
		logs.CtxError(ctx, "[Handler Game RestartGame] CoreRoomService.GetRoomInfoAndLockFromRoomId err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Game RestartGame] CoreRoomService.GetRoomInfoAndLockFromRoomId resp Success")
	// 释放房间锁
	defer func() {
		unlockErr := roomInfoUnlockFunc(roomInfo, err)
		if unlockErr != nil {
			logs.CtxError(ctx, "[Handler Game RestartGame] defer roomInfoUnlockFunc error: %#v", unlockErr)
			// 包装返回信息
			resp = util.ErrorWithMessage(ctx, unlockErr.Error())
			return
		}
		logs.CtxInfo(ctx, "[Handler Game RestartGame] defer roomInfoUnlockFunc resp Success")
	}()
	// 游戏服务：执行逻辑
	roomInfo, err = service.CoreGameService.RestartGame(ctx, openId, roomInfo)
	if err != nil {
		logs.CtxError(ctx, "[Handler Game RestartGame] CoreGameService.RestartGame err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Game RestartGame] CoreGameService.RestartGame resp Success")
	// TODO WebSocket

	/********************************** 下述为固定流程 *********************************************/
	// step5: 包装返回值
	resp = util.Success(ctx)
}
