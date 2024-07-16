package repository

import (
	"context"
	"fmt"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateChatRoom
func (repo *chatRoomRepository) CreateChatRoom(ctx context.Context, chatroom models.ChatRoom) (*models.ChatRoom, error) {
	result, err := repo.collection.InsertOne(ctx, chatroom)
	if err != nil || result.InsertedID == nil {
		return nil, fmt.Errorf("can't create Chatroom : %s", chatroom.Name)
	}
	return &chatroom, nil
}

// Get ChatRoom by roomname
func (repo *chatRoomRepository) ReadChatRoomByRoomName(ctx context.Context, roomname string) (*models.ChatRoom, error) {
	chatroom := models.ChatRoom{}
	filter := bson.M{
		"name": roomname,
	}

	foundChatroom := repo.collection.FindOne(ctx, filter)
	err := foundChatroom.Decode(&chatroom)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return nil, fmt.Errorf("can't find chatroom by roomname %s", roomname)
	}
	return &chatroom, nil
}

// Get ChatRoom by  uuid
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

// Update ChatRoom by uuid
func (repo *chatRoomRepository) UpdateChatRoomByUuid(ctx context.Context, uuid string, updateChatRoom models.ChatRoom) (*models.ChatRoom, error) {
	var chatroom models.ChatRoom
	filter := bson.M{
		"uuid": uuid,
	}

	foundChatroom := repo.collection.FindOne(ctx, filter)
	err := foundChatroom.Decode(&chatroom)
	if err != nil {
		return nil, fmt.Errorf("can't find this chatroom by uuid %s", uuid)
	}

	chatroom = updateChatRoom

	if _, err := repo.collection.UpdateOne(ctx, filter, chatroom); err != nil {
		return nil, fmt.Errorf("can't update ChatRoom of this uuid %s", uuid)
	}
	return &chatroom, nil
}

// Delete ChatRoom by uuid
func (repo *chatRoomRepository) DeleteChatRoomByUuid(ctx context.Context, uuid string) error {
	filter := bson.M{
		"uuid": uuid,
	}

	_, err := repo.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("can't Delete chatroom by uuid %s", uuid)
	}
	return nil
}
