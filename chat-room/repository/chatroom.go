package repository

import (
	"context"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type chatroomReposotiry interface {
	CreateChatroom(ctx context.Context, user models.User) (*models.User, error)
	ReadChatroomByUuid(ctx context.Context, uuid string) (*models.User, error)
	UpdateChatroomByUuid(ctx context.Context, uuid string, updateUser models.User) (*models.User, error)
	DeleteChatroomByUuid(ctx context.Context, uuid string) error
}

type mongoChatroomRepository struct {
	chatroomCollection *mongo.Collection
}
