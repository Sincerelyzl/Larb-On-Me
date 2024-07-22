package handler

import (
	"net/http"
	"strconv"

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
	// println("lom-token: " + lomToken)
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
	// fmt.Printf("\ndecrypt-token: %+v", modelLOM)

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

	userResponse := models.UserResponse{}
	uuidString, err := utils.UuidV7ToString(deletedUser.Uuid)
	if err != nil {
		errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	userResponse.Uuid = uuidString
	userResponse.Username = deletedUser.Username
	userResponse.PermissionGroup = deletedUser.PermissionGroup
	userResponse.CreatedAt = deletedUser.CreatedAt
	userResponse.UpdatedAt = deletedUser.UpdatedAt
	userResponse.DeletedAt = deletedUser.DeletedAt
	userResponse.ChatRoomsUuid = make([]string, 0)

	for _, chatRoomUuid := range deletedUser.ChatRoomsUuid {
		chatRoomUuidString, err := utils.UuidV7ToString(chatRoomUuid)
		if err != nil {
			errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
			c.JSON(errResponse.StatusCode, errResponse)
			return
		}
		userResponse.ChatRoomsUuid = append(userResponse.ChatRoomsUuid, chatRoomUuidString)
	}

	successResponse := utils.NewSuccessResponse(http.StatusOK, "user deleted", userResponse)
	c.JSON(http.StatusOK, successResponse)
}

func (h *userHandler) GetUsers(c *gin.Context) {
	page := int64(1)
	value, exist := c.GetQuery("page")

	if exist {
		var err error
		page, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			errResponse := utils.NewErrorResponse(http.StatusBadRequest, "page query must be a number")
			c.JSON(errResponse.StatusCode, errResponse)
			return
		}
	}

	users, err := h.userUsecase.GetUsers(c, page)

	usersResponse := []models.UserResponse{}
	for _, user := range users {
		userResponse := models.UserResponse{}
		uuidString, err := utils.UuidV7ToString(user.Uuid)
		if err != nil {
			errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
			c.JSON(errResponse.StatusCode, errResponse)
			return
		}
		userResponse.Uuid = uuidString
		userResponse.Username = user.Username
		userResponse.PermissionGroup = user.PermissionGroup
		userResponse.CreatedAt = user.CreatedAt
		userResponse.UpdatedAt = user.UpdatedAt
		userResponse.DeletedAt = user.DeletedAt
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
		usersResponse = append(usersResponse, userResponse)
	}

	if err != nil {
		errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	c.JSON(http.StatusOK, usersResponse)
}

func (h *userHandler) AddChatRoomUUID(c *gin.Context) {

	var reqBody models.UserAddChatRoomRequest

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		errResponse := utils.NewErrorResponse(http.StatusBadRequest, "invalid request body")
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

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

	err := h.userUsecase.AddChatRoomUUID(c, user.Uuid, reqBody.Uuid)
	if err != nil {
		errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	res := utils.NewSuccessResponse(http.StatusOK, "chatroom added", nil)
	c.JSON(res.StatusCode, res)
}
