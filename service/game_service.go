package service

import (
	"context"
	"math/rand"
	"sync"

	logs "lmf.mortal.com/GoLogs"

	"lmf.mortal.com/GoWxWhoIsTheSpy/model"
	"lmf.mortal.com/GoWxWhoIsTheSpy/util"
)

// 游戏服务
type GameService struct {
	userDao *model.UserDao // 访问用户信息
	roomDao *model.RoomDao // 访问房间信息
	wordDao *model.WordDao // 访问词汇信息
}

var gameServiceInstance *GameService  // 游戏服务实例
var gameServiceInstanceOnce sync.Once // 单例模式

// 获取游戏服务单例
func GameServiceInstance() *GameService {
	gameServiceInstanceOnce.Do(func() {
		gameServiceInstance = &GameService{
			userDao: model.UserDaoInstance(),
			roomDao: model.RoomDaoInstance(),
			wordDao: model.WordDaoInstance(),
		}
	})
	return gameServiceInstance
}

/**
内部接口：房间玩家准备/取消准备

req:
- playerOpenId: 玩家 openId

resp:
- roomInfoAfter: 房间信息（玩家加入后）
*/
func (s *GameService) roomInfoReadyPlayer(ctx context.Context, playerOpenId string, roomInfoBefore *model.RoomInfo) (roomInfoAfter *model.RoomInfo, err error) {
	logs.CtxInfo(ctx, "[Service GameService roomInfoReadyPlayer Req] req: %s, %s", playerOpenId, roomInfoBefore) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Service GameService roomInfoReadyPlayer Resp] resp: %s, %#v", roomInfoAfter, err) // 出口日志
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
		logs.CtxError(ctx, "[Service GameService roomInfoReadyPlayer] room has no player: %#v", playerOpenId)
		err = util.NewErrf("[Service GameService roomInfoReadyPlayer] room has no player: %#v", playerOpenId)
		return nil, err
	}
	// step2，判断玩家是否可准备，即房间为开放状态或等待状态可准备，房间为准备状态可取消准备
	if roomInfoBefore.State == model.RoomStateOpen { // 房间为开放状态，玩家可准备或取消准备
		curPlayerState := roomInfoBefore.PlayerList[index].State
		switch curPlayerState {
		case model.PlayerStateWait: // 玩家准备
			roomInfoBefore.PlayerList[index].State = model.PlayerStateReady
			break
		case model.PlayerStateReady: // 玩家取消准备
			roomInfoBefore.PlayerList[index].State = model.PlayerStateWait
			break
		default:
			logs.CtxFatal(ctx, "[Service GameService roomInfoReadyPlayer] unexpected error")
			panic("[Service GameService roomInfoReadyPlayer] unexpected error")
		}
	} else if roomInfoBefore.State == model.RoomStateWait { //房间为开放等待状态，玩家可准备或取消准备，若全部玩家准备则房间切换为准备状态
		curPlayerState := roomInfoBefore.PlayerList[index].State
		switch curPlayerState {
		case model.PlayerStateWait: // 玩家准备
			roomInfoBefore.PlayerList[index].State = model.PlayerStateReady
			var roomReady = true // 判断房间是否可准备
			for _, player := range roomInfoBefore.PlayerList {
				if player.State == model.PlayerStateWait { // 有一个玩家未准备，房间就不能准备
					roomReady = false
					break
				}
			}
			if roomReady {
				roomInfoBefore.State = model.RoomStateReady // 房间切换为准备状态
			}
			break
		case model.PlayerStateReady: // 玩家取消准备
			roomInfoBefore.PlayerList[index].State = model.PlayerStateWait
			break
		default:
			logs.CtxFatal(ctx, "[Service GameService roomInfoReadyPlayer] unexpected error")
			panic("[Service GameService roomInfoReadyPlayer] unexpected error")
		}
	} else if roomInfoBefore.State == model.RoomStateReady || roomInfoBefore.State == model.RoomStateClear { // 房间为准备状态或清算状态，玩家可取消准备，取消准备后房间切换为等待状态
		curPlayerState := roomInfoBefore.PlayerList[index].State
		switch curPlayerState {
		case model.PlayerStateWait: // 玩家准备
			logs.CtxFatal(ctx, "[Service GameService roomInfoReadyPlayer] unexpected error")
			panic("[Service GameService roomInfoReadyPlayer] unexpected error")
		case model.PlayerStateReady: // 玩家取消准备
			roomInfoBefore.PlayerList[index].State = model.PlayerStateWait
			roomInfoBefore.State = model.RoomStateWait
			break
		default:
			logs.CtxFatal(ctx, "[Service GameService roomInfoReadyPlayer] unexpected error")
			panic("[Service GameService roomInfoReadyPlayer] unexpected error")
		}
	} else { // 房间信息错误，玩家不可操作
		logs.CtxError(ctx, "[Service GameService roomInfoReadyPlayer] room can not ready player, room state: %#v", roomInfoBefore.State)
		err = util.NewErrf("[Service GameService roomInfoReadyPlayer] room can not ready player, room state: %#v", roomInfoBefore.State)
		return nil, err
	}
	// step3，赋值更改后的房间信息
	roomInfoAfter = roomInfoBefore
	/************************************ 核心逻辑 ****************************************/

	return roomInfoAfter, nil
}

