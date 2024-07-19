package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatRoom struct {
	Uuid         primitive.Binary   `json:"uuid" bson:"uuid"`
	OwnerUuid    primitive.Binary   `json:"owner_uuid" bson:"owner_uuid"`
	UsersUuid    []primitive.Binary `json:"users_uuid" bson:"users_uuid"`
	MessagesUuid []primitive.Binary `json:"messages_uuid" bson:"messages_uuid"`
	Name         string             `json:"name" bson:"name"`
	JoinCode     string             `json:"join_code" bson:"join_code"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt    time.Time          `json:"deleted_at" bson:"deleted_at"`
}

type CreateChatRoomRequest struct {
	Name string `json:"name" binding:"required"`
}

type JoinChatRoomRequest struct {
	JoinCode string `json:"join_code" bson:"join_code"`
}
