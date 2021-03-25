package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"lmf.mortal.com/GoWxWhoIsTheSpy/cconst"
	"sync"
	"time"

	"lmf.mortal.com/GoLogs"
	"lmf.mortal.com/GoWxWhoIsTheSpy/util"
)

/*************************************** Redis 常量信息 *********************************************/
const (
	Redis_RoomInfo_Lock_Prefix   = "roomInfo-lock-"   // Redis 锁的 ID 前缀
	Redis_RoomID_RoomInfo_Prefix = "roomId-roomInfo-" // Redis RoomId -> RoomInfo 键前缀
	Redis_OpenID_RoomID_Prefix   = "openId-roomId-"   // Redis OpenId -> RoomId 键前缀
)

const (
	RoomIDHexLen     = 4             // 房间 ID 末位十六进制长度，总长度为 20+4 位
	RoomIDRetryTimes = 3             // 若房间 ID 生成失败的重新尝试次数
	RoomIDTTLTime    = time.Hour * 8 // 房间 ID 的失效时间
)

const (
	RoomInfoLockRetryTimes = 3 // 房间信息 redis 锁重试次数
	// TODO: 测试确定重试时间
	RoomInfoLockWaitTime = 100 // ms, 房间信息 redis 锁重试间隔时间
)

/*************************************** Redis 常量信息 *********************************************/

/*************************************** 枚举常量信息 *********************************************/
type PlayerState string // 玩家状态
type PlayerRole string  // 玩家身份
type RoomState string   // 房间状态

const (
	PlayerStateWait  PlayerState = "Wait"  // 玩家状态 - 等待
	PlayerStateReady PlayerState = "Ready" // 玩家状态 - 准备

	PlayerRoleNormal PlayerRole = "Normal" // 玩家身份 - 正常
	PlayerRoleSpy    PlayerRole = "Spy"    // 玩家身份 - 卧底
	PlayerRoleBlank  PlayerRole = "Blank"  // 玩家身份 - 白板

	RoomStateOpen    RoomState = "Open"    // 房间状态 - 开放中（人未满）
	RoomStateWait    RoomState = "Wait"    // 房间状态 - 等待中（人已满）
	RoomStateReady   RoomState = "Ready"   // 房间状态 - 准备中（所有玩家已准备）
	RoomStatePlaying RoomState = "Playing" // 房间状态 - 游戏中（房主开始游戏）
	RoomStateClear   RoomState = "Clear"   // 房间状态 - 清算中（房主结束游戏）
)

/*************************************** 枚举常量信息 *********************************************/

/*************************************** 其他常量信息 *********************************************/
const RoomWordCacheLen = 10 // 房间词汇缓存大小

/*************************************** 其他常量信息 *********************************************/

/*************************************** 房间配置信息 *********************************************/
// 房间人数配置
type RoomSetting struct {
	TotalNum int `json:"total_num"` // 总人数
	SpyNum   int `json:"spy_num"`   // 卧底人数
	BlankNum int `json:"blank_num"` // 白板人数
}

// 玩家信息
type Player struct {
	OpenId    string      `json:"open_id"`    // openId
	NickName  string      `json:"nick_name"`  // 玩家用户名
	AvatarUrl string      `json:"avatar_url"` // 玩家头像
	State     PlayerState `json:"state"`      // 玩家状态
	Word      string      `json:"word"`       // 玩家本轮词汇
	Role      PlayerRole  `json:"role"`       // 玩家本轮角色
	IsSelf    bool        `json:"is_self"`    // 是否是玩家本人
}

// 房间信息
type RoomInfo struct {
	RoomId       string    `json:"room_id"`        // 房间 Id
	*RoomSetting `json:"room_setting"`             // 房间配置信息
	MasterOpenId string    `json:"master_open_id"` // 房主 openId
	IsMaster     bool      `json:"is_master"`      // 是否是房主本人
	PlayerList   []*Player `json:"player_list"`    // 玩家信息列表
	State        RoomState `json:"state"`          // 房间状态
	*Word        `json:"word"`                     // 房间本轮词汇
	BeginPlayer  string    `json:"begin_player"`   // 第一个玩家本轮用户名（第一个玩家一定不为白板，且词不为空）
	//TODO：WordSet      iset.Set  `json:"word_set"`       // 房间之前词汇集合（确保同一房间的词汇不重复），暂时不用，不进行重复词汇校检
	WordCache []*Word `json:"word_cache"` // 房间缓存词汇（每次获取词汇列表，之后从缓存中获取本轮词汇）
}

func (r *RoomInfo) String() string {
	bufPlayerList := bytes.Buffer{}
	for i, player := range r.PlayerList {
		bufPlayerList.WriteString(fmt.Sprintf("player(%d): %+v\n", i, player))
	}
	bufWordCache := bytes.Buffer{}
	for i, word := range r.WordCache {
		bufWordCache.WriteString(fmt.Sprintf("word(%d): %+v\n", i, word))
	}
	return fmt.Sprintf("RoomId: %s\nRoomSetting: %+v\nMasterOpenId: %s\nIsMaster: %+v\nPlayerList: \n%sState: %s\nWord: %+v\nBeginPlayer: %s\nWordCache:\n%s",
		r.RoomId, r.RoomSetting, r.MasterOpenId, r.IsMaster, bufPlayerList.String(), r.State, r.Word, r.BeginPlayer, bufWordCache.String())
}

