package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Uuid          primitive.Binary   `json:"uuid" bson:"uuid"`
	Username      string             `json:"username" bson:"username"`
	Password      string             `json:"password" bson:"password"`
	ChatRoomsUuid []primitive.Binary `json:"chatRoomUuid" bson:"chat_room_uuid"`
}

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	Uuid          string   `json:"uuid"`
	Username      string   `json:"username"`
	ChatRoomsUuid []string `json:"chat_room_uuid"`
}
