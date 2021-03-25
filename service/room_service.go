package service

import (
	"context"
	"math/rand"
	"sync"
	"time"

	logs "lmf.mortal.com/GoLogs"

	"lmf.mortal.com/GoWxWhoIsTheSpy/cconst"
	"lmf.mortal.com/GoWxWhoIsTheSpy/model"
	"lmf.mortal.com/GoWxWhoIsTheSpy/util"
)

// 房间服务
type RoomService struct {
	userDao *model.UserDao // 访问用户信息
	roomDao *model.RoomDao // 访问房间信息
}

var roomServiceInstance *RoomService  // 房间服务实例
var roomServiceInstanceOnce sync.Once // 单例模式

// 获取房间服务单例
func RoomServiceInstance() *RoomService {
	roomServiceInstanceOnce.Do(func() {
		roomServiceInstance = &RoomService{
			userDao: model.UserDaoInstance(),
			roomDao: model.RoomDaoInstance(),
		}
	})
	return roomServiceInstance
}

/**
内部接口：由用户 openId 新建玩家信息

req:
- openId: 用户 openId

resp:
- player: 玩家信息
*/
func (s *RoomService) createPlayerFromUserOpenId(ctx context.Context, openId string) (player *model.Player, err error) {
	logs.CtxInfo(ctx, "[Service RoomService createPlayerFromUserOpenId Req] req: %s", openId) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Service RoomService createPlayerFromUserOpenId Resp] resp: %#v, %#v", player, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	// step1，获取用户信息
	userInfo, err := s.userDao.QueryUser(ctx, openId)
	if err != nil {
		logs.CtxError(ctx, "[Service RoomService createPlayerFromUserOpenId] QueryUser error: %#v", err)
		err = util.NewErrf("[Service RoomService createPlayerFromUserOpenId] QueryUser error: %#v", err)
		return nil, err
	}

	// step2，构建玩家信息
	player = &model.Player{
		OpenId:    openId,
		NickName:  userInfo.NickName,
		AvatarUrl: userInfo.AvatarUrl,
		State:     model.PlayerStateWait,
		Word:      "", // 游戏开始时分配
		Role:      "", // 游戏开始时分配
	}
	/************************************ 核心逻辑 ****************************************/

	return player, nil
}

/**
内部接口：由房间配置新建房间信息

req:
- masterOpenId: 房主 openId
- roomSetting: 房间配置信息

resp:
- roomInfo: 房间信息
*/
func (s *RoomService) createRoomInfoFromRoomSetting(ctx context.Context, masterOpenId string, roomSetting *model.RoomSetting) (roomInfo *model.RoomInfo, err error) {
	logs.CtxInfo(ctx, "[Service RoomService createRoomInfoFromRoomSetting Req] req: %s, %#v", masterOpenId, roomSetting) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Service RoomService createRoomInfoFromRoomSetting Resp] resp: %s, %#v", roomInfo, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	// step1，检验房间配置是否合法
	if roomSetting.TotalNum < 2 || // 总人数应该大于 2
		roomSetting.SpyNum < 1 || // 卧底人数应该大于 1
		roomSetting.BlankNum < 0 || // 白板人数应该不小于 0
		roomSetting.TotalNum <= roomSetting.SpyNum+roomSetting.BlankNum { // 平民人数应该大于 0
		logs.CtxError(ctx, "[Service RoomService createRoomInfoFromRoomSetting] roomSetting invalid: %#v", roomSetting)
		err = util.NewErrf("[Service RoomService createRoomInfoFromRoomSetting] roomSetting invalid: %#v", roomSetting)
		return nil, err
	}
	// step2，构建玩家信息
	player, err := s.createPlayerFromUserOpenId(ctx, masterOpenId)
	if err != nil {
		logs.CtxError(ctx, "[Service RoomService createRoomInfoFromRoomSetting] createPlayerFromUserOpenId error: %#v", err)
		err = util.NewErrf("[Service RoomService createRoomInfoFromRoomSetting] createPlayerFromUserOpenId error: %#v", err)
		return nil, err
	}
	player.State = model.PlayerStateReady // 房主默认已准备

	// step3，构建房间信息
	roomInfo = &model.RoomInfo{
		RoomSetting:  roomSetting,
		MasterOpenId: masterOpenId,
		PlayerList:   []*model.Player{player}, // 只有房主一个玩家
		State:        model.RoomStateOpen,     // 创建时房间为开放状态
		Word:         &model.Word{},           // 游戏开始时分配
		WordCache:    []*model.Word{},
	}
	/************************************ 核心逻辑 ****************************************/

	return roomInfo, nil
}

