package util

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

//结构体转化为 Json Map
//func StructToInterfaceMap(v interface{}) map[string]interface{} {
//	var ret map[string]interface{}
//	bytes, _ := json.Marshal(v)
//	// TODO 忽略 json.Unmarshal 报错
//	_ = json.Unmarshal(bytes, &ret)
//	return ret
//}

func TestJsonMapString(t *testing.T) {
	m:=make(map[string]string)
	m["key"]= "value"
	mJsonString,err:=json.Marshal(m)
	fmt.Println(err)
	fmt.Println(string(mJsonString))
}

type TestUser struct {
	ID         int64     `json:"id" gorm:"column:id"`                   // 自增ID
	OpenId     string    `json:"openid" gorm:"column:openid"`           // openid，和微信接口的字段名一致
	NickName   string    `json:"nick_name" gorm:"column:nick_name"`     // 用户昵称
	Gender     int       `json:"gender" gorm:"column:gender"`           // 用户性别
	AvatarUrl  string    `json:"avatar_url" gorm:"column:avatar_url"`   // 用户头像
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"` // 创建时间
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time"` // 更新时间
}

func TestJsonStringAndJsonMapString(t *testing.T) {
	reqData := map[string]interface{}{
		"tempId": "20210303140727681706115F",
		"userInfo": TestUser{
			NickName:   "test_nick_name",
			Gender:     0,
			AvatarUrl:  "",
			CreateTime: time.Time{},
			UpdateTime: time.Time{},
		},
	}

	fmt.Println(StructToJsonString(reqData))
	fmt.Println(InterfaceMapToJsonString(reqData))
}