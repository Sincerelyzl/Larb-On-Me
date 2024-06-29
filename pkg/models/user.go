package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	User struct {
		UUID      primitive.Binary
		Username  string
		Password  string
		Name      string
		CreatedAt time.Time
		UpdatedAt time.Time
	}
)