/**
外部接口：新建房间

req:
- masterOpenId: 房主 openId
- roomSetting: 房间配置信息

resp:
- roomId: 房间 ID
- roomInfo: 房间信息
*/
func (s *RoomService) NewRoom(ctx context.Context, masterOpenId string, roomSetting *model.RoomSetting) (roomId string, roomInfo *model.RoomInfo, err error) {
	logs.CtxInfo(ctx, "[Service RoomService NewRoom Req] req: %s, %#v", masterOpenId, roomSetting) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Service RoomService NewRoom Resp] resp: %s, %s, %#v", roomId, roomInfo, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	// step1，构建房间信息
	roomInfo, err = s.createRoomInfoFromRoomSetting(ctx, masterOpenId, roomSetting)
	if err != nil {
		logs.CtxError(ctx, "[Service RoomService NewRoom] createRoomInfoFromRoomSetting error: %#v", err)
		err = util.NewErrf("[Service RoomService NewRoom] createRoomInfoFromRoomSetting error: %#v", err)
		return "", nil, err
	}
	// step2，存储房间信息到 redis，获取房间 ID
	roomId, err = s.roomDao.SetRoomInfoGenID(ctx, roomInfo)
	if err != nil {
		logs.CtxError(ctx, "[Service RoomService NewRoom] roomDao.SetRoomInfoGenID error: %#v", err)
		err = util.NewErrf("[Service RoomService NewRoom] roomDao.SetRoomInfoGenID error: %#v", err)
		return "", nil, err
	}
	// step3，构建房间 ID 和用户 openId 关联关系
	err = s.roomDao.SetRoomID(ctx, masterOpenId, roomId)
	if err != nil {
		logs.CtxError(ctx, "[Service RoomService NewRoom] roomDao.SetRoomID error: %#v", err)
		err = util.NewErrf("[Service RoomService NewRoom] roomDao.SetRoomID error: %#v", err)
		return "", nil, err
	}

	/************************************ 核心逻辑 ****************************************/

	return roomId, roomInfo, nil
}

/**
内部接口：房间信息新增玩家

req:
- playerOpenId: 玩家 openId
- roomInfoBefore: 房间信息（玩家加入前）

resp:
- roomInfoAfter: 房间信息（玩家加入后）
*/
func (s *RoomService) roomInfoAddPlayer(ctx context.Context, playerOpenId string, roomInfoBefore *model.RoomInfo) (roomInfoAfter *model.RoomInfo, err error) {
	logs.CtxInfo(ctx, "[Service RoomService roomInfoAddPlayer Req] req: %s, %s", playerOpenId, roomInfoBefore) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Service RoomService roomInfoAddPlayer Resp] resp: %s, %#v", roomInfoAfter, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	// step1，判断玩家是否已在房间中
	for _, player := range roomInfoBefore.PlayerList {
		if playerOpenId == player.OpenId {
			logs.CtxError(ctx, "[Service RoomService roomInfoAddPlayer] createPlayerFromUserOpenId error: user already exists")
			err = util.NewErrf("[Service RoomService roomInfoAddPlayer] createPlayerFromUserOpenId error: user already exists")
			return nil, err
		}
	}
	// step2，构建玩家信息
	player, err := s.createPlayerFromUserOpenId(ctx, playerOpenId)
	if err != nil {
		logs.CtxError(ctx, "[Service RoomService roomInfoAddPlayer] createPlayerFromUserOpenId error: %#v", err)
		err = util.NewErrf("[Service RoomService roomInfoAddPlayer] createPlayerFromUserOpenId error: %#v", err)
		return nil, err
	}
	// step3，判断房间是否可加入，即房间为开放状态
	if roomInfoBefore.State != model.RoomStateOpen {
		logs.CtxError(ctx, "[Service RoomService roomInfoAddPlayer] room can not add player, room state: %#v", roomInfoBefore.State)
		err = util.NewErrf("[Service RoomService roomInfoAddPlayer] room can not add player, room state: %#v", roomInfoBefore.State)
		return nil, err
	}
	// step4，玩家加入房间
	roomInfoBefore.PlayerList = append(roomInfoBefore.PlayerList, player)
	// step5，如果人数已满，更改房间状态为等待玩家准备
	if len(roomInfoBefore.PlayerList) == roomInfoBefore.TotalNum {
		roomInfoBefore.State = model.RoomStateWait
	}
	// step6，赋值更改后的房间信息
	roomInfoAfter = roomInfoBefore
	/************************************ 核心逻辑 ****************************************/

	return roomInfoAfter, nil
}

