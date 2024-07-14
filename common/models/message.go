package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	Uuid      primitive.Binary `json:"uuid" bson:"uuid"`
	UserUuid  primitive.Binary `json:"user_uuid" bson:"user_uuid"`
	Text      string           `json:"text" bson:"text"`
	CreatedAt time.Time        `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time        `json:"updated_at" bson:"updated_at"`
	DeletedAt time.Time        `json:"deleted_at" bson:"deleted_at"`
}
