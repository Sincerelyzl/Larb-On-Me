package repository

import (
	"context"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatroomReposotiry interface {
	CreateChatroom(ctx context.Context, chatroom models.ChatRoom) (*models.ChatRoom, error)
	ReadChatroomByUuid(ctx context.Context, uuid string) (*models.ChatRoom, error)
	ReadUserByRoomname(ctx context.Context, roomname string) (*models.ChatRoom, error)
	UpdateChatroomByUuid(ctx context.Context, uuid string, updateChatRoom models.ChatRoom) (*models.ChatRoom, error)
	DeleteChatroomByUuid(ctx context.Context, uuid string) error
}

type chatRoomRepository struct {
	chatroomCollection *mongo.Collection
}

func NewMongoChatroomRepository(collection *mongo.Collection) ChatroomReposotiry {
	return &chatRoomRepository{
		chatroomCollection: collection,
	}
}
