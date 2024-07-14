package repository

import (
	"context"
	"fmt"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (repo *userRepository) ReadUserByUuid(ctx context.Context, uuid string) (*models.User, error) {
	var user models.User
	filter := bson.M{
		"uuid": uuid,
	}

	foundUser := repo.userCollection.FindOne(ctx, filter)
	err := foundUser.Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("can't find user by uuid %s", uuid)
	}
	return &user, nil
}
