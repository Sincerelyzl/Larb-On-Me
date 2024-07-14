package handler

import (
	"net/http"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"github.com/Sincerelyzl/larb-on-me/common/utils"
	"github.com/Sincerelyzl/larb-on-me/user/usecase"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Register(c *gin.Context)
}

type userHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) UserHandler {
	return &userHandler{
		userUsecase: userUsecase,
	}
}

func (h *userHandler) Register(c *gin.Context) {

	var reqBody models.UserRegisterRequest

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		errResponse := utils.NewErrorResponse(http.StatusBadRequest, "Invalid request body")
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	user, err := h.userUsecase.CreateNewUser(c, reqBody.Username, reqBody.Password)
	if err != nil {
		errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	var userResponse models.UserResponse

	uuidString, err := utils.UuidV7ToString(user.Uuid)
	if err != nil {
		errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}
	userResponse.Uuid = uuidString
	userResponse.Username = user.Username

	userResponse.ChatRoomsUuid = make([]string, 0)
	for _, chatRoomUuid := range user.ChatRoomsUuid {
		chatRoomUuidString, err := utils.UuidV7ToString(chatRoomUuid)
		if err != nil {
			errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
			c.JSON(errResponse.StatusCode, errResponse)
			return
		}
		userResponse.ChatRoomsUuid = append(userResponse.ChatRoomsUuid, chatRoomUuidString)
	}

	c.JSON(http.StatusOK, userResponse)
}
