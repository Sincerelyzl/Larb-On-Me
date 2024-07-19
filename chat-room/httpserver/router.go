package httpserver

import (
	"github.com/Sincerelyzl/larb-on-me/chat-room/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func initRouter(chatRoomHander handler.ChatRoomHandler) *gin.Engine {
	r := gin.New()

	corsConfig := cors.New(
		cors.Config{
			AllowCredentials: true,
			AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3008", "http://localhost:3009"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Authorization", "Content-Type", "X-Correlation-ID"},
		},
	)
	r.Use(corsConfig)
	// TODO here Use middleware log
	setupUserRouteV1(r, chatRoomHander)
	return r
}
