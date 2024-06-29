package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Room struct {
		UUID          primitive.Binary
		Roomowneruuid primitive.Binary
		Members       []primitive.Binary
		Messages      []primitive.Binary
		Joincode      string
		DeletedAt     time.Time
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}
)
