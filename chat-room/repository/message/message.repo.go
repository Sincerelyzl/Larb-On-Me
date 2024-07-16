package repository

import (
	"context"
	"fmt"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateMessage
func (repo *messageRepository) CreateMessage(ctx context.Context, message models.Message) (*models.Message, error) {
	result, err := repo.collection.InsertOne(ctx, message)
	if err != nil || result.InsertedID == nil {
		return nil, fmt.Errorf("can't create Message : %s", message.Text)
	}
	return &message, nil
}

// Get Message by uuid
func (repo *messageRepository) ReadMessageByUuid(ctx context.Context, uuid string) (*models.Message, error) {
	var message models.Message
	filter := bson.M{
		"uuid": uuid,
	}

	foundMessage := repo.collection.FindOne(ctx, filter)
	err := foundMessage.Decode(&message)
	if err != nil {
		return nil, fmt.Errorf("can't find this message by uuid %s", uuid)
	}
	return &message, nil
}

// Update Message by uuid
func (repo *messageRepository) UpdateMessageByUuid(ctx context.Context, uuid string, updateMessage models.Message) (*models.Message, error) {
	var message models.Message
	filter := bson.M{
		"uuid": uuid,
	}

	foundMessage := repo.collection.FindOne(ctx, filter)
	err := foundMessage.Decode(&message)
	if err != nil {
		return nil, fmt.Errorf("can't find this message by uuid %s", uuid)
	}

	message = updateMessage

	_, err = repo.collection.UpdateOne(ctx, filter, updateMessage)
	if err != nil {
		return nil, fmt.Errorf("can't update this message by uuid %s", uuid)
	}
	return &updateMessage, nil
}

// Delete Message by uuid
func (repo *messageRepository) DeleteMessageByUuid(ctx context.Context, uuid string) error {
	filter := bson.M{
		"uuid": uuid,
	}

	_, err := repo.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("can't delete this message by uuid %s", uuid)
	}
	return nil
}