/**
外部接口：加入房间

req:
- openId: 用户 openId
- roomId: 房间 ID

resp:
- roomInfo: 房间信息
*/
func (s *RoomService) EnterRoom(ctx context.Context, openId string, roomInfoBefore *model.RoomInfo) (roomInfoAfter *model.RoomInfo, err error) {
	logs.CtxInfo(ctx, "[Service RoomService EnterRoom Req] req: %s, %s", openId, roomInfoBefore) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Service RoomService EnterRoom Resp] resp: %s, %#v", roomInfoAfter, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	// step1，房间中加入玩家
	roomInfoAfter, err = s.roomInfoAddPlayer(ctx, openId, roomInfoBefore)
	if err != nil {
		logs.CtxError(ctx, "[Service RoomService EnterRoom] roomInfoAddPlayer error: %#v", err)
		err = util.NewErrf("[Service RoomService EnterRoom] roomInfoAddPlayer error: %#v", err)
		return nil, err
	}
	// step2，玩家建立和房间的关系
	err = s.roomDao.SetRoomID(ctx, openId, roomInfoBefore.RoomId)
	if err != nil {
		logs.CtxError(ctx, "[Service RoomService EnterRoom] roomDao.SetRoomID error: %#v", err)
		err = util.NewErrf("[Service RoomService EnterRoom] roomDao.SetRoomID error: %#v", err)
		return nil, err
	}
	/************************************ 核心逻辑 ****************************************/

	return roomInfoAfter, nil
}

/**
内部接口：房间信息移除玩家

req:
- playerOpenId: 玩家 openId
- roomInfoBefore: 房间信息（玩家移除前）

resp:
- roomInfoAfter: 房间信息（玩家移除后）
*/
func (s *RoomService) roomInfoRemovePlayer(ctx context.Context, playerOpenId string, roomInfoBefore *model.RoomInfo) (roomInfoAfter *model.RoomInfo, err error) {
	logs.CtxInfo(ctx, "[Service RoomService roomInfoRemovePlayer Req] req: %s, %s", playerOpenId, roomInfoBefore) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Service RoomService roomInfoRemovePlayer Resp] resp: %s, %#v", roomInfoAfter, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	// step1，判断玩家是否在房间中
	var index = 0 // 玩家索引
	for _, player := range roomInfoBefore.PlayerList {
		if player.OpenId == playerOpenId { // 玩家在房间中，找到玩家的索引
			break
		}
		index++
	}
	if index == len(roomInfoBefore.PlayerList) { // 玩家不在房间中
		logs.CtxError(ctx, "[Service RoomService roomInfoRemovePlayer] room has no player: %#v", playerOpenId)
		err = util.NewErrf("[Service RoomService roomInfoRemovePlayer] room has no player: %#v", playerOpenId)
		return nil, err
	}
	// step2，判断玩家状态，已准备的不可退出
	if roomInfoBefore.PlayerList[index].State != model.PlayerStateWait {
		logs.CtxError(ctx, "[Service RoomService roomInfoRemovePlayer] player state error: %#v", roomInfoBefore.PlayerList[index].State)
		err = util.NewErrf("[Service RoomService roomInfoRemovePlayer] player state error: %#v", roomInfoBefore.PlayerList[index].State)
		return nil, err
	}
	// step3，判断房间是否可退出，即房间为开放状态或等待状态
	if roomInfoBefore.State != model.RoomStateOpen && roomInfoBefore.State != model.RoomStateWait {
		logs.CtxError(ctx, "[Service RoomService roomInfoRemovePlayer] room can not remove player, room state: %#v", roomInfoBefore.State)
		err = util.NewErrf("[Service RoomService roomInfoRemovePlayer] room can not remove player, room state: %#v", roomInfoBefore.State)
		return nil, err
	}
	// step4，移除玩家
	//if index == 0 {
	//	roomInfoBefore.PlayerList = roomInfoBefore.PlayerList[1:] // 将待删除玩家前后的索引连接起来
	//} else {
	roomInfoBefore.PlayerList = append(roomInfoBefore.PlayerList[:index], roomInfoBefore.PlayerList[index+1:]...) // 将待删除玩家前后的索引连接起来
	//}
	// step5，更改房间状态为开放状态
	roomInfoBefore.State = model.RoomStateOpen
	// step6，房主退出则转移房主
	if roomInfoBefore.MasterOpenId == playerOpenId && len(roomInfoBefore.PlayerList) > 0 {
		roomInfoBefore.MasterOpenId = roomInfoBefore.PlayerList[0].OpenId
	}
	// step7，赋值更改后的房间信息
	roomInfoAfter = roomInfoBefore
	/************************************ 核心逻辑 ****************************************/

	return roomInfoAfter, nil
}

