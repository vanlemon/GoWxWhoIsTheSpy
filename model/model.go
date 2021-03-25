package model

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"

	"lmf.mortal.com/GoWxWhoIsTheSpy/util"

	"lmf.mortal.com/GoLogs"
)

var (
	LimiterRedis *redis.Client // 限流器-缓存

	GatewayRedis *redis.Client // 网关 redis 锁-缓存

	UserDB    *gorm.DB      // 用户信息-数据库
	UserRedis *redis.Client // 用户临时码-缓存

	WordDB *gorm.DB // 谁是卧底词库-数据库

	RoomRedis *redis.Client // 用户房间信息-缓存
)

// 初始化所有模型
func InitModel(configJson *simplejson.Json) {
	var err error

	err = initGatewayRedis(configJson)
	if err != nil {
		logs.CtxFatal(logs.SysCtx, "[Bootstrap Model] init Gateway DB error: %#v", err)
	}

	err = initUserDB(configJson)
	if err != nil {
		logs.CtxFatal(logs.SysCtx, "[Bootstrap Model] init User DB error: %#v", err)
	}

	err = initUserRedis(configJson)
	if err != nil {
		logs.CtxFatal(logs.SysCtx, "[Bootstrap Model] init User Redis error: %#v", err)
	}

	err = initWordDB(configJson)
	if err != nil {
		logs.CtxFatal(logs.SysCtx, "[Bootstrap Model] init Word DB error: %#v", err)
	}

	err = initRoomRedis(configJson)
	if err != nil {
		logs.CtxFatal(logs.SysCtx, "[Bootstrap Model] init Room Redis error: %#v", err)
	}
}

func initGatewayRedis(configJson *simplejson.Json) error {
	configName := "gateway_redis"
	var err error
	GatewayRedis, err = util.NewRedisFromJSON(configJson, configName)
	return err
}

func initUserDB(configJson *simplejson.Json) error {
	configName := "user_mysql"
	var err error
	UserDB, err = util.NewGormFromJSON(configJson, configName)
	return err
}

func initUserRedis(configJson *simplejson.Json) error {
	configName := "user_redis"
	var err error
	UserRedis, err = util.NewRedisFromJSON(configJson, configName)
	return err
}

func initWordDB(configJson *simplejson.Json) error {
	configName := "word_mysql"
	var err error
	WordDB, err = util.NewGormFromJSON(configJson, configName)
	return err
}

func initRoomRedis(configJson *simplejson.Json) error {
	configName := "room_redis"
	var err error
	RoomRedis, err = util.NewRedisFromJSON(configJson, configName)
	return err
}
