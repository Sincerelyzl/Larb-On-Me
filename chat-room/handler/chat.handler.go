package handler

import (
	"net/http"

	"github.com/Sincerelyzl/larb-on-me/common/middleware"
	"github.com/Sincerelyzl/larb-on-me/common/models"
	"github.com/Sincerelyzl/larb-on-me/common/utils"
	"github.com/gin-gonic/gin"
)

func (chh *chatRoomHandler) CreateChatRoom(c *gin.Context) {
	value, _ := c.Get(middleware.LOMUserPrefix)
	token := c.GetHeader(middleware.LOMCookieAuthPrefix)
	user := value.(models.UserAuthenticationLOM)

	var reqBody models.CreateChatRoomRequest
	if err := c.ShouldBind(&reqBody); err != nil {
		errResponse := utils.NewErrorResponse(http.StatusBadRequest, "invalid request body")
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	chatRoom, err := chh.chatRoomUsecase.CreateChatRoom(c, token, user.Uuid, reqBody)
	if err != nil {
		errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	chatRoomUuidString, err := utils.UuidV7ToString(chatRoom.Uuid)
	if err != nil {
		errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}
	chatRoomOwnerUuidString, err := utils.UuidV7ToString(chatRoom.OwnerUuid)
	if err != nil {
		errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}
	chatRoomUsersUuid := make([]string, 0)
	for _, userUuid := range chatRoom.UsersUuid {
		userUuidString, err := utils.UuidV7ToString(userUuid)
		if err != nil {
			errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
			c.JSON(errResponse.StatusCode, errResponse)
			return
		}
		chatRoomUsersUuid = append(chatRoomUsersUuid, userUuidString)
	}

	chatRoomMessageUuid := make([]string, 0)

	chatRoomResponse := models.ChatRoomResponse{
		Uuid:         chatRoomUuidString,
		OwnerUuid:    chatRoomOwnerUuidString,
		UsersUuid:    chatRoomUsersUuid,
		MessagesUuid: chatRoomMessageUuid,
		JoinCode:     chatRoom.JoinCode,
		Name:         chatRoom.Name,
		CreatedAt:    chatRoom.CreatedAt,
		UpdatedAt:    chatRoom.UpdatedAt,
		DeletedAt:    chatRoom.DeletedAt,
	}

	c.JSON(http.StatusCreated, chatRoomResponse)
}

func (chh *chatRoomHandler) JoinChatRoomByJoinCode(c *gin.Context) {

}

func (chh *chatRoomHandler) GetChatRoomsByUser(c *gin.Context) {

}

func (chh *chatRoomHandler) LeaveChatRoom(c *gin.Context) {

}

func (chh *chatRoomHandler) DeleteChatRoom(c *gin.Context) {

}
