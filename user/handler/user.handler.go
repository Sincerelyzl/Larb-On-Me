package handler

import (
	"fmt"
	"net/http"

	"github.com/Sincerelyzl/larb-on-me/common/middleware"
	"github.com/Sincerelyzl/larb-on-me/common/models"
	"github.com/Sincerelyzl/larb-on-me/common/utils"
	"github.com/gin-gonic/gin"
)

func (h *userHandler) Login(c *gin.Context) {

	var reqBody models.UserLoginRequest
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
	userResponse.PermissionGroup = user.PermissionGroup
	userResponse.CreatedAt = user.CreatedAt
	userResponse.UpdatedAt = user.UpdatedAt
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

	lomToken, err := middleware.GenerateLOMKeys(userResponse)
	println("lom-token: " + lomToken)
	if err != nil {
		errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}
	modelLOM := &models.UserResponse{}
	err = middleware.ClaimsLOM(lomToken, modelLOM)
	if err != nil {
		errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}
	fmt.Printf("\ndecrypt-token: %+v", modelLOM)

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
	userResponse.PermissionGroup = user.PermissionGroup
	userResponse.Username = user.Username
	userResponse.CreatedAt = user.CreatedAt
	userResponse.UpdatedAt = user.UpdatedAt
	userResponse.ChatRoomsUuid = make([]string, 0)

	lomToken, err := middleware.GenerateLOMKeys(userResponse)
	if err != nil {
		errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	c.SetCookie(middleware.LOMCookieAuthPrefix, lomToken, int(middleware.LOMExpireTime), "/", "localhost", false, true)
	c.JSON(http.StatusOK, userResponse)
}

// Require: AuthenticationLOM
func (h *userHandler) ChangePassword(c *gin.Context) {
	user, exist := c.Get(middleware.LOMUserPrefix)
	if !exist {
		errResponse := utils.NewErrorResponse(http.StatusInternalServerError, "user not found")
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	userLOM := user.(models.UserAuthenticationLOM)

	var reqBody models.UserChangePasswordRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		errResponse := utils.NewErrorResponse(http.StatusBadRequest, "invalid request body")
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	if reqBody.OldPassword == reqBody.NewPassword {
		errResponse := utils.NewErrorResponse(http.StatusBadRequest, "don't use the same password")
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	if err := h.userUsecase.ChangePassword(c, userLOM.Uuid, reqBody.OldPassword, reqBody.NewPassword); err != nil {
		errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	successResponse := utils.NewSuccessResponse(http.StatusOK, "password changed", nil)
	c.JSON(http.StatusOK, successResponse)
}

func (h *userHandler) DeleteUser(c *gin.Context) {
	value, exist := c.Get(middleware.LOMUserPrefix)
	if !exist {
		errResponse := utils.NewErrorResponse(http.StatusInternalServerError, "user not found")
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	user, ok := value.(models.UserAuthenticationLOM)
	if !ok {
		errResponse := utils.NewErrorResponse(http.StatusInternalServerError, "user not found")
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	var reqBody models.UserDeleteByUuidRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		errResponse := utils.NewErrorResponse(http.StatusBadRequest, "invalid request body")
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	if user.Uuid == reqBody.Uuid {
		errResponse := utils.NewErrorResponse(http.StatusBadRequest, "can't delete yourself")
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	deletedUser, err := h.userUsecase.DeleteUser(c, reqBody.Uuid)
	if err != nil {
		errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	successResponse := utils.NewSuccessResponse(http.StatusOK, "user deleted", deletedUser)
	c.JSON(http.StatusOK, successResponse)
}