/**
外部接口：退出房间

req:
- openId: 用户 openId

resp:
*/
func (s *RoomService) ExitRoom(ctx context.Context, openId string, roomInfoBefore *model.RoomInfo) (roomInfoAfter *model.RoomInfo, err error) {
	logs.CtxInfo(ctx, "[Service RoomService ExitRoom Req] req: %s, %s", openId, roomInfoBefore) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Service RoomService ExitRoom Resp] resp: %s, %#v", roomInfoAfter, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	// step1，房间中移除玩家
	roomInfoAfter, err = s.roomInfoRemovePlayer(ctx, openId, roomInfoBefore)
	if err != nil {
		logs.CtxError(ctx, "[Service RoomService ExitRoom] roomInfoRemovePlayer error: %#v", err)
		err = util.NewErrf("[Service RoomService ExitRoom] roomInfoRemovePlayer error: %#v", err)
		return nil, err
	}
	// step2，删除房间关联关系
	err = s.roomDao.DelRoomID(ctx, openId)
	if err != nil {
		logs.CtxError(ctx, "[Service RoomService ExitRoom] s.roomDao.DelRoomID error: %#v", err)
		err = util.NewErrf("[Service RoomService ExitRoom] s.roomDao.DelRoomID error: %#v", err)
		return nil, err
	}
	/************************************ 核心逻辑 ****************************************/

	return roomInfoAfter, nil
}

/**
外部接口：获取房间信息

req:
- openId: 用户 openId

resp:
- roomInfo: 房间信息
*/
func (s *RoomService) RefreshRoom(ctx context.Context, openId string) (roomInfo *model.RoomInfo, err error) {
	logs.CtxInfo(ctx, "[Service RoomService RefreshRoom Req] req: %s", openId) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Service RoomService RefreshRoom Resp] resp: %s, %#v", roomInfo, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	// step1，获取用户所在房间 ID
	var roomId string
	roomId, err = s.GetRoomId(ctx, openId)
	if err == cconst.RedisNil { // 用户无房间，不报错
		logs.CtxInfo(ctx, "[Service RoomService RefreshRoom] player not in room")
		return nil, nil
	}
	if err != nil {
		logs.CtxError(ctx, "[Service RoomService RefreshRoom] GetRoomId error: %#v", err)
		err = util.NewErrf("[Service RoomService RefreshRoom] GetRoomId error: %#v", err)
		return nil, err
	}
	// step2，获取房间信息
	roomInfo, err = s.getRoomInfo(ctx, openId, roomId)
	if err == cconst.RedisNil { // 用户无房间
		logs.CtxInfo(ctx, "[Service RoomService RefreshRoom] room invalid")
		return nil, cconst.RedisNil
	}
	if err != nil {
		logs.CtxError(ctx, "[Service RoomService RefreshRoom] getRoomInfo error: %#v", err)
		err = util.NewErrf("[Service RoomService RefreshRoom] getRoomInfo error: %#v", err)
		return nil, err
	}
	//// step1，获取用户所在房间 ID
	//var roomId string
	//roomId, err = s.roomDao.GetRoomID(ctx, openId)
	//if err == cconst.RedisNil { // 用户无房间
	//	logs.CtxInfo(ctx, "[Service RoomService RefreshRoom] player not in room")
	//	return nil, cconst.RedisNil
	//}
	//if err != nil {
	//	logs.CtxError(ctx, "[Service RoomService RefreshRoom] GetRoomID error: %#v", err)
	//	err = util.NewErrf("[Service RoomService RefreshRoom] GetRoomID error: %#v", err)
	//	return nil, err
	//}
	//// step2，获取房间信息
	//roomInfo, err = s.roomDao.GetRoomInfo(ctx, roomId)
	//if err == cconst.RedisNil { // 房间已不存在，删除用户和房间的关联关系
	//	logs.CtxInfo(ctx, "[Service RoomService RefreshRoom] Room not exists, del openId -> roomId")
	//	err = s.roomDao.DelRoomID(ctx, openId) // 删除用户和房间的关联关系
	//	if err != nil {                        // 删除用户和房间的关联关系错误
	//		logs.CtxError(ctx, "[Service RoomService RefreshRoom] DelRoomID error: %#v", err)
	//		err = util.NewErrf("[Service RoomService RefreshRoom] DelRoomID error: %#v", err)
	//		return nil, err
	//	}
	//	return nil, cconst.RedisNil // 删除用户和房间的关联关系成功，返回房间无效错误信息
	//}
	//if err != nil { // 获取房间错误
	//	logs.CtxError(ctx, "[Service RoomService RefreshRoom] GetRoomInfo error: %#v", err)
	//	err = util.NewErrf("[Service RoomService RefreshRoom] GetRoomInfo error: %#v", err)
	//	return nil, err
	//}
	/************************************ 核心逻辑 ****************************************/

	return roomInfo, nil
}

