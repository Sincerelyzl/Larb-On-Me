package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func (repo *userRepository) CountUserByUsername(ctx context.Context, username string) (int64, error) {
	filter := bson.M{
		"username": username,
	}
	count, err := repo.userCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("can't count user by username %s", username)
	}
	return count, nil
}
