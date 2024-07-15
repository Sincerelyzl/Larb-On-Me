package repository

import (
	"context"
	"fmt"

	"github.com/Sincerelyzl/larb-on-me/common/models"
)

func (repo *chatRoomRepository) CreateChatRoom(ctx context.Context, chatroom models.ChatRoom) (*models.ChatRoom, error) {
	result, err := repo.collection.InsertOne(ctx, chatroom)
	if err != nil || result.InsertedID == nil {
		return nil, fmt.Errorf("can't create Chatroom : %s", chatroom.Name)
	}
	return &chatroom, nil
}
