package httpserver

import (
	"github.com/Sincerelyzl/larb-on-me/user/handler"
	"github.com/gin-gonic/gin"
)

type server struct {
	Router      *gin.Engine
	userHandler handler.UserHandler
}

func NewHTTPServer(userHandler handler.UserHandler) *server {
	r := initRouter(userHandler)
	return &server{
		Router:      r,
		userHandler: userHandler,
	}
}

func (s *server) Run(port string) error {
	return s.Router.Run(":" + port)
}
