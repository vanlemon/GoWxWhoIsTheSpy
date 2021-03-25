package client
//
//import (
//	"fmt"
//	logs "lmf.mortal.com/GoLogs"
//	"lmf.mortal.com/GoWxWhoIsTheSpy/config"
//	"lmf.mortal.com/GoWxWhoIsTheSpy/model"
//	"lmf.mortal.com/GoWxWhoIsTheSpy/util"
//	"math/rand"
//	"sync"
//	"testing"
//	"time"
//)
//
//var TestPlayerList []string // 测试玩家的临时 ID 集合
//
//var server = "http://something:9205" // 测试环境
////var server = "http://127.0.0.1:9205" // 开发环境
//
//func init() {
//	config.InitConfig("../conf/who_is_the_spy_dev.json")
//	logs.InitDefaultLogger(config.ConfigInstance.LogConfig)
//}
//
//// 初始化随机玩家
//func InitTestPlayerList(num int) {
//	for i := 0; i < num; i++ {
//		// 用户登录
//		postUrl := server + "/user/wx_login"
//		reqData := map[string]interface{}{
//			"code": "test_code",
//		}
//
//		reqDataJsonString := util.InterfaceMapToJsonString(reqData)
//		reqSign := util.GetSign(reqDataJsonString)
//
//		respData, _ := util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//		tempId := util.JsonStringToStringMap(respData)["Data"]
//
//		// 用户信息
//		postUrl = server + "/user/update_userInfo"
//		reqData = map[string]interface{}{
//			"tempId": tempId,
//			"userInfo": model.User{
//				NickName:   "test_nick_name",
//				Gender:     0,
//				AvatarUrl:  "test_url",
//				City:       "test_city",
//				Province:   "test_province",
//				Country:    "test_country",
//				Language:   "test_language",
//				CreateTime: time.Time{},
//				UpdateTime: time.Time{},
//			},
//		}
//
//		reqDataJsonString = util.InterfaceMapToJsonString(reqData)
//		reqSign = util.GetSign(reqDataJsonString)
//
//		respData, _ = util.HttpPost(postUrl, reqDataJsonString, reqSign)
//		fmt.Println("init:", respData)
//
//		// 添加用户到测试列表
//		TestPlayerList = append(TestPlayerList, tempId)
//	}
//
//	fmt.Println("TestPlayerList:", TestPlayerList)
//}
//
//func TestUserWxLogin(t *testing.T) {
//	postUrl := server + "/user/wx_login"
//	reqData := map[string]interface{}{
//		"code": "test_code",
//	}
//
//	reqDataJsonString := util.InterfaceMapToJsonString(reqData)
//	reqSign := util.GetSign(reqDataJsonString)
//
//	respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println("reqDataJsonString:", reqDataJsonString)
//		fmt.Println("respData:", respData)
//	}
//}
//
//func TestUserUpdateUserInfo(t *testing.T) {
//	postUrl := server + "/user/update_userInfo"
//	reqData := map[string]interface{}{
//		"tempId": "20210303140727681706115F",
//		"userInfo": model.User{
//			NickName:   "test_nick_name",
//			Gender:     0,
//			AvatarUrl:  "",
//			CreateTime: time.Time{},
//			UpdateTime: time.Time{},
//		},
//	}
//
//	reqDataJsonString := util.InterfaceMapToJsonString(reqData)
//	reqSign := util.GetSign(reqDataJsonString)
//
//	respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println("reqDataJsonString:", reqDataJsonString)
//		fmt.Println("respData:", respData)
//	}
//}
//
//func TestRoomNewRoom(t *testing.T) {
//	postUrl := server + "/room/new_room"
//	reqData := map[string]interface{}{
//		"tempId": "20210303140727681706115F",
//		"roomSetting": model.RoomSetting{
//			TotalNum: 6,
//			SpyNum:   2,
//			BlankNum: 0,
//		},
//	}
//
//	reqDataJsonString := util.InterfaceMapToJsonString(reqData)
//	reqSign := util.GetSign(reqDataJsonString)
//
//	logs.CtxInfo(logs.TestCtx(), "data: %s, sign: %s", reqDataJsonString, reqSign)
//
//	respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//	// TODO:RedisNil 应为 UserNotLogin
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println("reqDataJsonString:", reqDataJsonString)
//		fmt.Println("respData:", respData)
//	}
//}
//
//func TestRoomNewRoomAndEnterRoom(t *testing.T) {
//	// 初始化测试玩家
//	InitTestPlayerList(10)
//	// 创建房间
//	postUrl := server + "/room/new_room"
//	reqData := map[string]interface{}{
//		"tempId": TestPlayerList[0],
//		"roomSetting": model.RoomSetting{
//			TotalNum: 6,
//			SpyNum:   2,
//			BlankNum: 0,
//		},
//	}
//
//	reqDataJsonString := util.InterfaceMapToJsonString(reqData)
//	reqSign := util.GetSign(reqDataJsonString)
//
//	respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println("reqDataJsonString:", reqDataJsonString)
//		fmt.Println("respData:", respData)
//	}
//
//	roomId := util.JsonStringToInterfaceMap(respData)["Data"].(map[string]interface{})["roomId"].(string)
//	fmt.Println("roomId:", roomId)
//
//	wg := sync.WaitGroup{}
//	wg.Add(10)
//
//	// 玩家加入房间
//	postUrl = server + "/room/enter_room"
//	reqData = map[string]interface{}{
//		"tempId": TestPlayerList[1],
//		"roomId": roomId,
//	}
//
//	reqDataJsonString = util.InterfaceMapToJsonString(reqData)
//	reqSign = util.GetSign(reqDataJsonString)
//
//	respData, err = util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println("reqDataJsonString:", reqDataJsonString)
//		fmt.Println("respData:", respData)
//	}
//}
//
//// 单人重复加入房间
///**
//enter: {"Data":null,"LogId":"20210304154728805373F8C3","Message":"Success","Success":true}
//enter: {"Data":null,"LogId":"202103041547296658914E5D","Message":"[Service RoomService EnterRoom] roomInfoAddPlayer error: \u0026errors.errorString{s:\"[Service RoomService roomInfoAddPlayer] createPlayerFromUserOpenId error: user already exists\"}","Success":false}
//enter: {"Data":null,"LogId":"20210304154730193971965C","Message":"[Service RoomService EnterRoom] roomInfoAddPlayer error: \u0026errors.errorString{s:\"[Service RoomService roomInfoAddPlayer] createPlayerFromUserOpenId error: user already exists\"}","Success":false}
//enter: {"Data":null,"LogId":"202103041547304305765E92","Message":"[Service RoomService EnterRoom] roomInfoAddPlayer error: \u0026errors.errorString{s:\"[Service RoomService roomInfoAddPlayer] createPlayerFromUserOpenId error: user already exists\"}","Success":false}
//enter: {"Data":null,"LogId":"2021030415473088955151F2","Message":"[Service RoomService EnterRoom] roomInfoAddPlayer error: \u0026errors.errorString{s:\"[Service RoomService roomInfoAddPlayer] createPlayerFromUserOpenId error: user already exists\"}","Success":false}
//enter: {"Data":null,"LogId":"202103041547324084741B2A","Message":"[Service RoomService EnterRoom] roomInfoAddPlayer error: \u0026errors.errorString{s:\"[Service RoomService roomInfoAddPlayer] createPlayerFromUserOpenId error: user already exists\"}","Success":false}
//enter: {"Data":null,"LogId":"202103041547327727575FED","Message":"[Service RoomService EnterRoom] roomInfoAddPlayer error: \u0026errors.errorString{s:\"[Service RoomService roomInfoAddPlayer] createPlayerFromUserOpenId error: user already exists\"}","Success":false}
//enter: {"Data":null,"LogId":"202103041547362368530CF8","Message":"[Service RoomService EnterRoom] roomInfoAddPlayer error: \u0026errors.errorString{s:\"[Service RoomService roomInfoAddPlayer] createPlayerFromUserOpenId error: user already exists\"}","Success":false}
//enter: {"Data":null,"LogId":"20210304154736428522E42A","Message":"[Service RoomService EnterRoom] roomInfoAddPlayer error: \u0026errors.errorString{s:\"[Service RoomService roomInfoAddPlayer] createPlayerFromUserOpenId error: user already exists\"}","Success":false}
//*/
//func TestRoomNewRoomAndEnterRoomMany1(t *testing.T) {
//	// 初始化测试玩家
//	num := 10
//	InitTestPlayerList(num)
//	// 创建房间
//	postUrl := server + "/room/new_room"
//	reqData := map[string]interface{}{
//		"tempId": TestPlayerList[0],
//		"roomSetting": model.RoomSetting{
//			TotalNum: 6,
//			SpyNum:   2,
//			BlankNum: 0,
//		},
//	}
//
//	reqDataJsonString := util.InterfaceMapToJsonString(reqData)
//	reqSign := util.GetSign(reqDataJsonString)
//
//	respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println(respData)
//	}
//
//	roomId := util.JsonStringToInterfaceMap(respData)["Data"].(map[string]interface{})["roomId"].(string)
//	fmt.Println("roomId:", roomId)
//
//	wg := sync.WaitGroup{}
//	wg.Add(num - 1)
//
//	// 玩家加入房间
//	for i := 1; i < num; i++ {
//		i := i
//		go func() {
//			postUrl := server + "/room/enter_room"
//			reqData := map[string]interface{}{
//				"tempId": TestPlayerList[1],
//				"roomId": roomId,
//			}
//
//			reqDataJsonString := util.InterfaceMapToJsonString(reqData)
//			reqSign := util.GetSign(reqDataJsonString)
//
//			fmt.Println("enter(tempId):", i, TestPlayerList[i], reqSign)
//
//			time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)
//			respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//			if err != nil {
//				fmt.Println(err)
//			} else {
//				fmt.Println("enter:", respData)
//			}
//
//			wg.Done()
//		}()
//	}
//
//	wg.Wait()
//}
//
//// 多人加入房间
///**
//enter: {"Data":null,"LogId":"202103041548253905121949","Message":"Success","Success":true}
//enter: {"Data":null,"LogId":"202103041548262526745F56","Message":"Success","Success":true}
//enter: {"Data":null,"LogId":"202103041548267807279A1A","Message":"Success","Success":true}
//enter: {"Data":null,"LogId":"20210304154827149528D7D","Message":"Success","Success":true}
//enter: {"Data":null,"LogId":"202103041548274740769D65","Message":"Success","Success":true}
//enter: {"Data":null,"LogId":"202103041548289901599EB1","Message":"[Service RoomService EnterRoom] roomInfoAddPlayer error: \u0026errors.errorString{s:\"[Service RoomService roomInfoAddPlayer] room can not add player, room state: \\\"Wait\\\"\"}","Success":false}
//enter: {"Data":null,"LogId":"202103041548293578560816","Message":"[Service RoomService EnterRoom] roomInfoAddPlayer error: \u0026errors.errorString{s:\"[Service RoomService roomInfoAddPlayer] room can not add player, room state: \\\"Wait\\\"\"}","Success":false}
//enter: {"Data":null,"LogId":"20210304154832821292F8DA","Message":"[Service RoomService EnterRoom] roomInfoAddPlayer error: \u0026errors.errorString{s:\"[Service RoomService roomInfoAddPlayer] room can not add player, room state: \\\"Wait\\\"\"}","Success":false}
//*/
//func TestRoomNewRoomAndEnterRoomMany2(t *testing.T) {
//	// 初始化测试玩家
//	num := 10
//	InitTestPlayerList(num)
//	// 创建房间
//	postUrl := server + "/room/new_room"
//	reqData := map[string]interface{}{
//		"tempId": TestPlayerList[0],
//		"roomSetting": model.RoomSetting{
//			TotalNum: 6,
//			SpyNum:   2,
//			BlankNum: 0,
//		},
//	}
//
//	reqDataJsonString := util.InterfaceMapToJsonString(reqData)
//	reqSign := util.GetSign(reqDataJsonString)
//
//	respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println(respData)
//	}
//
//	roomId := util.JsonStringToInterfaceMap(respData)["Data"].(map[string]interface{})["roomId"].(string)
//	fmt.Println("roomId:", roomId)
//
//	wg := sync.WaitGroup{}
//	wg.Add(num - 1)
//
//	// 玩家加入房间
//	for i := 1; i < num; i++ {
//		i := i
//		go func() {
//			postUrl := server + "/room/enter_room"
//			reqData := map[string]interface{}{
//				"tempId": TestPlayerList[i],
//				"roomId": roomId,
//			}
//
//			reqDataJsonString := util.InterfaceMapToJsonString(reqData)
//			reqSign := util.GetSign(reqDataJsonString)
//
//			fmt.Println("enter(tempId):", i, TestPlayerList[i], reqSign)
//
//			time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)
//			respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//			if err != nil {
//				fmt.Println(err)
//			} else {
//				fmt.Println("enter:", respData)
//			}
//
//			wg.Done()
//		}()
//	}
//
//	wg.Wait()
//}
//
//// 获取房间信息
//func TestRoomRefreshRoom(t *testing.T) {
//	// 初始化测试玩家
//	num := 2
//	InitTestPlayerList(num)
//	// 创建房间
//	postUrl := server + "/room/new_room"
//	reqData := map[string]interface{}{
//		"tempId": TestPlayerList[0],
//		"roomSetting": model.RoomSetting{
//			TotalNum: 6,
//			SpyNum:   2,
//			BlankNum: 0,
//		},
//	}
//
//	reqDataJsonString := util.InterfaceMapToJsonString(reqData)
//	reqSign := util.GetSign(reqDataJsonString)
//
//	respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println(respData)
//	}
//
//	roomId := util.JsonStringToInterfaceMap(respData)["Data"].(map[string]interface{})["roomId"].(string)
//	fmt.Println("roomId:", roomId)
//
//	// 获取房间信息（在房间内）
//	postUrl = server + "/room/refresh_room"
//	reqData = map[string]interface{}{
//		"tempId": TestPlayerList[0],
//	}
//
//	reqDataJsonString = util.InterfaceMapToJsonString(reqData)
//	reqSign = util.GetSign(reqDataJsonString)
//
//	respData, err = util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println("reqDataJsonString:", reqDataJsonString)
//		fmt.Println("respData:", respData)
//	}
//
//	// 获取房间信息（不在房间内）
//	postUrl = server + "/room/refresh_room"
//	reqData = map[string]interface{}{
//		"tempId": TestPlayerList[1],
//	}
//
//	reqDataJsonString = util.InterfaceMapToJsonString(reqData)
//	reqSign = util.GetSign(reqDataJsonString)
//
//	respData, err = util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println(respData)
//	}
//}
//
//// 多人加入并退出房间（10）
//func TestRoomNewRoomAndEnterRoomAndExitRoomMany1(t *testing.T) {
//	// 初始化测试玩家
//	num := 10
//	InitTestPlayerList(num)
//	// 创建房间
//	postUrl := server + "/room/new_room"
//	reqData := map[string]interface{}{
//		"tempId": TestPlayerList[0],
//		"roomSetting": model.RoomSetting{
//			TotalNum: 6,
//			SpyNum:   2,
//			BlankNum: 0,
//		},
//	}
//
//	reqDataJsonString := util.InterfaceMapToJsonString(reqData)
//	reqSign := util.GetSign(reqDataJsonString)
//
//	respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println(respData)
//	}
//
//	roomId := util.JsonStringToInterfaceMap(respData)["Data"].(map[string]interface{})["roomId"].(string)
//	fmt.Println("roomId:", roomId)
//
//	wg := sync.WaitGroup{}
//	wg.Add(2 * (num - 1))
//
//	// 玩家加入房间
//	for i := 1; i < num; i++ {
//		i := i
//		go func() {
//			postUrl := server + "/room/enter_room"
//			reqData := map[string]interface{}{
//				"tempId": TestPlayerList[i],
//				"roomId": roomId,
//			}
//
//			reqDataJsonString := util.InterfaceMapToJsonString(reqData)
//			reqSign := util.GetSign(reqDataJsonString)
//
//			fmt.Println("enter(tempId):", i, TestPlayerList[i], reqSign)
//
//			time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)
//			respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//			if err != nil {
//				fmt.Println(err)
//			} else {
//				fmt.Println("enter:", respData)
//			}
//
//			wg.Done()
//		}()
//	}
//
//	// 玩家退出房间
//	for i := 1; i < num; i++ {
//		i := i
//		go func() {
//			postUrl := server + "/room/exit_room"
//			reqData := map[string]interface{}{
//				"tempId": TestPlayerList[i],
//				"roomId": roomId,
//			}
//
//			reqDataJsonString := util.InterfaceMapToJsonString(reqData)
//			reqSign := util.GetSign(reqDataJsonString)
//
//			fmt.Println("enter(tempId):", i, TestPlayerList[i], reqSign)
//
//			time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)
//			respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//			if err != nil {
//				fmt.Println(err)
//			} else {
//				fmt.Println("exit:", respData)
//				fmt.Println("reqDataJsonString:", reqDataJsonString)
//				fmt.Println("respData:", respData)
//			}
//
//			wg.Done()
//		}()
//	}
//
//	wg.Wait()
//}
//
//// 多人加入并退出房间（100）
//func TestRoomNewRoomAndEnterRoomAndExitRoomMany2(t *testing.T) {
//	// 初始化测试玩家
//	num := 100
//	InitTestPlayerList(num)
//	// 创建房间
//	postUrl := server + "/room/new_room"
//	reqData := map[string]interface{}{
//		"tempId": TestPlayerList[0],
//		"roomSetting": model.RoomSetting{
//			TotalNum: 60,
//			SpyNum:   2,
//			BlankNum: 0,
//		},
//	}
//
//	reqDataJsonString := util.InterfaceMapToJsonString(reqData)
//	reqSign := util.GetSign(reqDataJsonString)
//
//	respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println(respData)
//	}
//
//	roomId := util.JsonStringToInterfaceMap(respData)["Data"].(map[string]interface{})["roomId"].(string)
//	fmt.Println("roomId:", roomId)
//
//	wg := sync.WaitGroup{}
//	wg.Add(2 * (num - 1))
//
//	// 玩家加入房间
//	for i := 1; i < num; i++ {
//		i := i
//		go func() {
//			postUrl := server + "/room/enter_room"
//			reqData := map[string]interface{}{
//				"tempId": TestPlayerList[i],
//				"roomId": roomId,
//			}
//
//			reqDataJsonString := util.InterfaceMapToJsonString(reqData)
//			reqSign := util.GetSign(reqDataJsonString)
//
//			fmt.Println("enter(tempId):", i, TestPlayerList[i], reqSign)
//
//			time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)
//			respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//			if err != nil {
//				fmt.Println(err)
//			} else {
//				fmt.Println("enter:", respData)
//			}
//
//			wg.Done()
//		}()
//	}
//
//	// 玩家退出房间
//	for i := 1; i < num; i++ {
//		i := i
//		go func() {
//			postUrl := server + "/room/exit_room"
//			reqData := map[string]interface{}{
//				"tempId": TestPlayerList[i],
//				"roomId": roomId,
//			}
//
//			reqDataJsonString := util.InterfaceMapToJsonString(reqData)
//			reqSign := util.GetSign(reqDataJsonString)
//
//			fmt.Println("enter(tempId):", i, TestPlayerList[i], reqSign)
//
//			time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)
//			respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//			if err != nil {
//				fmt.Println(err)
//			} else {
//				fmt.Println("exit:", respData)
//			}
//
//			wg.Done()
//		}()
//	}
//
//	wg.Wait()
//}
//
//// 一人加入房间并准备
///**
//enter(tempId): 1 20210305111752600817DDB9 7c4681b429567a4e1913741413fb13cc
//enter: {"Data":null,"LogId":"202103051117526964740855","Message":"Success","Success":true}
//ready(tempId): 1 20210305111752600817DDB9 cb399cbbc5555db0f533ebb3b7d8461b
//ready: {"Data":null,"LogId":"20210305111752751709AA2F","Message":"Success","Success":true}
//*/
//func TestRoomNewRoomAndEnterRoomAndReady(t *testing.T) {
//	// 初始化测试玩家
//	num := 2
//	InitTestPlayerList(num)
//	// 创建房间
//	postUrl := server + "/room/new_room"
//	reqData := map[string]interface{}{
//		"tempId": TestPlayerList[0],
//		"roomSetting": model.RoomSetting{
//			TotalNum: 6,
//			SpyNum:   2,
//			BlankNum: 0,
//		},
//	}
//
//	reqDataJsonString := util.InterfaceMapToJsonString(reqData)
//	reqSign := util.GetSign(reqDataJsonString)
//
//	respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println(respData)
//	}
//
//	roomId := util.JsonStringToInterfaceMap(respData)["Data"].(map[string]interface{})["roomId"].(string)
//	fmt.Println("roomId:", roomId)
//
//	wg := sync.WaitGroup{}
//	wg.Add(num - 1)
//
//	// 玩家加入房间并准备
//	for i := 1; i < num; i++ {
//		i := i
//		go func() {
//			// 玩家加入房间
//			postUrl := server + "/room/enter_room"
//			reqData := map[string]interface{}{
//				"tempId": TestPlayerList[i],
//				"roomId": roomId,
//			}
//
//			reqDataJsonString := util.InterfaceMapToJsonString(reqData)
//			reqSign := util.GetSign(reqDataJsonString)
//
//			fmt.Println("enter(tempId):", i, TestPlayerList[i], reqSign)
//
//			//time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)
//			respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//			if err != nil {
//				fmt.Println(err)
//			} else {
//				fmt.Println("enter:", respData)
//			}
//
//			// 玩家准备
//			postUrl = server + "/game/ready_game"
//			reqData = map[string]interface{}{
//				"tempId": TestPlayerList[i],
//			}
//
//			reqDataJsonString = util.InterfaceMapToJsonString(reqData)
//			reqSign = util.GetSign(reqDataJsonString)
//
//			fmt.Println("ready(tempId):", i, TestPlayerList[i], reqSign)
//
//			//time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)
//			respData, err = util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//			if err != nil {
//				fmt.Println(err)
//			} else {
//				fmt.Println("ready:", respData)
//			}
//
//			wg.Done()
//		}()
//	}
//
//	wg.Wait()
//}
//
//// 多人加入房间并准备
//func TestRoomNewRoomAndEnterRoomAndReadyMany1(t *testing.T) {
//	// 初始化测试玩家
//	num := 10
//	InitTestPlayerList(num)
//	// 创建房间
//	postUrl := server + "/room/new_room"
//	reqData := map[string]interface{}{
//		"tempId": TestPlayerList[0],
//		"roomSetting": model.RoomSetting{
//			TotalNum: 6,
//			SpyNum:   2,
//			BlankNum: 0,
//		},
//	}
//
//	reqDataJsonString := util.InterfaceMapToJsonString(reqData)
//	reqSign := util.GetSign(reqDataJsonString)
//
//	respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println(respData)
//	}
//
//	roomId := util.JsonStringToInterfaceMap(respData)["Data"].(map[string]interface{})["roomId"].(string)
//	fmt.Println("roomId:", roomId)
//
//	wg := sync.WaitGroup{}
//	wg.Add(num - 1)
//
//	// 玩家加入房间并准备
//	for i := 1; i < num; i++ {
//		i := i
//		go func() {
//			// 玩家加入房间
//			postUrl := server + "/room/enter_room"
//			reqData := map[string]interface{}{
//				"tempId": TestPlayerList[i],
//				"roomId": roomId,
//			}
//
//			reqDataJsonString := util.InterfaceMapToJsonString(reqData)
//			reqSign := util.GetSign(reqDataJsonString)
//
//			fmt.Println("enter(tempId):", i, TestPlayerList[i], reqSign)
//
//			time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)
//			respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//			if err != nil {
//				fmt.Println(err)
//			} else {
//				fmt.Println("enter:", respData)
//			}
//
//			// 玩家准备
//			postUrl = server + "/game/ready_game"
//			reqData = map[string]interface{}{
//				"tempId": TestPlayerList[i],
//			}
//
//			reqDataJsonString = util.InterfaceMapToJsonString(reqData)
//			reqSign = util.GetSign(reqDataJsonString)
//
//			fmt.Println("ready(tempId):", i, TestPlayerList[i], reqSign)
//
//			time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)
//			respData, err = util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//			if err != nil {
//				fmt.Println(err)
//			} else {
//				fmt.Println("ready:", respData)
//			}
//
//			wg.Done()
//		}()
//	}
//
//	wg.Wait()
//
//	fmt.Println("roomId:", roomId)
//}
//
//// 多人加入房间并准备，后房主开始游戏
//func TestRoomNewRoomAndEnterRoomAndReadyAndStart(t *testing.T) {
//	// 初始化测试玩家
//	num := 10
//	InitTestPlayerList(num)
//	// 创建房间
//	postUrl := server + "/room/new_room"
//	reqData := map[string]interface{}{
//		"tempId": TestPlayerList[0],
//		"roomSetting": model.RoomSetting{
//			TotalNum: num,
//			SpyNum:   2,
//			BlankNum: 0,
//		},
//	}
//
//	reqDataJsonString := util.InterfaceMapToJsonString(reqData)
//	reqSign := util.GetSign(reqDataJsonString)
//
//	respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println(respData)
//	}
//
//	roomId := util.JsonStringToInterfaceMap(respData)["Data"].(map[string]interface{})["roomId"].(string)
//	fmt.Println("roomId:", roomId)
//
//	wg := sync.WaitGroup{}
//	wg.Add(num - 1)
//
//	// 玩家加入房间并准备
//	for i := 1; i < num; i++ {
//		i := i
//		go func() {
//			// 玩家加入房间
//			postUrl := server + "/room/enter_room"
//			reqData := map[string]interface{}{
//				"tempId": TestPlayerList[i],
//				"roomId": roomId,
//			}
//
//			reqDataJsonString := util.InterfaceMapToJsonString(reqData)
//			reqSign := util.GetSign(reqDataJsonString)
//
//			fmt.Println("enter(tempId):", i, TestPlayerList[i], reqSign)
//
//			time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)
//			respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//			if err != nil {
//				fmt.Println(err)
//			} else {
//				fmt.Println("enter:", respData)
//			}
//
//			// 玩家准备
//			postUrl = server + "/game/ready_game"
//			reqData = map[string]interface{}{
//				"tempId": TestPlayerList[i],
//			}
//
//			reqDataJsonString = util.InterfaceMapToJsonString(reqData)
//			reqSign = util.GetSign(reqDataJsonString)
//
//			fmt.Println("ready(tempId):", i, TestPlayerList[i], reqSign)
//
//			time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)
//			respData, err = util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//			if err != nil {
//				fmt.Println(err)
//			} else {
//				fmt.Println("ready:", respData)
//			}
//
//			wg.Done()
//		}()
//	}
//
//	wg.Wait()
//
//	// 房主开始游戏
//	postUrl = server + "/game/start_game"
//	reqData = map[string]interface{}{
//		"tempId": TestPlayerList[0],
//	}
//
//	reqDataJsonString = util.InterfaceMapToJsonString(reqData)
//	reqSign = util.GetSign(reqDataJsonString)
//
//	fmt.Println("start(tempId):", TestPlayerList[0], reqSign)
//
//	time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)
//	respData, err = util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println("start:", respData)
//	}
//
//	fmt.Println("roomId:", roomId)
//}
//
//// 多人加入房间并准备，后房主开始游戏，结束游戏，重开游戏，综合测试
//func Test_All_1(t *testing.T) {
//	// 初始化测试玩家
//	num := 10
//	InitTestPlayerList(num)
//	// 创建房间
//	postUrl := server + "/room/new_room"
//	reqData := map[string]interface{}{
//		"tempId": TestPlayerList[0],
//		"roomSetting": model.RoomSetting{
//			TotalNum: num,
//			SpyNum:   2,
//			BlankNum: 0,
//		},
//	}
//
//	reqDataJsonString := util.InterfaceMapToJsonString(reqData)
//	reqSign := util.GetSign(reqDataJsonString)
//
//	respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println(respData)
//	}
//
//	roomId := util.JsonStringToInterfaceMap(respData)["Data"].(map[string]interface{})["roomId"].(string)
//	fmt.Println("roomId:", roomId)
//
//	wg := sync.WaitGroup{}
//	wg.Add(num - 1)
//
//	// 玩家加入房间并准备
//	for i := 1; i < num; i++ {
//		i := i
//		go func() {
//			// 玩家加入房间
//			postUrl := server + "/room/enter_room"
//			reqData := map[string]interface{}{
//				"tempId": TestPlayerList[i],
//				"roomId": roomId,
//			}
//
//			reqDataJsonString := util.InterfaceMapToJsonString(reqData)
//			reqSign := util.GetSign(reqDataJsonString)
//
//			fmt.Println("enter(tempId):", i, TestPlayerList[i], reqSign)
//
//			time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)
//			respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//			if err != nil {
//				fmt.Println(err)
//			} else {
//				fmt.Println("enter:", respData)
//			}
//
//			// 玩家准备
//			postUrl = server + "/game/ready_game"
//			reqData = map[string]interface{}{
//				"tempId": TestPlayerList[i],
//			}
//
//			reqDataJsonString = util.InterfaceMapToJsonString(reqData)
//			reqSign = util.GetSign(reqDataJsonString)
//
//			fmt.Println("ready(tempId):", i, TestPlayerList[i], reqSign)
//
//			time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)
//			respData, err = util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//			if err != nil {
//				fmt.Println(err)
//			} else {
//				fmt.Println("ready:", respData)
//			}
//
//			wg.Done()
//		}()
//	}
//
//	wg.Wait()
//
//	for i := 0; i < 7; i++ {
//		// 房主开始游戏
//		postUrl = server + "/game/start_game"
//		reqData = map[string]interface{}{
//			"tempId": TestPlayerList[0],
//		}
//
//		reqDataJsonString = util.InterfaceMapToJsonString(reqData)
//		reqSign = util.GetSign(reqDataJsonString)
//
//		fmt.Println("start(tempId):", TestPlayerList[0], reqSign)
//
//		time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)
//		respData, err = util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//		if err != nil {
//			fmt.Println(err)
//		} else {
//			fmt.Println("start:", respData)
//		}
//
//		// 房主结束游戏
//		postUrl = server + "/game/end_game"
//		reqData = map[string]interface{}{
//			"tempId": TestPlayerList[0],
//		}
//
//		reqDataJsonString = util.InterfaceMapToJsonString(reqData)
//		reqSign = util.GetSign(reqDataJsonString)
//
//		fmt.Println("end(tempId):", TestPlayerList[0], reqSign)
//
//		time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)
//		respData, err = util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//		if err != nil {
//			fmt.Println(err)
//		} else {
//			fmt.Println("end:", respData)
//		}
//	}
//
//	// 房主开始游戏
//	postUrl = server + "/game/start_game"
//	reqData = map[string]interface{}{
//		"tempId": TestPlayerList[0],
//	}
//
//	reqDataJsonString = util.InterfaceMapToJsonString(reqData)
//	reqSign = util.GetSign(reqDataJsonString)
//
//	fmt.Println("start(tempId):", TestPlayerList[0], reqSign)
//
//	time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)
//	respData, err = util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println("start:", respData)
//	}
//
//	// 房主重新开始游戏
//	for i := 0; i < 7; i++ {
//		postUrl = server + "/game/restart_game"
//		reqData = map[string]interface{}{
//			"tempId": TestPlayerList[0],
//		}
//
//		reqDataJsonString = util.InterfaceMapToJsonString(reqData)
//		reqSign = util.GetSign(reqDataJsonString)
//
//		fmt.Println("restart(tempId):", TestPlayerList[0], reqSign)
//
//		time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)
//		respData, err = util.HttpPost(postUrl, reqDataJsonString, reqSign)
//
//		if err != nil {
//			fmt.Println(err)
//		} else {
//			fmt.Println("restart:", respData)
//		}
//	}
//
//	fmt.Println("roomId:", roomId)
//}
