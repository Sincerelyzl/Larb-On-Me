package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"github.com/Sincerelyzl/larb-on-me/common/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (repo *userRepository) DeleteUserByUuid(ctx context.Context, uuid string) (*time.Time, error) {
	// convert uuid string to uuid version 7
	uuidV7, err := utils.UuidV7FromString(uuid)
	if err != nil {
		return nil, err
	}

	// prepare filter with uuid
	filter := bson.M{
		"uuid": uuidV7,
	}

	// prepare delete time
	deletedAt := utils.GetNowUTCTime()

	// prepare filter update user
	filterUpdate := bson.M{
		"$set": bson.M{
			"deleted_at": deletedAt,
		},
	}

	// delete user by uuid
	_, err = repo.collection.UpdateOne(ctx, filter, filterUpdate)
	if err != nil {
		fmt.Println("error while deleting user by uuid: ", err)
		return nil, fmt.Errorf("can't delete user by uuid %s", uuid)
	}

	// return nil if success
	return &deletedAt, nil
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

func (repo *userRepository) ReadUsers(ctx context.Context, pagination *models.Pagination) ([]*models.User, error) {
	// Prepare a slice to store the users
	users := []*models.User{}

	// Set the options for the query
	findOptions := options.Find()
	findOptions.SetSkip(pagination.Offset)
	findOptions.SetLimit(pagination.Limit)

	// Define filter (empty in this case to get all users)
	filter := bson.M{}

	// Create a context with timeout for the query
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Perform the query
	cursor, err := repo.collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Println("error while fetching users:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate over the cursor and decode each document into a user model
	for cursor.Next(ctx) {
		user := &models.User{}
		if err := cursor.Decode(&user); err != nil {
			log.Println("error decoding user:", err)
			return nil, err
		}
		users = append(users, user)
	}

	// Check for any cursor errors
	if err := cursor.Err(); err != nil {
		log.Println("cursor error:", err)
		return nil, err
	}

	// Return the users
	return users, nil
}
