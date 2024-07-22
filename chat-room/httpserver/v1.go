package httpserver

import (
	"github.com/Sincerelyzl/larb-on-me/chat-room/handler"
	"github.com/Sincerelyzl/larb-on-me/common/middleware"
	"github.com/Sincerelyzl/larb-on-me/common/utils"
	"github.com/gin-gonic/gin"
)

func setupUserRouteV1(r *gin.Engine, chatRoomHandler handler.ChatRoomHandler) {
	chatRoomRoutesV1 := r.Group("/v1/chatroom")

	chatRoomRoutesV1.POST("/create", middleware.LogGlobal.Middleware(), middleware.AuthenticationLOM, chatRoomHandler.CreateChatRoom)

	// Health Check
	chatRoomRoutesV1.GET("/health", middleware.LogGlobal.Middleware(), utils.HealthCheck())
}
