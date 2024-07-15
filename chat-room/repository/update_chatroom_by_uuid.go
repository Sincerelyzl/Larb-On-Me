package repository

import (
	"context"
	"fmt"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (repo *chatRoomRepository) UpdateChatroomByUuid(ctx context.Context, uuid string, updateChatRoom models.ChatRoom) (*models.ChatRoom, error) {
	var chatroom models.ChatRoom
	filter := bson.M{
		"uuid": uuid,
	}

	foundChatroom := repo.chatroomCollection.FindOne(ctx, filter)
	err := foundChatroom.Decode(&chatroom)
	if err != nil {
		return nil, fmt.Errorf("can't find this chatroom by uuid %s", uuid)
	}

	chatroom = updateChatRoom

	if _, err := repo.chatroomCollection.UpdateOne(ctx, filter, chatroom); err != nil {
		return nil, fmt.Errorf("can't update ChatRoom of this uuid %s", uuid)
	}
	return &chatroom, nil
}