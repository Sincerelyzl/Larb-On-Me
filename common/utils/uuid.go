package utils

import (
	"errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewUuidV7() (primitive.Binary, error) {
	uuidV7, err := uuid.NewV7()
	if err != nil {
		return primitive.Binary{}, err
	}

	return primitive.Binary{
		Subtype: 0x04,
		Data:    uuidV7[:],
	}, nil
}

func UuidV7ToString(binaryUUID primitive.Binary) (string, error) {
	if binaryUUID.Subtype != 0x04 || len(binaryUUID.Data) != 16 {
		return "", errors.New("invalid UUIDv7")
	}

	return uuid.UUID(binaryUUID.Data).String(), nil
}

func UuidV7FromString(uuidStr string) (primitive.Binary, error) {
	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return primitive.Binary{}, err
	}

	return primitive.Binary{
		Subtype: 0x04,
		Data:    parsedUUID[:],
	}, nil
}
