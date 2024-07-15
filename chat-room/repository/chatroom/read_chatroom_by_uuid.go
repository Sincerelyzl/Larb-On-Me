package repository

import (
	"context"
	"fmt"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (repo *chatRoomRepository) ReadChatRoomByUuid(ctx context.Context, uuid string) (*models.ChatRoom, error) {
	var chatroom models.ChatRoom
	filter := bson.M{
		"uuid": uuid,
	}

	foundChatroom := repo.collection.FindOne(ctx, filter)
	err := foundChatroom.Decode(&chatroom)
	if err != nil {
		return nil, fmt.Errorf("can't find this chatroom by uuid %s", uuid)
	}
	return &chatroom, nil
}
