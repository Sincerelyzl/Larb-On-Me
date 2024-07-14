package repository

import (
	"context"
	"fmt"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (repo *userRepository) UpdateUserByUuid(ctx context.Context, uuid string, updateUser models.User) (*models.User, error) {
	var user models.User
	filter := bson.M{
		"uuid": uuid,
	}

	foundUser := repo.userCollection.FindOne(ctx, filter)
	err := foundUser.Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("can't find user by uuid for update %s", uuid)
	}

	user = updateUser

	if _, err := repo.userCollection.UpdateOne(ctx, filter, user); err != nil {
		return nil, fmt.Errorf("can't update User of this uuid %s", uuid)
	}
	return &user, nil

}
