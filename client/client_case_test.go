package client

import (
	"fmt"
	logs "lmf.mortal.com/GoLogs"
	"lmf.mortal.com/GoWxWhoIsTheSpy/config"
	"lmf.mortal.com/GoWxWhoIsTheSpy/model"
	"lmf.mortal.com/GoWxWhoIsTheSpy/util"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var TestPlayerList []string // 测试玩家的临时 ID 集合

//var server = "http://something:9205" // 测试环境
var server = "http://127.0.0.1:9205" // 开发环境
//var server = "https://coding8zz.com" // 线上环境

func init() {
	config.InitConfig("../conf/who_is_the_spy_dev.json")
	logs.InitDefaultLogger(config.ConfigInstance.LogConfig)
}

// 初始化随机玩家
func InitTestPlayerList(num int) {
	for i := 0; i < num; i++ {
		// 用户登录
		postUrl := server + "/user/wx_login"
		reqData := map[string]interface{}{
			"code": "test_code",
		}

		reqDataJsonString := util.InterfaceMapToJsonString(reqData)
		reqSign := util.GetSign(reqDataJsonString)

		respData, _ := util.HttpPost(postUrl, reqDataJsonString, reqSign)

		tempId := util.JsonStringToStringMap(respData)["Data"]

		// 用户信息
		postUrl = server + "/user/update_userInfo"
		reqData = map[string]interface{}{
			"tempId": tempId,
			"userInfo": model.User{
				NickName:   "test_nick_name_" + fmt.Sprint(i),
				Gender:     0,
				AvatarUrl:  "test_url_" + fmt.Sprint(i),
				City:       "test_city_" + fmt.Sprint(i),
				Province:   "test_province_" + fmt.Sprint(i),
				Country:    "test_country_" + fmt.Sprint(i),
				Language:   "test_language_" + fmt.Sprint(i),
				CreateTime: time.Time{},
				UpdateTime: time.Time{},
			},
		}

		reqDataJsonString = util.InterfaceMapToJsonString(reqData)
		reqSign = util.GetSign(reqDataJsonString)

		respData, _ = util.HttpPost(postUrl, reqDataJsonString, reqSign)
		fmt.Println("init:", respData)

		// 添加用户到测试列表
		TestPlayerList = append(TestPlayerList, tempId)
	}

	fmt.Println("TestPlayerList:", TestPlayerList)
}

// 测试用户登录和创建用户信息
func TestUserWxLoginAndUpdateUserInfo(t *testing.T) {
	InitTestPlayerList(10)
}

// 多人加入房间并准备，后房主开始游戏，结束游戏，重开游戏，综合测试
func Test_All_1(t *testing.T) {
	// 初始化测试玩家
	num := 10
	InitTestPlayerList(num)
	// 创建房间
	postUrl := server + "/room/new_room"
	reqData := map[string]interface{}{
		"tempId": TestPlayerList[0],
		"roomSetting": model.RoomSetting{
			TotalNum: num,
			SpyNum:   2,
			BlankNum: 0,
		},
	}

	reqDataJsonString := util.InterfaceMapToJsonString(reqData)
	reqSign := util.GetSign(reqDataJsonString)

	respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(respData)
	}

	roomId := util.JsonStringToInterfaceMap(respData)["Data"].(map[string]interface{})["roomId"].(string)
	fmt.Println("roomId:", roomId)

	wg := sync.WaitGroup{}
	wg.Add(num - 1)

	// 玩家加入房间并准备
	for i := 1; i < num; i++ {
		i := i
		go func() {
			// 玩家加入房间
			postUrl := server + "/room/enter_room"
			reqData := map[string]interface{}{
				"tempId": TestPlayerList[i],
				"roomId": roomId,
			}

			reqDataJsonString := util.InterfaceMapToJsonString(reqData)
			reqSign := util.GetSign(reqDataJsonString)

			fmt.Println("enter(tempId):", i, TestPlayerList[i], reqSign)

			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)

			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("enter:", respData)
			}

			// 玩家准备
			postUrl = server + "/game/ready_game"
			reqData = map[string]interface{}{
				"tempId": TestPlayerList[i],
			}

			reqDataJsonString = util.InterfaceMapToJsonString(reqData)
			reqSign = util.GetSign(reqDataJsonString)

			fmt.Println("ready(tempId):", i, TestPlayerList[i], reqSign)

			time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
			respData, err = util.HttpPost(postUrl, reqDataJsonString, reqSign)

			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("ready:", respData)
			}

			wg.Done()
		}()
	}

	wg.Wait()

	for i := 0; i < 7; i++ {
		// 房主开始游戏
		postUrl = server + "/game/start_game"
		reqData = map[string]interface{}{
			"tempId": TestPlayerList[0],
		}

		reqDataJsonString = util.InterfaceMapToJsonString(reqData)
		reqSign = util.GetSign(reqDataJsonString)

		fmt.Println("start(tempId):", TestPlayerList[0], reqSign)

		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		respData, err = util.HttpPost(postUrl, reqDataJsonString, reqSign)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("start:", respData)
		}

		// 房主结束游戏
		postUrl = server + "/game/end_game"
		reqData = map[string]interface{}{
			"tempId": TestPlayerList[0],
		}

		reqDataJsonString = util.InterfaceMapToJsonString(reqData)
		reqSign = util.GetSign(reqDataJsonString)

		fmt.Println("end(tempId):", TestPlayerList[0], reqSign)

		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		respData, err = util.HttpPost(postUrl, reqDataJsonString, reqSign)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("end:", respData)
		}
	}

	// 房主开始游戏
	postUrl = server + "/game/start_game"
	reqData = map[string]interface{}{
		"tempId": TestPlayerList[0],
	}

	reqDataJsonString = util.InterfaceMapToJsonString(reqData)
	reqSign = util.GetSign(reqDataJsonString)

	fmt.Println("start(tempId):", TestPlayerList[0], reqSign)

	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	respData, err = util.HttpPost(postUrl, reqDataJsonString, reqSign)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("start:", respData)
	}

	// 房主重新开始游戏
	for i := 0; i < 7; i++ {
		postUrl = server + "/game/restart_game"
		reqData = map[string]interface{}{
			"tempId": TestPlayerList[0],
		}

		reqDataJsonString = util.InterfaceMapToJsonString(reqData)
		reqSign = util.GetSign(reqDataJsonString)

		fmt.Println("restart(tempId):", TestPlayerList[0], reqSign)

		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		respData, err = util.HttpPost(postUrl, reqDataJsonString, reqSign)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("restart:", respData)
		}
	}

	fmt.Println("roomId:", roomId)
}

