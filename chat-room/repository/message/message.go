package repository

import (
	"context"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageReposotiry interface {
	CreateMessage(ctx context.Context, message models.Message) (*models.Message, error)
	ReadMessageByUuid(ctx context.Context, uuid string) (*models.Message, error)
	UpdateMessageByUuid(ctx context.Context, uuid string, updateMessage models.Message) (*models.Message, error)
	DeleteMessageByUuid(ctx context.Context, uuid string) error
}

type messageRepository struct {
	collection *mongo.Collection
}
