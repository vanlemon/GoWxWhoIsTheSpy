package config

import (
	"log"
	"testing"
)

func TestConfigInit(t *testing.T) {
	InitConfig("../conf/who_is_the_spy_dev.json")
	log.Printf("%#v\n", ConfigInstance)
	log.Printf("%#v\n", ConfigJson)
}