/**
内部接口：获取房间信息，若房间已不存在，删除用户和房间的关联关系

req:
- openId: 用户 openId（删除用户和房间的关联关系时使用）（ openId 可能为空，此时只是单独查询房间信息）
- roomId: 房间 ID

resp:
- roomInfo: 房间信息
*/
func (s *RoomService) getRoomInfo(ctx context.Context, openId, roomId string) (roomInfo *model.RoomInfo, err error) {
	logs.CtxInfo(ctx, "[Service RoomService getRoomInfo Req] req: %s, %s", openId, roomId) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Service RoomService getRoomInfo Resp] resp: %s, %#v", roomInfo, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	// step1，获取房间信息
	roomInfo, err = s.roomDao.GetRoomInfo(ctx, roomId)
	if err == cconst.RedisNil { // 房间已不存在，删除用户和房间的关联关系
		if !util.IsNil(openId) { // openId 可能为空，此时只是单独查询房间信息
			logs.CtxInfo(ctx, "[Service RoomService getRoomInfo] Room not exists, del openId -> roomId")
			err = s.roomDao.DelRoomID(ctx, openId) // 删除用户和房间的关联关系
			if err != nil {                        // 删除用户和房间的关联关系错误
				logs.CtxError(ctx, "[Service RoomService getRoomInfo] DelRoomID error: %#v", err)
				err = util.NewErrf("[Service RoomService getRoomInfo] DelRoomID error: %#v", err)
				return nil, err
			}
		}
		return nil, cconst.RedisNil // 删除用户和房间的关联关系成功，返回房间无效错误信息
	}
	if err != nil { // 获取房间错误
		logs.CtxError(ctx, "[Service RoomService getRoomInfo] GetRoomInfo error: %#v", err)
		err = util.NewErrf("[Service RoomService getRoomInfo] GetRoomInfo error: %#v", err)
		return nil, err
	}
	/************************************ 核心逻辑 ****************************************/

	return roomInfo, nil
}

/**
外部接口：获取房间 Id

req:
- openId: 用户 openId

resp:
- roomId: 房间 Id
*/
func (s *RoomService) GetRoomId(ctx context.Context, openId string) (roomId string, err error) {
	logs.CtxInfo(ctx, "[Service RoomService GetRoomId Req] req: %s", openId) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Service RoomService GetRoomId Resp] resp: %#v, %#v", roomId, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	// step1，获取用户所在房间 ID
	roomId, err = s.roomDao.GetRoomID(ctx, openId)
	if err == cconst.RedisNil { // 用户无房间
		logs.CtxInfo(ctx, "[Service RoomService GetRoomId] player not in room")
		return "", cconst.RedisNil
	}
	if err != nil {
		logs.CtxError(ctx, "[Service RoomService GetRoomId] GetRoomID error: %#v", err)
		err = util.NewErrf("[Service RoomService GetRoomId] GetRoomID error: %#v", err)
		return "", err
	}
	/************************************ 核心逻辑 ****************************************/

	return roomId, nil
}