/**
外部接口：准备游戏/取消准备游戏

req:
- playerOpenId: 玩家 openId

resp:
*/
func (s *GameService) ReadyGame(ctx context.Context, playerOpenId string, roomInfoBefore *model.RoomInfo) (roomInfoAfter *model.RoomInfo, err error) {
	logs.CtxInfo(ctx, "[Service GameService ReadyGame Req] req: %s, %s", playerOpenId, roomInfoBefore) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Service GameService ReadyGame Resp] resp: %s, %#v", roomInfoAfter, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	// step1，玩家准备
	roomInfoAfter, err = s.roomInfoReadyPlayer(ctx, playerOpenId, roomInfoBefore)
	if err != nil {
		logs.CtxError(ctx, "[Service GameService ReadyGame] roomInfoReadyPlayer error: %#v", err)
		err = util.NewErrf("[Service GameService ReadyGame] roomInfoReadyPlayer error: %#v", err)
		return nil, err
	}
	/************************************ 核心逻辑 ****************************************/

	return roomInfoAfter, nil
}

/**
内部接口：开始游戏

req:
- masterOpenId: 房主 openId

resp:
- roomInfoAfter: 房间信息（玩家加入后）
*/
func (s *GameService) roomInfoStartGame(ctx context.Context, masterOpenId string, roomInfoBefore *model.RoomInfo) (roomInfoAfter *model.RoomInfo, err error) {
	logs.CtxInfo(ctx, "[Service GameService roomInfoStartGame Req] req: %s, %s", masterOpenId, roomInfoBefore) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Service GameService roomInfoStartGame Resp] resp: %s, %#v", roomInfoAfter, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	// step1，玩家是否是房主
	if masterOpenId != roomInfoBefore.MasterOpenId { // 玩家不是房主
		logs.CtxError(ctx, "[Service GameService roomInfoStartGame] room not master: %#v", masterOpenId)
		err = util.NewErrf("[Service GameService roomInfoStartGame] room not master: %#v", masterOpenId)
		return nil, err
	}
	// step2，房间状态是否为准备中或清算中
	if roomInfoBefore.State != model.RoomStateReady && roomInfoBefore.State != model.RoomStateClear {
		logs.CtxError(ctx, "[Service GameService roomInfoStartGame] room state not ready or clear: %#v", roomInfoBefore.State)
		err = util.NewErrf("[Service GameService roomInfoStartGame] room state not ready or clear: %#v", roomInfoBefore.State)
		return nil, err
	}
	// step3，获取本轮游戏词汇
	if len(roomInfoBefore.WordCache) != 0 { // 如果房间存在词汇缓存则获取缓存的第一个词汇
		roomInfoBefore.Word = roomInfoBefore.WordCache[0]       // 获取词汇缓存
		roomInfoBefore.WordCache = roomInfoBefore.WordCache[1:] // 从词汇缓存中移除词汇
	} else { // 否则重新获取词汇缓存
		wordList, err := s.wordDao.RandomQueryWordList(ctx, model.RoomWordCacheLen)
		if err != nil {
			logs.CtxError(ctx, "[Service GameService roomInfoStartGame] RandomQueryWordList error: %#v", err)
			err = util.NewErrf("[Service GameService roomInfoStartGame] RandomQueryWordList error: %#v", err)
			return nil, err
		}
		roomInfoBefore.Word = wordList[0]       // 词汇缓存
		roomInfoBefore.WordCache = wordList[1:] // 更新词汇缓存
	}
	// step4，分配词汇给各个玩家
	indexRandomList := rand.Perm(roomInfoBefore.RoomSetting.TotalNum) // 随机排序后，顺序依次为卧底、白板、平民
	for i := 0; i < roomInfoBefore.RoomSetting.SpyNum; i++ {
		roomInfoBefore.PlayerList[indexRandomList[i]].Role = model.PlayerRoleSpy
		roomInfoBefore.PlayerList[indexRandomList[i]].Word = roomInfoBefore.Word.SpyWord
	}
	for i := roomInfoBefore.RoomSetting.SpyNum; i < roomInfoBefore.RoomSetting.SpyNum+roomInfoBefore.RoomSetting.BlankNum; i++ {
		roomInfoBefore.PlayerList[indexRandomList[i]].Role = model.PlayerRoleBlank
		roomInfoBefore.PlayerList[indexRandomList[i]].Word = roomInfoBefore.Word.BlankWord
	}
	for i := roomInfoBefore.RoomSetting.SpyNum + roomInfoBefore.RoomSetting.BlankNum; i < roomInfoBefore.RoomSetting.TotalNum; i++ {
		roomInfoBefore.PlayerList[indexRandomList[i]].Role = model.PlayerRoleNormal
		roomInfoBefore.PlayerList[indexRandomList[i]].Word = roomInfoBefore.Word.NormalWord
	}
	// step5，确定第一个玩家
	for _, index := range indexRandomList {
		player := roomInfoBefore.PlayerList[index]
		if player.Role != model.PlayerRoleBlank && !util.IsNil(player.Word) { // 第一个玩家一定不为白板，且词不为空）
			roomInfoBefore.BeginPlayer = player.NickName
			break
		}
	}
	// step6，更改房间状态为游戏中
	roomInfoBefore.State = model.RoomStatePlaying
	// step7，赋值更改后的房间信息
	roomInfoAfter = roomInfoBefore
	/************************************ 核心逻辑 ****************************************/

	return roomInfoAfter, nil
}

