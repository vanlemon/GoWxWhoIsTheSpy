package model

import (
	"context"
	"github.com/go-redis/redis"
	"lmf.mortal.com/GoWxWhoIsTheSpy/cconst"
	"sync"
	"time"

	"lmf.mortal.com/GoLogs"
	"lmf.mortal.com/GoWxWhoIsTheSpy/util"
)

/*************************************** Redis 常量信息 *********************************************/
const (
	Redis_TempID_OpenID_Prefix = "tempId-openId-" // Redis TempId -> OpenId 键前缀
)

const (
	S3idHexLen     = 4             // 临时 ID 末位十六进制长度，总长度为 20+4 位
	S3idRetryTimes = 3             // 若临时 ID 生成失败的重新尝试次数
	S3idTTLTime    = time.Hour * 8 // 临时 ID 的失效时间
)

/*************************************** Redis 常量信息 *********************************************/

// 微信接口返回的 openid 结构体
type OpenIdAndSessionKey struct {
	OpenId     string `json:"openid"` //与微信请求返回值一致
	SessionKey string `json:"session_key"`
}

// 用户临时 ID 数据访问层
type Session3idDao struct { // key - value
	// 用户临时 ID - 用户 openId
}

var session3idDao *Session3idDao // 用户临时 ID 数据访问层实例
var session3idDaoOnce sync.Once  // 单例模式

// 获取用户临时 ID 数据访问层单例
func Session3idDaoInstance() *Session3idDao {
	session3idDaoOnce.Do(func() {
		session3idDao = &Session3idDao{}
	})
	return session3idDao
}

/**
生成 tempId，创建 tempId -> openId 关联关系

Set: tempId -> openId
*/
func (d *Session3idDao) SetOpenIdGenID(ctx context.Context, openId string) (tempId string, err error) {
	logs.CtxInfo(ctx, "[Model Redis SetOpenIdGenID Req] req: %s", openId) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Model Redis SetOpenIdGenID Resp] resp: %s, %#v", tempId, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	for i := 0; i < S3idRetryTimes; i++ { // 重试多次
		// 生成 tempId
		tempId = util.GenUUID(S3idHexLen)
		logs.CtxInfo(ctx, "[Model Redis SetOpenIdGenID] GenUUID: %s", tempId)
		// 写入 Redis
		var success bool
		success, err = UserRedis.SetNX(cconst.RedisKeyPrefix(Redis_TempID_OpenID_Prefix, tempId), openId, S3idTTLTime).Result()
		// 访问 Redis 失败
		if err != nil {
			logs.CtxError(ctx, "[Model Redis SetOpenIdGenID] SetNX error: %#v", err)
			continue
		}
		// tempId 重复
		if !success {
			logs.CtxError(ctx, "[Model Redis SetOpenIdGenID] SetNX Duplicate")
			err = util.NewErrf("[Model Redis SetOpenIdGenID] SetNX Duplicate") // 更新 error 信息
			continue
		}
		// 生成 tempId 成功
		logs.CtxInfo(ctx, "[Model Redis SetOpenIdGenID] SetNX Success: %#v", tempId)
		// 成功
		return tempId, nil
	}
	/************************************ 核心逻辑 ****************************************/

	// 失败
	return "", err
}

/**
由 tempId 获取 openId

Get: tempId -> openId
*/
func (d *Session3idDao) GetOpenId(ctx context.Context, tempId string) (openId string, err error) {
	logs.CtxInfo(ctx, "[Model Redis GetOpenId Req] req: %s", tempId) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Model Redis GetOpenId Resp] resp: %s, %#v", openId, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	openId, err = UserRedis.Get(cconst.RedisKeyPrefix(Redis_TempID_OpenID_Prefix, tempId)).Result()
	if err == redis.Nil {
		logs.CtxError(ctx, "[Model Redis GetOpenId] Get error: %#v", err) // TODO：降级为 Wran
		return "", cconst.RedisNil
	}
	if err != nil {
		logs.CtxError(ctx, "[Model Redis GetOpenId] Get error: %#v", err)
		return "", err
	}
	/************************************ 核心逻辑 ****************************************/

	return openId, nil
}
