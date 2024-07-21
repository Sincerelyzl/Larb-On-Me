package handler

import (
	"github.com/Sincerelyzl/larb-on-me/user/usecase"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	ChangePassword(c *gin.Context)
	DeleteUser(c *gin.Context)
	GetUsers(c *gin.Context)
	AddChatRoomUUID(c *gin.Context)
}

type userHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) UserHandler {
	return &userHandler{
		userUsecase: userUsecase,
	}
}
