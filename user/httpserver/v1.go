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
	userRoutesV1.POST("/register", middleware.LogGlobal.Middleware(), userHandler.Register)

	// R : Read
	userRoutesV1.POST("/login", middleware.LogGlobal.Middleware(), userHandler.Login)
	userRoutesV1.GET("/get", middleware.LogGlobal.Middleware(), middleware.AuthenticationLOM, middleware.Authorization(map[string]bool{
		"owner":      true,
		"superadmin": true,
		"admin":      true,
	}), userHandler.GetUsers)

	// U : Update
	userRoutesV1.PATCH("/change.password", middleware.LogGlobal.Middleware(), middleware.AuthenticationLOM, userHandler.ChangePassword)
	userRoutesV1.PATCH("/add.chatroom.uuid", middleware.LogGlobal.Middleware(), middleware.AuthenticationLOM, userHandler.AddChatRoomUUID)

	// D : Delete
	userRoutesV1.DELETE("/delete", middleware.LogGlobal.Middleware(), middleware.AuthenticationLOM, middleware.Authorization(map[string]bool{
		"owner":      true,
		"superadmin": true,
	}), userHandler.DeleteUser)

	// Health Check
	userRoutesV1.GET("/health", middleware.LogGlobal.Middleware(), utils.HealthCheck())
}
