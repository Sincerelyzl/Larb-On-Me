package httpserver

import (
	"github.com/Sincerelyzl/larb-on-me/common/middleware"
	"github.com/Sincerelyzl/larb-on-me/user/handler"
	"github.com/gin-gonic/gin"
)

func setupUserRouteV1(r *gin.Engine, userHandler handler.UserHandler) {
	userRoutesV1 := r.Group("/v1/user")

	userRoutesV1.POST("/register", userHandler.Register)
	userRoutesV1.POST("/login", userHandler.Login)
	userRoutesV1.PATCH("/change.password", middleware.AuthenticationLOM, userHandler.ChangePassword)
}
