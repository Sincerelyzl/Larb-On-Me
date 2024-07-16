package handler

import (
	"net/http"

	"github.com/Sincerelyzl/larb-on-me/common/middleware"
	"github.com/Sincerelyzl/larb-on-me/common/models"
	"github.com/Sincerelyzl/larb-on-me/common/utils"
	userModel "github.com/Sincerelyzl/larb-on-me/user/models"
	"github.com/gin-gonic/gin"
)

func (h *userHandler) Login(c *gin.Context) {
	var reqBody userModel.UserLoginRequest

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		errResponse := utils.NewErrorResponse(http.StatusBadRequest, "invalid request body")
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	user, err := h.userUsecase.Login(c, reqBody.Username, reqBody.Password)
	if err != nil {
		errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	var userResponse models.UserResponse

	uuidString, errParsingUuid := utils.UuidV7ToString(user.Uuid)
	if errParsingUuid != nil {
		errResponse := utils.NewErrorResponse(http.StatusInternalServerError, errParsingUuid.Error())
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}
	userResponse.Uuid = uuidString
	userResponse.Username = user.Username
	userResponse.CreatedAt = user.CreatedAt
	userResponse.UpdatedAt = user.UpdatedAt

	for _, chatRoomUuid := range user.ChatRoomsUuid {
		chatRoomUuidString, err := utils.UuidV7ToString(chatRoomUuid)
		if err != nil {
			errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
			c.JSON(errResponse.StatusCode, errResponse)
			return
		}
		userResponse.ChatRoomsUuid = append(userResponse.ChatRoomsUuid, chatRoomUuidString)
	}

	lomUser := models.LOMUser{
		Username: userResponse.Username,
		Role:     "user",
	}
	lomToken, err := middleware.GenerateLOMKeys(lomUser)
	if err != nil {
		errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	c.SetCookie(middleware.LOMCookieAuthPrefix, lomToken, int(middleware.LOMExpireTime), "/", "localhost", false, true)
	c.JSON(http.StatusOK, userResponse)
}

func (h *userHandler) Register(c *gin.Context) {

	var reqBody models.UserRegisterRequest

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		errResponse := utils.NewErrorResponse(http.StatusBadRequest, "invalid request body")
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
	userResponse.CreatedAt = user.CreatedAt
	userResponse.UpdatedAt = user.UpdatedAt

	for _, chatRoomUuid := range user.ChatRoomsUuid {
		chatRoomUuidString, err := utils.UuidV7ToString(chatRoomUuid)
		if err != nil {
			errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
			c.JSON(errResponse.StatusCode, errResponse)
			return
		}
		userResponse.ChatRoomsUuid = append(userResponse.ChatRoomsUuid, chatRoomUuidString)
	}

	lomToken, err := middleware.GenerateLOMKeys(userResponse)
	if err != nil {
		errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	c.SetCookie(middleware.LOMCookieAuthPrefix, lomToken, int(middleware.LOMExpireTime), "/", "localhost", false, true)
	c.JSON(http.StatusOK, userResponse)
}

func (h *userHandler) ChangePassword(c *gin.Context) {
	// reqBody := userModel.UserChangePasswordRequest{}

}
