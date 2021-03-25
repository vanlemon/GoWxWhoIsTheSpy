package handler

import (
	"github.com/gin-gonic/gin"
)

// HTTP 服务接口，所有 HTTP handler 都需要实现这个接口
type Handler interface {
	Register(engine *gin.Engine) // 将 HTTP 服务注册到 gin 服务器上
}

// 初始化 HTTP 服务，即将所有服务注册到 gin 服务器上
func InitHandler(engine *gin.Engine) {
	handlers := []Handler{
		NewLoginHandler(),
		NewRoomHandler(),
		NewGameHandler(),
	}

	// 遍历服务列表，注册到 HTTP 服务器上
	for _, handler := range handlers {
		handler.Register(engine)
	}
}
