package repository

import (
	"context"
	"fmt"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (repo *chatRoomRepository) ReadUserByRoomname(ctx context.Context, roomname string) (*models.ChatRoom, error) {
	chatroom := models.ChatRoom{}
	filter := bson.M{
		"name": roomname,
	}

	foundChatroom := repo.chatroomCollection.FindOne(ctx, filter)
	err := foundChatroom.Decode(&chatroom)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return nil, fmt.Errorf("can't find chatroom by roomname %s", roomname)
	}
	return &chatroom, nil
}
