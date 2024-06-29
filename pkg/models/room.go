package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Message struct {
		UUID      primitive.Binary
		UserUUID  primitive.Binary
		Message   string
		CreatedAt time.Time
	}
)