/**
外部接口：开始游戏

req:
- masterOpenId: 房主 openId

resp:
*/
func (s *GameService) StartGame(ctx context.Context, masterOpenId string, roomInfoBefore *model.RoomInfo) (roomInfoAfter *model.RoomInfo, err error) {
	logs.CtxInfo(ctx, "[Service GameService StartGame Req] req: %s, %s", masterOpenId, roomInfoBefore) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Service GameService StartGame Resp] resp: %s, %#v", roomInfoAfter, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	// step1，开始游戏
	roomInfoAfter, err = s.roomInfoStartGame(ctx, masterOpenId, roomInfoBefore)
	if err != nil {
		logs.CtxError(ctx, "[Service GameService StartGame] roomInfoStartGame error: %#v", err)
		err = util.NewErrf("[Service GameService StartGame] roomInfoStartGame error: %#v", err)
		return nil, err
	}
	/************************************ 核心逻辑 ****************************************/

	return roomInfoAfter, nil
}

/**
内部接口：结束游戏

req:
- masterOpenId: 房主 openId

resp:
- roomInfoAfter: 房间信息（玩家加入后）
*/
func (s *GameService) roomInfoEndGame(ctx context.Context, masterOpenId string, roomInfoBefore *model.RoomInfo) (roomInfoAfter *model.RoomInfo, err error) {
	logs.CtxInfo(ctx, "[Service GameService roomInfoEndGame Req] req: %s, %s", masterOpenId, roomInfoBefore) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Service GameService roomInfoEndGame Resp] resp: %s, %#v", roomInfoAfter, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	// step1，玩家是否是房主
	if masterOpenId != roomInfoBefore.MasterOpenId { // 玩家不是房主
		logs.CtxError(ctx, "[Service GameService roomInfoEndGame] room not master: %#v", masterOpenId)
		err = util.NewErrf("[Service GameService roomInfoEndGame] room not master: %#v", masterOpenId)
		return nil, err
	}
	// step2，房间状态是否为游戏中
	if roomInfoBefore.State != model.RoomStatePlaying {
		logs.CtxError(ctx, "[Service GameService roomInfoEndGame] room state not playing: %#v", roomInfoBefore.State)
		err = util.NewErrf("[Service GameService roomInfoEndGame] room state not playing: %#v", roomInfoBefore.State)
		return nil, err
	}
	//// step3，更改房间状态为开放中或等待中，或准备中
	//if len(roomInfoBefore.PlayerList) == roomInfoBefore.TotalNum { // 人已满
	//	roomInfoBefore.State = model.RoomStateWait
	//	var roomReady = true // 判断房间是否可准备
	//	for _, player := range roomInfoBefore.PlayerList {
	//		if player.State == model.PlayerStateWait { // 有一个玩家未准备，房间就不能准备
	//			roomReady = false
	//			break
	//		}
	//	}
	//	if roomReady {
	//		roomInfoBefore.State = model.RoomStateReady // 房间切换为准备状态
	//	}
	//} else {
	//	roomInfoBefore.State = model.RoomStateOpen // 人未满
	//}
	// step3，更改房间状态为清算中
	roomInfoBefore.State = model.RoomStateClear
	// step4，赋值更改后的房间信息
	roomInfoAfter = roomInfoBefore
	/************************************ 核心逻辑 ****************************************/

	return roomInfoAfter, nil
}

