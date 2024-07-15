package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func (repo *chatRoomRepository) DeleteChatroomByUuid(ctx context.Context, uuid string) error {
	filter := bson.M{
		"uuid": uuid,
	}

	_, err := repo.chatroomCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("can't Delete chatroom by uuid %s", uuid)
	}
	return nil
}
