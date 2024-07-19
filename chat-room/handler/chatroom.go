package handler

import (
	"github.com/Sincerelyzl/larb-on-me/chat-room/usecase"
	"github.com/gin-gonic/gin"
)

type ChatRoomHandler interface {
	CreateChatRoom(c *gin.Context)
	JoinChatRoomByJoinCode(c *gin.Context)
	GetChatRoomsByUser(c *gin.Context)
	LeaveChatRoom(c *gin.Context)
	DeleteChatRoom(c *gin.Context)
}

type chatRoomHandler struct {
	chatRoomUsecase usecase.ChatRoomUsecase
}

func NewChatRoomHandler(chatroomUsecase usecase.ChatRoomUsecase) ChatRoomHandler {
	return &chatRoomHandler{
		chatRoomUsecase: chatroomUsecase,
	}
}
