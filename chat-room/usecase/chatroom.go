package usecase

import (
	"context"

	repository "github.com/Sincerelyzl/larb-on-me/chat-room/repository/chatroom"
	"github.com/Sincerelyzl/larb-on-me/common/models"
	"github.com/Sincerelyzl/larb-on-me/discovery/consul"
)

type ChatRoomUsecase interface {
	CreateChatRoom(ctx context.Context, lomToken string, createByUerUuid string, chatroom models.CreateChatRoomRequest) (*models.ChatRoom, error)
	JoinChatRoomByJoinCode(ctx context.Context, joinUser models.User, joinCode models.JoinChatRoomRequest) (*models.ChatRoom, error)
	GetChatRoomsByUser(ctx context.Context, user models.User) ([]models.ChatRoom, error)
	LeaveChatRoom(ctx context.Context, leaveUser models.User, chatRoomUuid string) (*models.ChatRoom, error)
	DeleteChatRoom(ctx context.Context, chatRoomOwner models.User, chatRoomUuid string) error
}

type chatRoomUsecase struct {
	registry     *consul.Registry
	chatRoomRepo repository.ChatRoomReposotiry
}

func NewChatRoomUsecase(chatRoomRepo repository.ChatRoomReposotiry, registry *consul.Registry) ChatRoomUsecase {
	return &chatRoomUsecase{
		registry:     registry,
		chatRoomRepo: chatRoomRepo,
	}
}
