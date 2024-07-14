package httpserver

import (
	"github.com/Sincerelyzl/larb-on-me/user/handler"
	"github.com/gin-gonic/gin"
)

func setupUserRouteV1(r *gin.Engine, userHandler handler.UserHandler) {
	userRoutesV1 := r.Group("/v1")

	userRoutesV1.POST("/register", userHandler.Register)
}
