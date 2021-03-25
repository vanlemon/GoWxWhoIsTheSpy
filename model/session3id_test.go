package model

import (
	"fmt"
	"testing"

	logs "lmf.mortal.com/GoLogs"
)

func TestSession3idDao_SetAndGet(t *testing.T) {
	tempId, err := Session3idDaoInstance().SetOpenIdGenID(logs.TestCtx(), "openid")
	fmt.Println(tempId, err)
	openid, err := Session3idDaoInstance().GetOpenId(logs.TestCtx(), tempId)
	fmt.Println(openid, err)
}

func BenchmarkSession3idDao_SetAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tempId, err := Session3idDaoInstance().SetOpenIdGenID(logs.TestCtx(), "openid")
		fmt.Println(tempId, err)
		openid, err := Session3idDaoInstance().GetOpenId(logs.TestCtx(), tempId)
		fmt.Println(openid, err)
	}
}
