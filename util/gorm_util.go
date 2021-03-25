package util

import (
	"fmt"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/jinzhu/gorm"

	"lmf.mortal.com/GoWxWhoIsTheSpy/config"

	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql 引擎
)

// 从配置 json 中创建 gorm 数据库实例
func NewGormFromJSON(js *simplejson.Json, configName string) (*gorm.DB, error) {
	conf := js.Get(configName)

	// TODO：忽略配置解析错误
	host, _ := conf.Get("host").String()
	port, _ := conf.Get("port").Int()

	databaseName, _ := conf.Get("database").String()
	settings, _ := conf.Get("settings").String()

	username, _ := conf.Get("username").String()
	password, _ := conf.Get("password").String()

	return newGormInstance(host, port, username, password, databaseName, settings)
}

// 从参数中创建 gorm 数据库实例
func newGormInstance(host string, port int, username, password, databaseName, settings string) (*gorm.DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", username, password, host, port, databaseName, settings) // 拼接数据库连接字符串

	db, err := gorm.Open("mysql", connStr) // 检测是否连接成功
	if err != nil {
		return nil, err
	}

	if config.IsDev() { // 开发模式打开数据库 debug 模式
		db.LogMode(true)
	}

	db.DB().SetConnMaxLifetime(300 * time.Second) // 最大连接空闲时间

	// 连接数 = ((核心数 * 2) + 有效磁盘数)
	db.DB().SetMaxIdleConns(3) // 最大空闲连接数
	db.DB().SetMaxOpenConns(3) // 最大开启连接数

	return db, nil
}