/**
外部接口：获取房间信息和房间锁

req:
- openId: openId（可选参数，若不为空，在房间失效时删除 openId->roomId 的关联关系）
- roomId: 房间 ID（必须为 openId 关联的 roomId）
- serviceName: 服务名

resp:
- roomInfo: 房间信息
*/
func (s *RoomService) GetRoomInfoAndLockFromRoomId(ctx context.Context, openId, roomId, serviceName string) (roomInfo *model.RoomInfo, roomInfoUnlockFunc func(roomInfo *model.RoomInfo, roomInfoErr error) (err error), err error) {
	logs.CtxInfo(ctx, "[Service RoomService GetRoomInfoAndLockFromRoomId Req] req: %s, %s, %s", openId, roomId, serviceName) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Service RoomService GetRoomInfoAndLockFromRoomId Resp] resp: %s, %#v, %#v",roomInfo, roomInfoUnlockFunc, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	// step0，默认 roomInfoUnlockFunc，更新房间信息后解锁
	roomInfoUnlockFunc = func(roomInfo *model.RoomInfo, roomInfoErr error) (err error) {
		logs.CtxInfo(ctx, "[Service RoomService GetRoomInfoAndLockFromRoomId] roomInfoUnlockFunc, Lock Failed") // 加锁失败，无需解锁，也无需更新房间信息
		return nil
	}
	// step1，获取用户所在房间 ID
	// step2，记录加锁时间和解锁时间，预估每次加锁时长
	var lockTime, unlockTime = time.Now(), time.Now() // 记录加锁时间和解锁时间，预估每次加锁时长
	// step3，尝试获取房间信息锁，重复多次，每次失败后等待若干时间后重试
	for i := 0; i < model.RoomInfoLockRetryTimes; i++ {
		// step3.1，尝试获取房间信息锁
		err = s.roomDao.LockRoomInfo(ctx, roomId)
		if err == cconst.RedisLockDuplicate { // 获取房间信息锁失败，等待若干时间后重试
			time.Sleep(time.Duration(rand.Int63n(model.RoomInfoLockWaitTime)) * time.Millisecond)
			logs.CtxError(ctx, "[Service RoomService GetRoomInfoAndLockFromRoomId] LockRoomInfo Duplicate")
			err = util.NewErrf("[Service RoomService GetRoomInfoAndLockFromRoomId] LockRoomInfo Duplicate")
			continue
		}
		if err != nil {
			logs.CtxError(ctx, "[Service RoomService GetRoomInfoAndLockFromRoomId] LockRoomInfo error: %#v", err)
			err = util.NewErrf("[Service RoomService GetRoomInfoAndLockFromRoomId] LockRoomInfo error: %#v", err)
			return nil, roomInfoUnlockFunc, err
		}
		// step3.2，获取房间信息锁成功，定义 defer 解锁函数
		logs.CtxInfo(ctx, "[Service RoomService GetRoomInfoAndLockFromRoomId] LockRoomInfo Success")
		lockTime = time.Now() // 记录加锁时间，预估每次加锁时长
		roomInfoUnlockFunc = func(roomInfo *model.RoomInfo, roomInfoErr error) (err error) {
			// step1，加锁成功，需要解锁
			logs.CtxInfo(ctx, "[Service RoomService GetRoomInfoAndLockFromRoomId] roomInfoUnlockFunc, Lock Success")
			// step2，无论是否更新房间信息成功，都需要解锁
			defer func() {
				err = s.roomDao.UnLockRoomInfo(ctx, roomId)
				if err != nil { // 解锁失败
					logs.CtxError(ctx, "[Service RoomService GetRoomInfoAndLockFromRoomId] UnLockRoomInfo error: %#v", err)
					err = util.NewErrf("[Service RoomService GetRoomInfoAndLockFromRoomId] UnLockRoomInfo error: %#v", err)
				}
				// 解锁成功
				logs.CtxInfo(ctx, "[Service RoomService GetRoomInfoAndLockFromRoomId] roomInfoUnlockFunc, Unlock Success")
				unlockTime = time.Now()              // 记录解锁时间，预估每次加锁时长
				costTime := unlockTime.Sub(lockTime) // 耗时
				// TODO：记录加锁耗时便于后续优化
				logs.CtxInfo(ctx, "[Service RoomService GetRoomInfoAndLockFromRoomId(%s)] LockRoomInfo to UnLockRoomInfo cost time: %#v ms", serviceName, costTime/time.Millisecond)
			}()
			// step3，更新房间信息
			if roomInfoErr != nil || roomInfo == nil || util.IsNil(roomInfo.RoomId) { // 若出现错误则不更新房间信息
				logs.CtxError(ctx, "[Service GameService GetRoomInfoAndLockFromRoomId(%s)] roomInfo error: %#v", serviceName, roomInfoErr)
				err = util.NewErrf("[Service GameService GetRoomInfoAndLockFromRoomId(%s)] roomInfo error: %#v", serviceName, roomInfoErr)
				return err
			}
			if len(roomInfo.PlayerList) == 0 { // 若玩家都退出则删除房间信息
				err = s.roomDao.DelRoomInfo(ctx, roomId)
				if err != nil {
					logs.CtxError(ctx, "[Service GameService GetRoomInfoAndLockFromRoomId(%s)] roomDao.DelRoomInfo error: %#v", err)
					err = util.NewErrf("[Service GameService GetRoomInfoAndLockFromRoomId(%s)] roomDao.DelRoomInfo error: %#v", err)
					return err
				}
			} else { // 否则更新房间信息
				err = s.roomDao.SetRoomInfo(ctx, roomInfo.RoomId, roomInfo)
				if err != nil {
					logs.CtxError(ctx, "[Service GameService GetRoomInfoAndLockFromRoomId(%s)] roomDao.SetRoomInfo error: %#v", err)
					err = util.NewErrf("[Service GameService GetRoomInfoAndLockFromRoomId(%s)] roomDao.SetRoomInfo error: %#v", err)
					return err
				}
			}
			return nil
		}
		// step3.3，获取房间信息
		roomInfo, err = s.getRoomInfo(ctx, openId, roomId)
		if err != nil {
			logs.CtxError(ctx, "[Service RoomService GetRoomInfoAndLockFromRoomId] getRoomInfo error: %#v", err)
			err = util.NewErrf("[Service RoomService GetRoomInfoAndLockFromRoomId] getRoomInfo error: %#v", err)
			return nil, roomInfoUnlockFunc, err
		}
		// step4，成功
		return roomInfo, roomInfoUnlockFunc, nil
	}
	/************************************ 核心逻辑 ****************************************/

	// step4，失败
	return nil, roomInfoUnlockFunc, err
}