func BenchmarkTest_All_1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// 初始化测试玩家
		num := 10
		InitTestPlayerList(num)
		// 创建房间
		postUrl := server + "/room/new_room"
		reqData := map[string]interface{}{
			"tempId": TestPlayerList[0],
			"roomSetting": model.RoomSetting{
				TotalNum: num,
				SpyNum:   2,
				BlankNum: 0,
			},
		}

		reqDataJsonString := util.InterfaceMapToJsonString(reqData)
		reqSign := util.GetSign(reqDataJsonString)

		respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(respData)
		}

		roomId := util.JsonStringToInterfaceMap(respData)["Data"].(map[string]interface{})["roomId"].(string)
		fmt.Println("roomId:", roomId)

		wg := sync.WaitGroup{}
		wg.Add(num - 1)

		// 玩家加入房间并准备
		for i := 1; i < num; i++ {
			i := i
			go func() {
				// 玩家加入房间
				postUrl := server + "/room/enter_room"
				reqData := map[string]interface{}{
					"tempId": TestPlayerList[i],
					"roomId": roomId,
				}

				reqDataJsonString := util.InterfaceMapToJsonString(reqData)
				reqSign := util.GetSign(reqDataJsonString)

				fmt.Println("enter(tempId):", i, TestPlayerList[i], reqSign)

				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
				respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)

				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("enter:", respData)
				}

				// 玩家准备
				postUrl = server + "/game/ready_game"
				reqData = map[string]interface{}{
					"tempId": TestPlayerList[i],
				}

				reqDataJsonString = util.InterfaceMapToJsonString(reqData)
				reqSign = util.GetSign(reqDataJsonString)

				fmt.Println("ready(tempId):", i, TestPlayerList[i], reqSign)

				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
				respData, err = util.HttpPost(postUrl, reqDataJsonString, reqSign)

				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("ready:", respData)
				}

				wg.Done()
			}()
		}

		wg.Wait()

		for i := 0; i < 7; i++ {
			// 房主开始游戏
			postUrl = server + "/game/start_game"
			reqData = map[string]interface{}{
				"tempId": TestPlayerList[0],
			}

			reqDataJsonString = util.InterfaceMapToJsonString(reqData)
			reqSign = util.GetSign(reqDataJsonString)

			fmt.Println("start(tempId):", TestPlayerList[0], reqSign)

			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			respData, err = util.HttpPost(postUrl, reqDataJsonString, reqSign)

			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("start:", respData)
			}

			// 房主结束游戏
			postUrl = server + "/game/end_game"
			reqData = map[string]interface{}{
				"tempId": TestPlayerList[0],
			}

			reqDataJsonString = util.InterfaceMapToJsonString(reqData)
			reqSign = util.GetSign(reqDataJsonString)

			fmt.Println("end(tempId):", TestPlayerList[0], reqSign)

			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			respData, err = util.HttpPost(postUrl, reqDataJsonString, reqSign)

			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("end:", respData)
			}
		}

		// 房主开始游戏
		postUrl = server + "/game/start_game"
		reqData = map[string]interface{}{
			"tempId": TestPlayerList[0],
		}

		reqDataJsonString = util.InterfaceMapToJsonString(reqData)
		reqSign = util.GetSign(reqDataJsonString)

		fmt.Println("start(tempId):", TestPlayerList[0], reqSign)

		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		respData, err = util.HttpPost(postUrl, reqDataJsonString, reqSign)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("start:", respData)
		}

		// 房主重新开始游戏
		for i := 0; i < 7; i++ {
			postUrl = server + "/game/restart_game"
			reqData = map[string]interface{}{
				"tempId": TestPlayerList[0],
			}

			reqDataJsonString = util.InterfaceMapToJsonString(reqData)
			reqSign = util.GetSign(reqDataJsonString)

			fmt.Println("restart(tempId):", TestPlayerList[0], reqSign)

			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			respData, err = util.HttpPost(postUrl, reqDataJsonString, reqSign)

			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("restart:", respData)
			}
		}

		fmt.Println("roomId:", roomId)
	}
}

