package main

import (
	"github.com/gin-gonic/gin"
	limiter "lmf.mortal.com/GoLimiter"
	logs "lmf.mortal.com/GoLogs"
	"lmf.mortal.com/GoWxWhoIsTheSpy/config"
	"lmf.mortal.com/GoWxWhoIsTheSpy/handler"
	"lmf.mortal.com/GoWxWhoIsTheSpy/middleware"
	"lmf.mortal.com/GoWxWhoIsTheSpy/model"
	"lmf.mortal.com/GoWxWhoIsTheSpy/service"
	"lmf.mortal.com/GoWxWhoIsTheSpy/util"
	"log"
	"os"
)

func init() {
	// step1，初始化配置文件，若失败则直接启动失败
	if len(os.Args) > 0 {
		switch os.Args[1] {
		case config.Dev:
			config.InitConfig(util.GetExecPath() + "/../conf/who_is_the_spy_dev.json")
			break
		case config.Prod:
			config.InitConfig(util.GetExecPath() + "/../conf/who_is_the_spy_prod.json")
			break
		case config.Test:
			config.InitConfig(util.GetExecPath() + "/../conf/who_is_the_spy_test.json")
			break
		default: // 默认配置
			log.Printf("[Bootstrap Server] Running in Default Dev Env\n")
			config.InitConfig(util.GetExecPath() + "/../conf/who_is_the_spy_dev.json")
			break
		}
		log.Printf("[Bootstrap Server] Running in %s Env\n", os.Args[1])
	} else {
		log.Printf("[Bootstrap Server] Running in Default Dev Env\n")
		config.InitConfig(util.GetExecPath() + "/../conf/who_is_the_spy_dev.json")
	}

	//config.InitConfig("./conf/who_is_the_spy_dev.json")

	// step2，初始化日志服务
	logs.InitDefaultLogger(config.ConfigInstance.LogConfig)

	// step3，初始化数据模型连接，若失败则直接启动失败
	model.InitModel(config.ConfigJson)

	// step4，初始化限流器中间件
	limiter.InitOverLoadMiddleWare(
		config.ConfigInstance.LimiterConfig,
		config.ConfigJson.Get("limiter_config"),
		model.LimiterRedis)

	// step5，初始化核心服务（单例）
	service.InitService()

	// step6，打印所有配置
	logs.CtxInfo(logs.SysCtx, "[Bootstrap Server] config: %#v", config.ConfigInstance)
	logs.CtxInfo(logs.SysCtx, "[Bootstrap Server] configJson: %#v", config.ConfigJson)
}

func main() {
	// step1，初始化 gin 服务器，设置运行模式
	if config.IsDev() {
		gin.SetMode(gin.DebugMode)
	} else if config.IsTest() {
		gin.SetMode(gin.TestMode)
	} else if config.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		logs.CtxFatal(logs.SysCtx, "[Bootstrap Server] unexpected config env: %#v", config.ConfigInstance.Env)
	}
	server := gin.Default() // 设置运行模式后再初始化

	// step2，初始化各个中间件
	middleware.InitMiddleware(server)

	// step3，初始化各个 HTTP 服务
	handler.InitHandler(server)

	// step4，启动 gin 服务器
	err := server.Run(":9205")
	if err != nil {
		logs.CtxFatal(logs.SysCtx, "[Bootstrap Server] GIN Server run error: %#v", err)
	}
}
