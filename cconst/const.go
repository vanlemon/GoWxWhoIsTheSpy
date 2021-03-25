package cconst

import "lmf.mortal.com/GoWxWhoIsTheSpy/util"

/************************************ 键值常量 ****************************************/
const (
	CtxKey     = "ctx_key"      // gin 请求上下文中 ctx 的键
	ReqDataKey = "req_data_key" // gin 请求上下文中请求参数
)


// 获取带前缀的 Redis 键
func RedisKeyPrefix(prefix, redisKey string) string {
	return prefix + redisKey
}

/************************************ 键值常量 ****************************************/

/************************************ 错误常量 ****************************************/
var RedisLockDuplicate = util.NewErrf("RedisLockDuplicate") // Redis 锁重入错误
var RedisNil = util.NewErrf("RedisNil")                     // Redis 键不存在错误
/************************************ 错误常量 ****************************************/