/**
外部接口：结束游戏

req:
- masterOpenId: 房主 openId

resp:
*/
func (s *GameService) EndGame(ctx context.Context, masterOpenId string, roomInfoBefore *model.RoomInfo) (roomInfoAfter *model.RoomInfo, err error) {
	logs.CtxInfo(ctx, "[Service GameService EndGame Req] req: %s, %s", masterOpenId, roomInfoBefore) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Service GameService EndGame Resp] resp: %s, %#v", roomInfoBefore, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	// step1，结束游戏
	roomInfoAfter, err = s.roomInfoEndGame(ctx, masterOpenId, roomInfoBefore)
	if err != nil {
		logs.CtxError(ctx, "[Service GameService EndGame] roomInfoEndGame error: %#v", err)
		err = util.NewErrf("[Service GameService EndGame] roomInfoEndGame error: %#v", err)
		return nil, err
	}
	/************************************ 核心逻辑 ****************************************/

	return roomInfoAfter, nil
}

/**
内部接口：开始游戏

req:
- masterOpenId: 房主 openId

resp:
- roomInfoAfter: 房间信息（玩家加入后）
*/
func (s *GameService) roomInfoRestartGame(ctx context.Context, masterOpenId string, roomInfoBefore *model.RoomInfo) (roomInfoAfter *model.RoomInfo, err error) {
	logs.CtxInfo(ctx, "[Service GameService roomInfoRestartGame Req] req: %s, %s", masterOpenId, roomInfoBefore) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Service GameService roomInfoRestartGame Resp] resp: %s, %#v", roomInfoAfter, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	// step1，玩家是否是房主
	if masterOpenId != roomInfoBefore.MasterOpenId { // 玩家不是房主
		logs.CtxError(ctx, "[Service GameService roomInfoRestartGame] room not master: %#v", masterOpenId)
		err = util.NewErrf("[Service GameService roomInfoRestartGame] room not master: %#v", masterOpenId)
		return nil, err
	}
	// step2，房间状态是否为游戏中（重开游戏机制基于游戏中玩家不可退出房间）
	if roomInfoBefore.State != model.RoomStatePlaying {
		logs.CtxError(ctx, "[Service GameService roomInfoRestartGame] room state not ready: %#v", roomInfoBefore.State)
		err = util.NewErrf("[Service GameService roomInfoRestartGame] room state not ready: %#v", roomInfoBefore.State)
		return nil, err
	}
	// step3，获取本轮游戏词汇
	if len(roomInfoBefore.WordCache) != 0 { // 如果房间存在词汇缓存则获取缓存的第一个词汇
		roomInfoBefore.Word = roomInfoBefore.WordCache[0]       // 获取词汇缓存
		roomInfoBefore.WordCache = roomInfoBefore.WordCache[1:] // 从词汇缓存中移除词汇
	} else { // 否则重新获取词汇缓存
		wordList, err := s.wordDao.RandomQueryWordList(ctx, model.RoomWordCacheLen)
		if err != nil {
			logs.CtxError(ctx, "[Service GameService roomInfoRestartGame] RandomQueryWordList error: %#v", err)
			err = util.NewErrf("[Service GameService roomInfoRestartGame] RandomQueryWordList error: %#v", err)
			return nil, err
		}
		roomInfoBefore.Word = wordList[0]       // 词汇缓存
		roomInfoBefore.WordCache = wordList[1:] // 更新词汇缓存
	}
	// step4，分配词汇给各个玩家
	indexRandomList := rand.Perm(roomInfoBefore.RoomSetting.TotalNum) // 随机排序后，顺序依次为卧底、白板、平民
	for i := 0; i < roomInfoBefore.RoomSetting.SpyNum; i++ {
		roomInfoBefore.PlayerList[indexRandomList[i]].Role = model.PlayerRoleSpy
		roomInfoBefore.PlayerList[indexRandomList[i]].Word = roomInfoBefore.Word.SpyWord
	}
	for i := roomInfoBefore.RoomSetting.SpyNum; i < roomInfoBefore.RoomSetting.SpyNum+roomInfoBefore.RoomSetting.BlankNum; i++ {
		roomInfoBefore.PlayerList[indexRandomList[i]].Role = model.PlayerRoleBlank
		roomInfoBefore.PlayerList[indexRandomList[i]].Word = roomInfoBefore.Word.BlankWord
	}
	for i := roomInfoBefore.RoomSetting.SpyNum + roomInfoBefore.RoomSetting.BlankNum; i < roomInfoBefore.RoomSetting.TotalNum; i++ {
		roomInfoBefore.PlayerList[indexRandomList[i]].Role = model.PlayerRoleNormal
		roomInfoBefore.PlayerList[indexRandomList[i]].Word = roomInfoBefore.Word.NormalWord
	}
	// step5，确定第一个玩家
	for _, index := range indexRandomList {
		player := roomInfoBefore.PlayerList[index]
		if player.Role != model.PlayerRoleBlank && !util.IsNil(player.Word) { // 第一个玩家一定不为白板，且词不为空）
			roomInfoBefore.BeginPlayer = player.NickName
			break
		}
	}
	// step6，更改房间状态为游戏中
	roomInfoBefore.State = model.RoomStatePlaying
	// step7，赋值更改后的房间信息
	roomInfoAfter = roomInfoBefore
	/************************************ 核心逻辑 ****************************************/

	return roomInfoAfter, nil
}

/**
外部接口：重新开始游戏

req:
- masterOpenId: 房主 openId

resp:
*/
func (s *GameService) RestartGame(ctx context.Context, masterOpenId string, roomInfoBefore *model.RoomInfo) (roomInfoAfter *model.RoomInfo, err error) {
	logs.CtxInfo(ctx, "[Service GameService RestartGame Req] req: %s, %s", masterOpenId, roomInfoBefore) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Service GameService RestartGame Resp] resp: %s, %#v", roomInfoBefore, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	// step1，重新开始游戏
	roomInfoAfter, err = s.roomInfoRestartGame(ctx, masterOpenId, roomInfoBefore)
	if err != nil {
		logs.CtxError(ctx, "[Service GameService RestartGame] roomInfoRestartGame error: %#v", err)
		err = util.NewErrf("[Service GameService RestartGame] roomInfoRestartGame error: %#v", err)
		return nil, err
	}
	/************************************ 核心逻辑 ****************************************/

	return roomInfoAfter, nil
}
