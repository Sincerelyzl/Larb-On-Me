package repository

import (
	"context"
	"fmt"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"github.com/Sincerelyzl/larb-on-me/common/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (repo *userRepository) UpdateUser(ctx context.Context, filter, updateUser models.User) (*models.User, error) {
	_, err := repo.collection.UpdateOne(ctx, filter, updateUser)
	if err != nil {
		return nil, fmt.Errorf("can't update user by filter %v", filter)
	}
	return &updateUser, nil
}

func (repo *userRepository) UpdateUserByUuid(ctx context.Context, uuid string, updateUser models.User) (*models.User, error) {
	// convert uuid string to uuid version 7
	uuidV7, err := utils.UuidV7FromString(uuid)
	if err != nil {
		return nil, err
	}

	// prepare filter with uuid
	filter := bson.M{
		"uuid": uuidV7,
	}

	// find user by uuid
	foundUser := repo.collection.FindOne(ctx, filter)
	errorFindUser := foundUser.Err()
	if errorFindUser != nil {
		if errorFindUser == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("not found user by uuid for update %s", uuid)
		}
		return nil, errorFindUser
	}

	// update time of modified
	updateUser.UpdatedAt = utils.GetNowUTCTime()

	// prepare filter update user
	filterUpdate := bson.M{
		"$set": updateUser,
	}

	// update user by uuid
	if _, err := repo.collection.UpdateOne(ctx, filter, filterUpdate); err != nil {
		return nil, fmt.Errorf("can't update User of this uuid %s", uuid)
	}
	return &updateUser, nil
}

func (repo *userRepository) ReadUserByUuid(ctx context.Context, uuid string) (*models.User, error) {
	// prepare user model
	var user models.User

	// convert uuid string to uuid version 7
	uuidV7, err := utils.UuidV7FromString(uuid)
	if err != nil {
		return nil, err
	}

	// prepare filter with uuid
	filter := bson.M{
		"uuid": uuidV7,
	}

	// find user by uuid
	foundUser := repo.collection.FindOne(ctx, filter)
	err = foundUser.Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("can't find user by uuid %s", uuid)
	}

	// return user model
	return &user, nil
}

func (repo *userRepository) ReadUserByUsername(ctx context.Context, username string) (*models.User, error) {
	// prepare user model
	user := models.User{}

	// prepare filter with username
	filter := bson.M{
		"username": username,
	}

	// find user by username
	foundUser := repo.collection.FindOne(ctx, filter)
	err := foundUser.Decode(&user)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return nil, fmt.Errorf("can't find user by username %s", username)
	}

	// return user model
	return &user, nil
}

func (repo *userRepository) DeleteUserByUuid(ctx context.Context, uuid string) error {
	// convert uuid string to uuid version 7
	uuidV7, err := utils.UuidV7FromString(uuid)
	if err != nil {
		return err
	}

	// prepare filter with uuid
	filter := bson.M{
		"uuid": uuidV7,
	}

	// prepare filter update user
	filterUpdate := bson.M{
		"$set": bson.M{
			"$cond": bson.M{
				"if": bson.M{
					"$eq": bson.A{
						"$deleted_at", nil,
					},
				},
				"then": bson.M{
					"deleted_at": utils.GetNowUTCTime(),
				},
			},
		},
	}

	// delete user by uuid
	_, err = repo.collection.UpdateOne(ctx, filter, filterUpdate)
	if err != nil {
		return fmt.Errorf("can't delete user by uuid %s", uuid)
	}

	// return nil if success
	return nil
}

func (repo *userRepository) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	// insert user to database
	result, err := repo.collection.InsertOne(ctx, user)
	if err != nil || result.InsertedID == nil {
		return nil, fmt.Errorf("can't insert user: %s", user.Username)
	}

	// return user model
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
