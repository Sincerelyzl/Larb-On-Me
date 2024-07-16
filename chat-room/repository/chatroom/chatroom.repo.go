package repository

import (
	"context"
	"fmt"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"github.com/Sincerelyzl/larb-on-me/common/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("can't find this chatroom by this name %s", roomname)
		}
		return nil, fmt.Errorf("error occurred while searching for chatroom: %v", err)
	}
	return &chatroom, nil
}

// Get ChatRoom by  uuid
func (repo *chatRoomRepository) ReadChatRoomByUuid(ctx context.Context, uuid string) (*models.ChatRoom, error) {
	var chatroom models.ChatRoom

	uuidHash, err := utils.UuidV7FromString(uuid)
	if err != nil {
		return nil, fmt.Errorf("can't convert uuid string %s to uuid hash", uuid)
	}

	filter := bson.M{
		"uuid": uuidHash,
	}

	foundChatroom := repo.collection.FindOne(ctx, filter)
	err = foundChatroom.Decode(&chatroom)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("can't find this chatroom by this uuid %s", uuid)
		}
		return nil, fmt.Errorf("error occurred while searching for chatroom: %v", err)
	}
	return &chatroom, nil
}

// Read ChatRoom by join code
func (repo *chatRoomRepository) ReadChatRoomByJoinCode(ctx context.Context, joinCode string) (*models.ChatRoom, error) {
	var chatroom models.ChatRoom
	filter := bson.M{
		"join_code": joinCode,
	}

	foundChatroom := repo.collection.FindOne(ctx, filter)
	err := foundChatroom.Decode(&chatroom)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("can't find this chatroom by this joincode %s", joinCode)
		}
		return nil, fmt.Errorf("error occurred while searching for chatroom: %v", err)
	}
	return &chatroom, nil
}

// Update ChatRoom by uuid
func (repo *chatRoomRepository) UpdateChatRoomByUuid(ctx context.Context, uuid string, updateChatRoom models.ChatRoom) (*models.ChatRoom, error) {

	var chatroom models.ChatRoom

	uuidHash, err := utils.UuidV7FromString(uuid)
	if err != nil {
		return nil, fmt.Errorf("can't convert uuid string %s to uuid hash", uuid)
	}

	filter := bson.M{
		"uuid": uuidHash,
	}

	foundChatroom := repo.collection.FindOne(ctx, filter)
	err = foundChatroom.Decode(&chatroom)
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

	uuidHash, err := utils.UuidV7FromString(uuid)
	if err != nil {
		return fmt.Errorf("can't convert uuid string %s to uuid hash", uuid)
	}

	filter := bson.M{
		"uuid": uuidHash,
	}

	_, err = repo.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("can't Delete chatroom by uuid %s", uuid)
	}
	return nil
}
