package httpserver

import (
	"github.com/Sincerelyzl/larb-on-me/user/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func initRouter(userHandler handler.UserHandler) *gin.Engine {
	r := gin.New()
	r.SetTrustedProxies(nil)
	corsConfig := cors.New(
		cors.Config{
			AllowCredentials: true,
			AllowOrigins:     []string{"https://localhost:3000", "https://localhost:3008", "https://localhost:3009"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Authorization", "Content-Type", "X-Correlation-ID"},
		},
	)
	r.Use(corsConfig)
	// TODO here Use middleware log
	setupUserRouteV1(r, userHandler)
	return r
}
