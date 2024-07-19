package handler

import (
	"net/http"

	"github.com/Sincerelyzl/larb-on-me/chat-room/models"
	"github.com/Sincerelyzl/larb-on-me/common/utils"
	"github.com/gin-gonic/gin"
)

func (chh *chatRoomHandler) CreateChatRoom(c *gin.Context) {

	var reqBody models.CreateChatRoomRequest
	if err := c.ShouldBind(&reqBody); err != nil {
		errResponse := utils.NewErrorResponse(http.StatusBadRequest, "invalid request body")
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	chatRoom, err := chh.chatRoomUsecase.CreateChatRoom(c, reqBody)
	if err != nil {
		errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	c.JSON(http.StatusCreated, chatRoom)

}
func (chh *chatRoomHandler) JoinChatRoomByJoinCode(c *gin.Context) {

}
func (chh *chatRoomHandler) GetChatRoomsByUser(c *gin.Context) {

}
func (chh *chatRoomHandler) LeaveChatRoom(c *gin.Context) {

}
func (chh *chatRoomHandler) DeleteChatRoom(c *gin.Context) {

}
