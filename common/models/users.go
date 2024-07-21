package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	DefaultPermissionGroup = "user"
)

type User struct {
	Uuid            primitive.Binary   `json:"uuid" bson:"uuid"`
	ChatRoomsUuid   []primitive.Binary `json:"chat_rooms_uuid" bson:"chat_rooms_uuid"`
	Username        string             `json:"username" bson:"username"`
	Password        string             `json:"password" bson:"password"`
	PermissionGroup string             `json:"permission_group" bson:"permission_group"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt       *time.Time         `json:"deleted_at" bson:"deleted_at"`
}

type UserAuthenticationLOM struct {
	Uuid            string     `json:"uuid"`
	ChatRoomsUuid   []string   `json:"chat_rooms_uuid"`
	Username        string     `json:"username"`
	Password        string     `json:"password"`
	PermissionGroup string     `json:"permission_group"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
}

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	Uuid            string     `json:"uuid"`
	ChatRoomsUuid   []string   `json:"chat_rooms_uuid"`
	Username        string     `json:"username"`
	PermissionGroup string     `json:"permission_group"`
	CreatedAt       time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" bson:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at" bson:"deleted_at"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type UserDeleteByUuidRequest struct {
	Uuid string `json:"uuid" binding:"required"`
}
