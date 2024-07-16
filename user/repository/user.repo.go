package repository

import (
	"context"
	"fmt"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (repo *userRepository) UpdateUserByUuid(ctx context.Context, uuid string, updateUser models.User) (*models.User, error) {
	var user models.User // prepare a user model
	filter := bson.M{
		"uuid": uuid,
	}

	foundUser := repo.collection.FindOne(ctx, filter)
	err := foundUser.Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("can't find user by uuid for update %s", uuid)
	}

	user = updateUser

	if _, err := repo.collection.UpdateOne(ctx, filter, user); err != nil {
		return nil, fmt.Errorf("can't update User of this uuid %s", uuid)
	}
	return &user, nil
}

func (repo *userRepository) ReadUserByUuid(ctx context.Context, uuid string) (*models.User, error) {
	var user models.User
	filter := bson.M{
		"uuid": uuid,
	}

	foundUser := repo.collection.FindOne(ctx, filter)
	err := foundUser.Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("can't find user by uuid %s", uuid)
	}
	return &user, nil
}

func (repo *userRepository) ReadUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := models.User{}
	filter := bson.M{
		"username": username,
	}

	foundUser := repo.collection.FindOne(ctx, filter)
	err := foundUser.Decode(&user)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return nil, fmt.Errorf("can't find user by username %s", username)
	}
	return &user, nil
}

func (repo *userRepository) DeleteUserByUuid(ctx context.Context, uuid string) error {
	filter := bson.M{
		"uuid": uuid,
	}

	_, err := repo.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("can't Delete user by uuid %s", uuid)
	}
	return nil
}

func (repo *userRepository) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	result, err := repo.collection.InsertOne(ctx, user)
	if err != nil || result.InsertedID == nil {
		return nil, fmt.Errorf("can't insert user: %s", user.Username)
	}
	return &user, nil
}

func (repo *userRepository) CountUserByUsername(ctx context.Context, username string) (int64, error) {
	filter := bson.M{
		"username": username,
	}
	count, err := repo.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("can't count user by username %s", username)
	}
	return count, nil
}
