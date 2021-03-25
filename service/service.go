package service

var (
	CoreUserService *UserService // 核心用户服务
	CoreRoomService *RoomService // 核心房间服务
	CoreGameService *GameService // 核心游戏服务
)

func InitService() {
	CoreUserService = UserServiceInstance()
	CoreRoomService = RoomServiceInstance()
	CoreGameService = GameServiceInstance()
}
