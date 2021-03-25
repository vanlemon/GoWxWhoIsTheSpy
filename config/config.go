package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	limiter "lmf.mortal.com/GoLimiter"
	logs "lmf.mortal.com/GoLogs"

	"github.com/bitly/go-simplejson"
)

const (
	Dev  = "dev"
	Prod = "prod"
	Test = "test"
)

type UserConfig struct {
	WxLoginUrl string `json:"wx_login_url"`
	AppId      string `json:"app_id"`
	AppSecret  string `json:"app_secret"`
}

// 创建一个配置结构体，包含所有的配置对象
type Config struct {
	Env                   string   `json:"env"`          // 服务运行环境
	GatewayList           []string `json:"gateway_list"` // 网关服务列表
	Salt                  string   `json:"salt"`         // 请求秘钥
	logs.LogConfig        `json:"log_config"`            // 日志服务配置
	limiter.LimiterConfig `json:"limiter_config"`        // 限流服务配置
	UserConfig            `json:"user_config"`           // 用户服务配置
}

var (
	// ConfigInstance 当前环境配置信息
	ConfigInstance *Config
	// origin File 将配置信息解成json
	ConfigJson *simplejson.Json
)

// InitConfig 初始化配置文件
// 日志服务此时尚未启动，所有的启动错误通过 panic 汇报
// 配置初始化失败，则系统无法启动，直接 panic
func InitConfig(file string) {
	confContent, err := ioutil.ReadFile(file) // 读取文件信息
	if err != nil {
		panic(fmt.Sprintf("[Init Config] create new config error: %#v\n", err))
	}

	var conf Config
	err = json.Unmarshal(confContent, &conf) // 解析到 Config 结构体
	if err != nil {
		panic(fmt.Sprintf("[Init Config] json unmarshal error: %#v\n", err))
	}
	ConfigInstance = &conf // 赋值到 ConfigInstance

	confJson, err := simplejson.NewJson(confContent) // 解析到 Json
	if err != nil {
		panic(fmt.Sprintf("[Init Config] json unmarshal error: %#v\n", err))
	}
	ConfigJson = confJson // 赋值到 ConfigJson
}

func IsDev() bool {
	return ConfigInstance.Env == Dev
}

func IsProd() bool {
	return ConfigInstance.Env == Prod
}

func IsTest() bool {
	return ConfigInstance.Env == Test
}

func GatewayExists(methodKey string) bool {
	for _, eachItem := range ConfigInstance.GatewayList {
		if eachItem == methodKey {
			return true
		}
	}
	return false
}

func GetSalt() string {
	return ConfigInstance.Salt
}
