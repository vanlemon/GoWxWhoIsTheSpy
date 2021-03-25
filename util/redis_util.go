package util

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-redis/redis"
)

// 从配置 json 中创建 go-redis 数据库实例
func NewRedisFromJSON(js *simplejson.Json, configName string) (*redis.Client, error) {
	conf := js.Get(configName)

	// TODO：忽略配置解析错误
	host, _ := conf.Get("host").String()
	port, _ := conf.Get("port").Int()

	password, _ := conf.Get("password").String()
	db, _ := conf.Get("db").Int()

	return newRedisInstance(host, port, password, db)
}

// 从参数中创建 go-redis 数据库实例
func newRedisInstance(host string, port int, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     CombineIpAndPort(host, port),
		Password: password,
		DB:       db,
	})

	_, err := client.Ping().Result() // 检测是否连接成功

	if err != nil {
		return nil, err
	}

	return client, err
}
