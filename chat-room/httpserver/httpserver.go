package httpserver

import (
	"github.com/Sincerelyzl/larb-on-me/chat-room/handler"
	"github.com/gin-gonic/gin"
)

type server struct {
	Router          *gin.Engine
	chatRoomHandler handler.ChatRoomHandler
}

func NewHTTPServer(chatRoomHandler handler.ChatRoomHandler) *server {
	r := initRouter(chatRoomHandler)
	return &server{
		Router:          r,
		chatRoomHandler: chatRoomHandler,
	}
}

func (s *server) Run(port string) error {
	return s.Router.Run(port)
}