/**
内部接口：根据用户 openId 封装房间信息

req:
- openId: 用户 openId
- roomInfoBefore: 房间信息

resp:
- roomInfoAfter: 房间信息
*/
func (s *RoomService) UpdateRoomInfoWithOpenId(ctx context.Context, openId string, roomInfoBefore *model.RoomInfo) (roomInfoAfter *model.RoomInfo) {
	logs.CtxInfo(ctx, "[Service RoomService UpdateRoomInfoWithOpenId Req] req: %s, %s", openId, roomInfoBefore) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Service RoomService UpdateRoomInfoWithOpenId Resp] resp: %s, %#v", roomInfoAfter) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	if roomInfoBefore == nil {
		return nil
	}
	// step1，设置房主信息
	if openId == roomInfoBefore.MasterOpenId {
		roomInfoBefore.IsMaster = true
	}
	// step2，设置玩家信息
	for _, player := range roomInfoBefore.PlayerList {
		if openId == player.OpenId {
			player.IsSelf = true
			break
		}
	}
	// step3，去除敏感信息，TODO，直接修改指针会有问题，去除 enter_room 的 room_info 返回值后，lock 和 UpdateRoomInfoWithOpenId 不能同时调用，会覆盖掩码到数据库
	//roomInfoBefore.MasterOpenId = "***"
	//for _, player := range roomInfoBefore.PlayerList {
	//	player.OpenId = "***"
	//}
	// step4，覆盖更新信息
	roomInfoAfter = roomInfoBefore
	/************************************ 核心逻辑 ****************************************/

	return roomInfoAfter
}
