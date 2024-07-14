package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func (repo *userRepository) DeleteUserByUuid(ctx context.Context, uuid string) error {
	filter := bson.M{
		"uuid": uuid,
	}

	_, err := repo.userCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("can't Delete user by uuid %s", uuid)
	}
	return nil
}