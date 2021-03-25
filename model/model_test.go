package model

import (
	logs "lmf.mortal.com/GoLogs"

	"lmf.mortal.com/GoWxWhoIsTheSpy/config"
)

func init() {
	config.InitConfig("../conf/who_is_the_spy_dev.json")
	logs.InitDefaultLogger(config.ConfigInstance.LogConfig)
	InitModel(config.ConfigJson)
}