// 多人加入房间并退出房间
func Test_All_2(t *testing.T) {
	// 初始化测试玩家
	num := 10
	InitTestPlayerList(num)
	// 创建房间
	postUrl := server + "/room/new_room"
	reqData := map[string]interface{}{
		"tempId": TestPlayerList[0],
		"roomSetting": model.RoomSetting{
			TotalNum: num,
			SpyNum:   2,
			BlankNum: 0,
		},
	}

	reqDataJsonString := util.InterfaceMapToJsonString(reqData)
	reqSign := util.GetSign(reqDataJsonString)

	respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(respData)
	}

	roomId := util.JsonStringToInterfaceMap(respData)["Data"].(map[string]interface{})["roomId"].(string)
	fmt.Println("roomId:", roomId)

	wg := sync.WaitGroup{}
	wg.Add(num - 1)

	// 玩家加入房间
	for i := 1; i < num; i++ {
		i := i
		go func() {
			// 玩家加入房间
			postUrl := server + "/room/enter_room"
			reqData := map[string]interface{}{
				"tempId": TestPlayerList[i],
				"roomId": roomId,
			}

			reqDataJsonString := util.InterfaceMapToJsonString(reqData)
			reqSign := util.GetSign(reqDataJsonString)

			fmt.Println("enter(tempId):", i, TestPlayerList[i], reqSign)

			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)

			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("enter:", respData)
			}

			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("roomId:", roomId)
	time.Sleep(time.Second * 5)

	wg = sync.WaitGroup{}
	wg.Add(num)

	// 玩家退出房间
	for i := 0; i < num; i++ {
		i := i
		go func() {
			// 玩家退出房间
			postUrl := server + "/room/exit_room"
			reqData := map[string]interface{}{
				"tempId": TestPlayerList[i],
			}

			reqDataJsonString := util.InterfaceMapToJsonString(reqData)
			reqSign := util.GetSign(reqDataJsonString)

			fmt.Println("exit(tempId):", i, TestPlayerList[i], reqSign)

			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)

			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("exit:", respData)
			}

			wg.Done()
		}()
	}

	wg.Wait()

	fmt.Println("roomId:", roomId)
}

func TestRefreshRoom(t *testing.T) {
	// 初始化测试玩家
	num := 1
	InitTestPlayerList(num)
	// 创建房间
	postUrl := server + "/room/refresh_room"
	reqData := map[string]interface{}{
		"tempId": TestPlayerList[0],
	}

	reqDataJsonString := util.InterfaceMapToJsonString(reqData)
	reqSign := util.GetSign(reqDataJsonString)

	respData, err := util.HttpPost(postUrl, reqDataJsonString, reqSign)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(respData)
	}
}