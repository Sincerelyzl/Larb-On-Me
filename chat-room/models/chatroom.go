package models

type CreateChatRoomRequest struct {
	Name string `json:"name" binding:"required"`
}

type JoinChatRoomRequest struct {
	JoinCode string `json:"join_code" binding:"required"`
}
