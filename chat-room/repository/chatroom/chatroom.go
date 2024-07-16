package repository

import (
	"context"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatRoomReposotiry interface {
	CreateChatRoom(ctx context.Context, chatroom models.ChatRoom) (*models.ChatRoom, error)
	ReadChatRoomByUuid(ctx context.Context, uuid string) (*models.ChatRoom, error)
	ReadChatRoomByRoomName(ctx context.Context, roomname string) (*models.ChatRoom, error)
	ReadChatRoomByJoinCode(ctx context.Context, joinCode string) (*models.ChatRoom, error)
	UpdateChatRoomByUuid(ctx context.Context, uuid string, updateChatRoom models.ChatRoom) (*models.ChatRoom, error)
	DeleteChatRoomByUuid(ctx context.Context, uuid string) error
}

type chatRoomRepository struct {
	collection *mongo.Collection
}

func NewMongoChatroomRepository(collection *mongo.Collection) ChatRoomReposotiry {
	return &chatRoomRepository{
		collection: collection,
	}
}
