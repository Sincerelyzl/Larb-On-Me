package httpserver

import (
	"github.com/Sincerelyzl/larb-on-me/chat-room/handler"
	"github.com/gin-gonic/gin"
)

func setupUserRouteV1(r *gin.Engine, chatRoomHandler handler.ChatRoomHandler) {
	chatRoomRoutesV1 := r.Group("/v1/chatroom")

	chatRoomRoutesV1.POST("/create", chatRoomHandler.CreateChatRoom)
}