/*************************************** 房间配置信息 *********************************************/

// 游戏房间数据访问层
type RoomDao struct { // key - value
	// 用户 openId - 房间 ID
	// 房间 ID - 房间信息
}

var roomDao *RoomDao      // 游戏房间数据访问层实例
var roomDaoOnce sync.Once // 单例模式

// 获取游戏房间数据访问层单例
func RoomDaoInstance() *RoomDao {
	roomDaoOnce.Do(func() {
		roomDao = &RoomDao{}
	})
	return roomDao
}

/**
生成 roomId，创建 roomId -> roomInfo 关联关系

Set: roomId -> roomInfo
*/
func (d *RoomDao) SetRoomInfoGenID(ctx context.Context, roomInfo *RoomInfo) (roomId string, err error) {
	logs.CtxInfo(ctx, "[Model Redis SetRoomInfoGenID Req] req: %s", roomInfo) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Model Redis SetRoomInfoGenID Resp] resp: %s, %#v", roomId, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	for i := 0; i < RoomIDRetryTimes; i++ { // 重试多次
		// 生成 roomId
		roomId = util.GenUUID(RoomIDHexLen)
		logs.CtxInfo(ctx, "[Model Redis SetRoomInfoGenID] GenUUID: %s", roomId)
		roomInfo.RoomId = roomId // 写入 roomId
		// 写入 Redis
		var success bool
		success, err = RoomRedis.SetNX(cconst.RedisKeyPrefix(Redis_RoomID_RoomInfo_Prefix, roomId), util.StructToJsonString(roomInfo), RoomIDTTLTime).Result()
		// 访问 Redis 失败
		if err != nil {
			logs.CtxError(ctx, "[Model Redis SetRoomInfoGenID] SetNX error: %#v", err)
			continue
		}
		// tempId 重复
		if !success {
			logs.CtxError(ctx, "[Model Redis SetRoomInfoGenID] SetNX Duplicate")
			err = util.NewErrf("[Model Redis SetRoomInfoGenID] SetNX Duplicate") // 更新 error 信息
			continue
		}
		// 生成 tempId 成功
		logs.CtxInfo(ctx, "[Model Redis SetRoomInfoGenID] SetNX Success: %#v", roomId)
		// 成功
		return roomId, nil
	}
	/************************************ 核心逻辑 ****************************************/

	// 失败
	return "", err
}

/**
由 roomId 获取 roomInfo

Get: roomId -> roomInfo
*/
func (d *RoomDao) GetRoomInfo(ctx context.Context, roomId string) (roomInfo *RoomInfo, err error) {
	logs.CtxInfo(ctx, "[Model Redis GetRoomInfo Req] req: %s", roomId) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Model Redis GetRoomInfo Resp] resp: %s, %#v", roomInfo, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	var roomInfoValue string // roomInfo 值
	roomInfoValue, err = RoomRedis.Get(cconst.RedisKeyPrefix(Redis_RoomID_RoomInfo_Prefix, roomId)).Result()
	if err == redis.Nil {
		logs.CtxError(ctx, "[Model Redis GetRoomInfo] Get error: %#v", err) // TODO：降级为 Warn
		return nil, cconst.RedisNil
	}
	if err != nil {
		logs.CtxError(ctx, "[Model Redis GetRoomInfo] Get error: %#v", err)
		return nil, err
	}
	var roomInfoData RoomInfo
	err = json.Unmarshal([]byte(roomInfoValue), &roomInfoData)
	if err != nil {
		logs.CtxError(ctx, "[Model Redis GetRoomInfo] roomInfoValue: %#v error: %#v", roomInfoValue, err)
	}
	roomInfo = &roomInfoData // 赋值给指针值
	/************************************ 核心逻辑 ****************************************/

	return roomInfo, nil
}

/**
创建 roomId -> roomInfo 关联关系（需要在 GetRoomInfo 之前就加锁，在 SetRoomInfo 之后解锁）

Set: roomId -> roomInfo
*/
func (d *RoomDao) SetRoomInfo(ctx context.Context, roomId string, roomInfo *RoomInfo) (err error) {
	logs.CtxInfo(ctx, "[Model Redis SetRoomInfo Req] req: %s, %s", roomId, roomInfo) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Model Redis SetRoomInfo Resp] resp: %#v", err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	err = RoomRedis.Set(cconst.RedisKeyPrefix(Redis_RoomID_RoomInfo_Prefix, roomId), util.StructToJsonString(roomInfo), RoomIDTTLTime).Err() // 此处不使用 SetNX，会覆盖掉房间信息，所以需要加锁后使用
	if err != nil {
		logs.CtxError(ctx, "[Model Redis SetRoomInfo] Set error: %#v", err)
		return err
	}
	/************************************ 核心逻辑 ****************************************/

	return nil
}

