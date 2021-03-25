package handler

import (
	"context"
	"encoding/json"

	"net/http"

	logs "lmf.mortal.com/GoLogs"

	"lmf.mortal.com/GoWxWhoIsTheSpy/cconst"
	"lmf.mortal.com/GoWxWhoIsTheSpy/model"
	"lmf.mortal.com/GoWxWhoIsTheSpy/service"
	"lmf.mortal.com/GoWxWhoIsTheSpy/util"

	"github.com/gin-gonic/gin"
)

/**
房间服务：
1. new_room：创建房间
2. enter_room：加入房间
3. exit_room：退出房间
4. refresh_room：获取房间信息
*/
type RoomHandler struct {
}

func NewRoomHandler() *RoomHandler {
	return &RoomHandler{}
}

func (h *RoomHandler) Register(engine *gin.Engine) {
	room := engine.Group("room")
	{
		room.POST("new_room", NewRoom)
		room.POST("enter_room", EnterRoom)
		room.POST("exit_room", ExitRoom)
		room.POST("refresh_room", RefreshRoom)
	}
}

func NewRoom(gctx *gin.Context) {
	// step1: 获取请求 ctx，请求数据，请求唯一标识
	// 1.1 请求 ctx
	ctxValue, exists := gctx.Get(cconst.CtxKey)
	ctx, ok := ctxValue.(context.Context)
	if !exists || !ok {
		logs.CtxFatal(ctx, "[Handler Room NewRoom] ctx not exists in gin ctx")
		panic("[Handler Room NewRoom] ctx not exists in gin ctx")
	}
	// 1.2 请求数据
	reqDataValue, exists := gctx.Get(cconst.ReqDataKey)
	reqData, ok := reqDataValue.(map[string]string)
	if !exists || !ok {
		logs.CtxFatal(ctx, "[Handler Room NewRoom] req_data not exists in gin ctx")
		panic("[Handler Room NewRoom] req_data not exists in gin ctx")
	}
	// 1.3 请求唯一标识
	requestKey := util.CtxGetString(ctx, util.RequestKey)

	// step2: 定义返回数据，并打印返回日志，删除重入锁
	var resp map[string]interface{}
	defer func() {
		// 打印返回日志
		logs.CtxInfo(ctx, "[Handler Room NewRoom] resp: %#v", resp)
		// 删除重入锁
		// Redis DEL 命令用于删除已存在的键。不存在的 key 会被忽略。
		err := model.GatewayRedis.Del(requestKey).Err()
		if err != nil {
			logs.CtxError(ctx, "[Handler Room NewRoom] [Redis] redis lock del error: %#v, for requestKey: %#v", err, requestKey)
		}
		// 返回值写入 gctx
		gctx.JSON(http.StatusOK, resp)
	}()
	/********************************** 上述为固定流程 *********************************************/

	// step3: 解析请求数据并校检参数
	// 参数：用户临时 ID
	tempId, ok := reqData["tempId"]
	if !ok || util.IsNil(tempId) {
		logs.CtxError(ctx, "[Handler Room NewRoom] reqData ( %#v ) no param `tempId`", reqData)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, "[Handler Room NewRoom] reqData ( %#v ) no param `tempId`", reqData)
		return
	}
	logs.CtxInfo(ctx, "[Handler Room NewRoom] reqData param `tempId`: %#v", tempId)
	// 参数：房间信息
	roomSettingValue, ok := reqData["roomSetting"]
	if !ok || util.IsNil(roomSettingValue) {
		logs.CtxError(ctx, "[Handler Room NewRoom] reqData ( %#v ) no param `roomSetting`", reqData)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, "[Handler Room NewRoom] reqData ( %#v ) no param `roomSetting`", reqData)
		return
	}
	logs.CtxInfo(ctx, "[Handler Room NewRoom] reqData param `roomSetting`: %#v", roomSettingValue)
	// 解析参数：房间信息
	roomSetting := &model.RoomSetting{}
	err := json.Unmarshal([]byte(roomSettingValue), roomSetting)
	if err != nil {
		logs.CtxError(ctx, "[Handler Room NewRoom] reqData ( %#v ) param `roomSetting` err: %#v", reqData, err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, "[Handler Room NewRoom] reqData ( %#v ) param `roomSetting` err: %#v", reqData, err)
		return
	}
	logs.CtxInfo(ctx, "[Handler Room NewRoom] reqData param `roomSetting`: %#v", roomSetting)

	// step4: 核心服务层执行请求
	// 用户服务：获取用户 openId
	openId, err := service.CoreUserService.GetUserOpenId(ctx, tempId)
	if err != nil {
		logs.CtxError(ctx, "[Handler Room NewRoom] CoreUserService.GetUserOpenId err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Room NewRoom] CoreUserService.GetUserOpenId resp Success")
	// 房间服务：执行逻辑
	roomId, roomInfo, err := service.CoreRoomService.NewRoom(ctx, openId, roomSetting)
	if err != nil {
		logs.CtxError(ctx, "[Handler Room NewRoom] CoreRoomService.NewRoom err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Room NewRoom] CoreRoomService.NewRoom resp Success")

	/********************************** 下述为固定流程 *********************************************/
	// step5: 包装返回值
	respData := map[string]interface{}{
		"roomId":   roomId,
		"roomInfo": service.CoreRoomService.UpdateRoomInfoWithOpenId(ctx, openId, roomInfo),// 房间服务：根据用户 openId 封装房间信息
	}
	resp = util.SuccessWithData(ctx, respData)
}

func EnterRoom(gctx *gin.Context) {
	// step1: 获取请求 ctx，请求数据，请求唯一标识
	// 1.1 请求 ctx
	ctxValue, exists := gctx.Get(cconst.CtxKey)
	ctx, ok := ctxValue.(context.Context)
	if !exists || !ok {
		logs.CtxFatal(ctx, "[Handler Room EnterRoom] ctx not exists in gin ctx")
		panic("[Handler Room EnterRoom] ctx not exists in gin ctx")
	}
	// 1.2 请求数据
	reqDataValue, exists := gctx.Get(cconst.ReqDataKey)
	reqData, ok := reqDataValue.(map[string]string)
	if !exists || !ok {
		logs.CtxFatal(ctx, "[Handler Room EnterRoom] req_data not exists in gin ctx")
		panic("[Handler Room EnterRoom] req_data not exists in gin ctx")
	}
	// 1.3 请求唯一标识
	requestKey := util.CtxGetString(ctx, util.RequestKey)
	// 1.4 接口唯一标识
	methodKey := util.CtxGetString(ctx, util.MethodKey)

	// step2: 定义返回数据，并打印返回日志，删除重入锁
	var resp map[string]interface{}
	defer func() {
		// 打印返回日志
		logs.CtxInfo(ctx, "[Handler Room EnterRoom] resp: %#v", resp)
		// 删除重入锁
		// Redis DEL 命令用于删除已存在的键。不存在的 key 会被忽略。
		err := model.GatewayRedis.Del(requestKey).Err()
		if err != nil {
			logs.CtxError(ctx, "[Handler Room EnterRoom] [Redis] redis lock del error: %#v, for requestKey: %#v", err, requestKey)
		}
		// 返回值写入 gctx
		gctx.JSON(http.StatusOK, resp)
	}()
	/********************************** 上述为固定流程 *********************************************/

	// step3: 解析请求数据并校检参数
	// 参数：用户临时 ID
	tempId, ok := reqData["tempId"]
	if !ok || util.IsNil(tempId) {
		logs.CtxError(ctx, "[Handler Room EnterRoom] reqData ( %#v ) no param `tempId`", reqData)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, "[Handler Room EnterRoom] reqData ( %#v ) no param `tempId`", reqData)
		return
	}
	logs.CtxInfo(ctx, "[Handler Room EnterRoom] reqData param `tempId`: %#v", tempId)
	// 参数：房间 ID
	roomId, ok := reqData["roomId"]
	if !ok || util.IsNil(roomId) {
		logs.CtxError(ctx, "[Handler Room EnterRoom] reqData ( %#v ) no param `roomId`", reqData)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, "[Handler Room EnterRoom] reqData ( %#v ) no param `roomId`", reqData)
		return
	}
	logs.CtxInfo(ctx, "[Handler Room EnterRoom] reqData param `roomId`: %#v", roomId)

	// step4: 核心服务层执行请求
	// 用户服务：获取用户 openId
	openId, err := service.CoreUserService.GetUserOpenId(ctx, tempId)
	if err != nil {
		logs.CtxError(ctx, "[Handler Room EnterRoom] CoreUserService.GetUserOpenId err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Room EnterRoom] CoreUserService.GetUserOpenId resp Success")
	// 房间服务：获取房间信息和房间锁
	roomInfo, roomInfoUnlockFunc, err := service.CoreRoomService.GetRoomInfoAndLockFromRoomId(ctx, "", roomId, methodKey) // TODO：cconst.RedisNil
	if err != nil {
		logs.CtxError(ctx, "[Handler Room EnterRoom] CoreRoomService.GetRoomInfoAndLockFromRoomId err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Room EnterRoom] CoreRoomService.GetRoomInfoAndLockFromRoomId resp Success")
	// 释放房间锁
	defer func() {
		unlockErr := roomInfoUnlockFunc(roomInfo, err)
		if unlockErr != nil {
			logs.CtxError(ctx, "[Handler Room EnterRoom] defer roomInfoUnlockFunc error: %#v", unlockErr)
			// 包装返回信息
			resp = util.ErrorWithMessage(ctx, unlockErr.Error())
			return
		}
		logs.CtxInfo(ctx, "[Handler Room EnterRoom] defer roomInfoUnlockFunc resp Success")
	}()
	// 房间服务：执行逻辑
	roomInfo, err = service.CoreRoomService.EnterRoom(ctx, openId, roomInfo)
	if err != nil {
		logs.CtxError(ctx, "[Handler Room EnterRoom] CoreRoomService.EnterRoom err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Room EnterRoom] CoreRoomService.EnterRoom resp Success")

	/********************************** 下述为固定流程 *********************************************/
	// step5: 包装返回值
	respData := map[string]interface{}{
		"roomInfo": service.CoreRoomService.UpdateRoomInfoWithOpenId(ctx, openId, roomInfo),// 房间服务：根据用户 openId 封装房间信息
	}
	resp = util.SuccessWithData(ctx, respData)
}

func ExitRoom(gctx *gin.Context) {
	// step1: 获取请求 ctx，请求数据，请求唯一标识
	// 1.1 请求 ctx
	ctxValue, exists := gctx.Get(cconst.CtxKey)
	ctx, ok := ctxValue.(context.Context)
	if !exists || !ok {
		logs.CtxFatal(ctx, "[Handler Room ExitRoom] ctx not exists in gin ctx")
		panic("[Handler Room ExitRoom] ctx not exists in gin ctx")
	}
	// 1.2 请求数据
	reqDataValue, exists := gctx.Get(cconst.ReqDataKey)
	reqData, ok := reqDataValue.(map[string]string)
	if !exists || !ok {
		logs.CtxFatal(ctx, "[Handler Room ExitRoom] req_data not exists in gin ctx")
		panic("[Handler Room ExitRoom] req_data not exists in gin ctx")
	}
	// 1.3 请求唯一标识
	requestKey := util.CtxGetString(ctx, util.RequestKey)
	// 1.4 接口唯一标识
	methodKey := util.CtxGetString(ctx, util.MethodKey)

	// step2: 定义返回数据，并打印返回日志，删除重入锁
	var resp map[string]interface{}
	defer func() {
		// 打印返回日志
		logs.CtxInfo(ctx, "[Handler Room ExitRoom] resp: %#v", resp)
		// 删除重入锁
		// Redis DEL 命令用于删除已存在的键。不存在的 key 会被忽略。
		err := model.GatewayRedis.Del(requestKey).Err()
		if err != nil {
			logs.CtxError(ctx, "[Handler Room ExitRoom] [Redis] redis lock del error: %#v, for requestKey: %#v", err, requestKey)
		}
		// 返回值写入 gctx
		gctx.JSON(http.StatusOK, resp)
	}()
	/********************************** 上述为固定流程 *********************************************/

	// step3: 解析请求数据并校检参数
	// 参数：用户临时 ID
	tempId, ok := reqData["tempId"]
	if !ok || util.IsNil(tempId) {
		logs.CtxError(ctx, "[Handler Room ExitRoom] reqData ( %#v ) no param `tempId`", reqData)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, "[Handler Room ExitRoom] reqData ( %#v ) no param `tempId`", reqData)
		return
	}
	logs.CtxInfo(ctx, "[Handler Room ExitRoom] reqData param `tempId`: %#v", tempId)

	// step4: 核心服务层执行请求
	// 用户服务：获取用户 openId
	openId, err := service.CoreUserService.GetUserOpenId(ctx, tempId)
	if err != nil {
		logs.CtxError(ctx, "[Handler Room ExitRoom] CoreUserService.GetUserOpenId err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Room ExitRoom] CoreUserService.GetUserOpenId resp Success")
	// 房间服务：获取房间 Id
	roomId, err := service.CoreRoomService.GetRoomId(ctx, openId)
	if err == cconst.RedisNil {
		logs.CtxError(ctx, "[Handler Room ExitRoom] CoreRoomService.GetRoomId RoomInvalid")
		// 包装返回信息
		resp = util.RoomInvalid(ctx)
		return
	}
	if err != nil {
		logs.CtxError(ctx, "[Handler Room ExitRoom] CoreRoomService.GetRoomId err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	// 房间服务：获取房间信息和房间锁
	roomInfo, roomInfoUnlockFunc, err := service.CoreRoomService.GetRoomInfoAndLockFromRoomId(ctx, "", roomId, methodKey) // TODO：cconst.RedisNil
	if err != nil {
		logs.CtxError(ctx, "[Handler Room ExitRoom] CoreRoomService.GetRoomInfoAndLockFromRoomId err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Room ExitRoom] CoreRoomService.GetRoomInfoAndLockFromRoomId resp Success")
	// 释放房间锁
	defer func() {
		unlockErr := roomInfoUnlockFunc(roomInfo, err)
		if unlockErr != nil {
			logs.CtxError(ctx, "[Handler Room ExitRoom] defer roomInfoUnlockFunc error: %#v", unlockErr)
			// 包装返回信息
			resp = util.ErrorWithMessage(ctx, unlockErr.Error())
			return
		}
		logs.CtxInfo(ctx, "[Handler Room ExitRoom] defer roomInfoUnlockFunc resp Success")
	}()
	// 房间服务：执行逻辑
	roomInfo, err = service.CoreRoomService.ExitRoom(ctx, openId, roomInfo)
	if err != nil {
		logs.CtxError(ctx, "[Handler Room ExitRoom] CoreRoomService.ExitRoom err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Room ExitRoom] CoreRoomService.ExitRoom resp Success")

	/********************************** 下述为固定流程 *********************************************/
	// step5: 包装返回值
	resp = util.Success(ctx)
}

func RefreshRoom(gctx *gin.Context) {
	// step1: 获取请求 ctx，请求数据，请求唯一标识
	// 1.1 请求 ctx
	ctxValue, exists := gctx.Get(cconst.CtxKey)
	ctx, ok := ctxValue.(context.Context)
	if !exists || !ok {
		logs.CtxFatal(ctx, "[Handler Room RefreshRoom] ctx not exists in gin ctx")
		panic("[Handler Room RefreshRoom] ctx not exists in gin ctx")
	}
	// 1.2 请求数据
	reqDataValue, exists := gctx.Get(cconst.ReqDataKey)
	reqData, ok := reqDataValue.(map[string]string)
	if !exists || !ok {
		logs.CtxFatal(ctx, "[Handler Room RefreshRoom] req_data not exists in gin ctx")
		panic("[Handler Room RefreshRoom] req_data not exists in gin ctx")
	}
	// 1.3 请求唯一标识
	requestKey := util.CtxGetString(ctx, util.RequestKey)

	// step2: 定义返回数据，并打印返回日志，删除重入锁
	var resp map[string]interface{}
	defer func() {
		// 打印返回日志
		logs.CtxInfo(ctx, "[Handler Room RefreshRoom] resp: %#v", resp)
		// 删除重入锁
		// Redis DEL 命令用于删除已存在的键。不存在的 key 会被忽略。
		err := model.GatewayRedis.Del(requestKey).Err()
		if err != nil {
			logs.CtxError(ctx, "[Handler Room RefreshRoom] [Redis] redis lock del error: %#v, for requestKey: %#v", err, requestKey)
		}
		// 返回值写入 gctx
		gctx.JSON(http.StatusOK, resp)
	}()
	/********************************** 上述为固定流程 *********************************************/

	// step3: 解析请求数据并校检参数
	// 参数：用户临时 ID
	tempId, ok := reqData["tempId"]
	if !ok || util.IsNil(tempId) {
		logs.CtxError(ctx, "[Handler Room RefreshRoom] reqData ( %#v ) no param `tempId`", reqData)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, "[Handler Room RefreshRoom] reqData ( %#v ) no param `tempId`", reqData)
		return
	}
	logs.CtxInfo(ctx, "[Handler Room RefreshRoom] reqData param `tempId`: %#v", tempId)

	// step4: 核心服务层执行请求
	// 用户服务：获取用户 openId
	openId, err := service.CoreUserService.GetUserOpenId(ctx, tempId)
	if err != nil {
		logs.CtxError(ctx, "[Handler Room RefreshRoom] CoreUserService.GetUserOpenId err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Room RefreshRoom] CoreUserService.GetUserOpenId resp Success")
	// 房间服务：执行逻辑
	roomInfo, err := service.CoreRoomService.RefreshRoom(ctx, openId)
	if err == cconst.RedisNil { // 用户所在房间已失效，报错 RoomInvalid
		logs.CtxError(ctx, "[Handler Room RefreshRoom] CoreRoomService.RefreshRoom RoomInvalid")
		// 包装返回信息
		resp = util.RoomInvalid(ctx)
		return
	}
	if err != nil {
		logs.CtxError(ctx, "[Handler Room RefreshRoom] CoreRoomService.RefreshRoom err: %#v", err)
		// 包装返回信息
		resp = util.ErrorWithMessage(ctx, err.Error())
		return
	}
	logs.CtxInfo(ctx, "[Handler Room RefreshRoom] CoreRoomService.RefreshRoom resp Success")

	/********************************** 下述为固定流程 *********************************************/
	// step5: 包装返回值
	respData := map[string]interface{}{
		"roomInfo": service.CoreRoomService.UpdateRoomInfoWithOpenId(ctx, openId, roomInfo),// 房间服务：根据用户 openId 封装房间信息
	}
	resp = util.SuccessWithData(ctx, respData)
}
