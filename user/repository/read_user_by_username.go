package repository

import (
	"context"
	"fmt"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (repo *userRepository) ReadUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := models.User{}
	filter := bson.M{
		"username": username,
	}

	foundUser := repo.userCollection.FindOne(ctx, filter)
	err := foundUser.Decode(&user)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return nil, fmt.Errorf("can't find user by username %s", username)
	}
	return &user, nil
}
