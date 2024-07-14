package repository

import (
	"context"
	"fmt"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"github.com/Sincerelyzl/larb-on-me/common/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userReposotiry interface {
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
	ReadUserByUuid(ctx context.Context, uuid string) (*models.User, error)
	UpdateUserByUuid(ctx context.Context, uuid string, updateUser models.User) (*models.User, error)
	DeleteUserByUuid(ctx context.Context, uuid string) error
}

type mongoUserRepository struct {
	userCollection *mongo.Collection
}

func NewMongoUserRepository(collection *mongo.Collection) *mongoUserRepository {
	return &mongoUserRepository{
		userCollection: collection,
	}
}

func (repo *mongoUserRepository) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	newUuidV7, err := utils.NewUuidV7()
	if err != nil {
		return nil, fmt.Errorf("can't generate new UUIDv7: %s", err)
	}
	// start default initialize
	user.Uuid = &newUuidV7
	chatRoomUuid := []primitive.Binary{}
	// end default initialize

	user.ChatRoomUuid = &chatRoomUuid
	result, err := repo.userCollection.InsertOne(ctx, user)
	if err != nil || result.InsertedID == nil {
		return nil, fmt.Errorf("can't insert user: %s", user.Username)
	}
	return &user, nil
}

func (repo *mongoUserRepository) ReadUserByUuid(ctx context.Context, uuid string) (*models.User, error) {
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

// TODO implement UpdateUser
func (repo *mongoUserRepository) UpdateUserByUuid(ctx context.Context, uuid string, updateUser models.User) (*models.User, error) {
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

// TODO implement DeleteUser
func (repo *mongoUserRepository) DeleteUserByUuid(ctx context.Context, uuid string) error {
	filter := bson.M{
		"uuid": uuid,
	}

	_, err := repo.userCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("can't Delete user by uuid %s", uuid)
	}
	return nil
}
