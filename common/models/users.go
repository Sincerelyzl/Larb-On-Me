package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Uuid         primitive.Binary   `json:"uuid" bson:"uuid"`
	Username     string             `json:"username" bson:"username"`
	Password     string             `json:"password" bson:"password"`
	ChatRoomUuid []primitive.Binary `json:"chatRoomUuid" bson:"chat_room_uuid"`
}
