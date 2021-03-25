package util

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetSign(t *testing.T) {
	data := make(map[string]interface{})
	data["key"] = "value"
	reqDataJsonString, err := json.Marshal(data)
	fmt.Println("reqDataJsonString:", string(reqDataJsonString))
	if err != nil {
		fmt.Println("err:", err)
	}
	sign := GetSign(string(reqDataJsonString))
	fmt.Println("sign:", sign)
}
