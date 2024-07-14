package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Uuid          primitive.Binary   `json:"uuid" bson:"uuid"`
	ChatRoomsUuid []primitive.Binary `json:"chat_rooms_uuid" bson:"chat_rooms_uuid"`
	Username      string             `json:"username" bson:"username"`
	Password      string             `json:"password" bson:"password"`
}

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	Uuid          string   `json:"uuid"`
	ChatRoomsUuid []string `json:"chat_rooms_uuid"`
	Username      string   `json:"username"`
}
