package httpserver

import (
	"github.com/Sincerelyzl/larb-on-me/common/middleware"
	"github.com/Sincerelyzl/larb-on-me/common/utils"
	"github.com/Sincerelyzl/larb-on-me/user/handler"
	"github.com/gin-gonic/gin"
)

func setupUserRouteV1(r *gin.Engine, userHandler handler.UserHandler) {
	userRoutesV1 := r.Group("/v1/user")

	// C : Create
	userRoutesV1.POST("/register", userHandler.Register)

	// R : Read
	userRoutesV1.POST("/login", userHandler.Login)
	userRoutesV1.GET("/get", middleware.AuthenticationLOM, middleware.Authorization(map[string]bool{
		"owner":      true,
		"superadmin": true,
		"admin":      true,
	}), userHandler.GetUsers)

	// U : Update
	userRoutesV1.PATCH("/change.password", middleware.AuthenticationLOM, userHandler.ChangePassword)

	// D : Delete
	userRoutesV1.DELETE("/delete", middleware.AuthenticationLOM, middleware.Authorization(map[string]bool{
		"owner":      true,
		"superadmin": true,
	}), userHandler.DeleteUser)

	userRoutesV1.GET("/health", utils.HealthCheck())
}
