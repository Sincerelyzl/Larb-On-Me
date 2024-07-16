package usecase

import (
	"context"

	repository "github.com/Sincerelyzl/larb-on-me/chat-room/repository/chatroom"
	"github.com/Sincerelyzl/larb-on-me/common/models"
)

type ChatRoomUsecase interface {
	CreateChatRoom(ctx context.Context, chatroom models.CreateChatRoomRequest) (*models.ChatRoom, error)
}

type chatRoomUsecase struct {
	chatRoomRepo repository.ChatRoomReposotiry
}

func NewChatRoomUsecase(chatRoomRepo repository.ChatRoomReposotiry) ChatRoomUsecase {
	return &chatRoomUsecase{
		chatRoomRepo: chatRoomRepo,
	}
}
