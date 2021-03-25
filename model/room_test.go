package model

import (
	"fmt"
	logs "lmf.mortal.com/GoLogs"
	"testing"
)

func TestRoomDao_SetRoomInfoGenID(t *testing.T) {
	roomInfo := &RoomInfo{
		RoomId:       "",
		RoomSetting:  &RoomSetting{},
		MasterOpenId: "",
		PlayerList: []*Player{
			&Player{},
			&Player{},
		},
		State:       "",
		Word:        &Word{},
		BeginPlayer: "",
		WordCache: []*Word{
			&Word{},
			&Word{},
		},
	}
	roomId, err := RoomDaoInstance().SetRoomInfoGenID(logs.TestCtx(), roomInfo)
	fmt.Println(roomId, err)
}

func TestRoomDao_GetRoomInfo(t *testing.T) {
	roomInfo, err := RoomDaoInstance().GetRoomInfo(logs.TestCtx(), "202103031545124745831323")
	fmt.Println(roomInfo, err)
}

func TestRoomInfo_String(t *testing.T) {
	roomInfo := &RoomInfo{
		RoomId:       "",
		RoomSetting:  &RoomSetting{},
		MasterOpenId: "",
		PlayerList: []*Player{
			&Player{},
			&Player{},
		},
		State:       "",
		Word:        &Word{},
		BeginPlayer: "",
		WordCache: []*Word{
			&Word{},
			&Word{},
		},
	}
	fmt.Println(roomInfo)
}