/**
删除 roomId -> roomInfo 关联关系

Del: roomId -> roomInfo
*/
func (d *RoomDao) DelRoomInfo(ctx context.Context, roomId string) (err error) {
	logs.CtxInfo(ctx, "[Model Redis DelRoomInfo Req] req: %#v", roomId) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Model Redis DelRoomInfo Resp] resp: %#v", err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	err = RoomRedis.Del(cconst.RedisKeyPrefix(Redis_RoomID_RoomInfo_Prefix, roomId)).Err()
	if err != nil {
		logs.CtxError(ctx, "[Model Redis DelRoomInfo] Set error: %#v", err)
		return err
	}
	/************************************ 核心逻辑 ****************************************/

	return nil
}

/**
roomInfo 加锁，无法根据 roomId 查看或修改 roomInfo
*/
func (d *RoomDao) LockRoomInfo(ctx context.Context, roomId string) (err error) {
	logs.CtxInfo(ctx, "[Model Redis LockRoomInfo Req] req: %#v", roomId) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Model Redis LockRoomInfo Resp] resp: %#v", err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	success, err := RoomRedis.SetNX(cconst.RedisKeyPrefix(Redis_RoomInfo_Lock_Prefix, roomId), nil, 0).Result()
	if err != nil {
		logs.CtxError(ctx, "[Model Redis LockRoomInfo] SetNX error: %#v", err)
		return err
	}
	if !success {
		logs.CtxInfo(ctx, "[Model Redis LockRoomInfo] SetNX duplicate")
		return cconst.RedisLockDuplicate
	}
	/************************************ 核心逻辑 ****************************************/

	return nil
}

/**
roomInfo 解锁
*/
func (d *RoomDao) UnLockRoomInfo(ctx context.Context, roomId string) (err error) {
	logs.CtxInfo(ctx, "[Model Redis UnLockRoomInfo Req] req: %#v", roomId) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Model Redis UnLockRoomInfo Resp] resp: %#v", err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	err = RoomRedis.Del(cconst.RedisKeyPrefix(Redis_RoomInfo_Lock_Prefix, roomId)).Err()
	if err != nil {
		logs.CtxError(ctx, "[Model Redis UnLockRoomInfo] Del error: %#v", err)
		return err
	}
	/************************************ 核心逻辑 ****************************************/

	return nil
}

/**
创建 openId -> roomId 关联关系

Set: openId -> roomId
*/
func (d *RoomDao) SetRoomID(ctx context.Context, openId, roomId string) (err error) {
	logs.CtxInfo(ctx, "[Model Redis SetRoomID Req] req: %#v, %#v", openId, roomId) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Model Redis SetRoomID Resp] resp: %#v", err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	err = RoomRedis.Set(cconst.RedisKeyPrefix(Redis_OpenID_RoomID_Prefix, openId), roomId, RoomIDTTLTime).Err() // 此处不使用 SetNX，会覆盖掉用户已有的房间关联
	if err != nil {
		logs.CtxError(ctx, "[Model Redis SetRoomID] Set error: %#v", err)
		return err
	}
	/************************************ 核心逻辑 ****************************************/

	return nil
}

/**
删除 openId -> roomId 关联关系

Set: openId -> roomId
*/
func (d *RoomDao) DelRoomID(ctx context.Context, openId string) (err error) {
	logs.CtxInfo(ctx, "[Model Redis DelRoomID Req] req: %#v", openId) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Model Redis DelRoomID Resp] resp: %#v", err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	err = RoomRedis.Del(cconst.RedisKeyPrefix(Redis_OpenID_RoomID_Prefix, openId)).Err()
	if err != nil {
		logs.CtxError(ctx, "[Model Redis DelRoomID] Set error: %#v", err)
		return err
	}
	/************************************ 核心逻辑 ****************************************/

	return nil
}

/**
由 openId 获取 roomId

Set: openId -> roomId
*/
func (d *RoomDao) GetRoomID(ctx context.Context, openId string) (roomId string, err error) {
	logs.CtxInfo(ctx, "[Model Redis GetRoomID Req] req: %s", openId) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Model Redis GetRoomID Resp] resp: %s, %#v", roomId, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	roomId, err = RoomRedis.Get(cconst.RedisKeyPrefix(Redis_OpenID_RoomID_Prefix, openId)).Result()
	if err == redis.Nil {
		logs.CtxError(ctx, "[Model Redis GetRoomID] Get error: %#v", err) // TODO：降级为 Warn
		return "", cconst.RedisNil
	}
	if err != nil {
		logs.CtxError(ctx, "[Model Redis GetRoomID] Get error: %#v", err)
		return "", err
	}
	/************************************ 核心逻辑 ****************************************/

	return roomId, nil
}
