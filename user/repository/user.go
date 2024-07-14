package repository

import (
	"context"
	"fmt"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserReposotiry interface {
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
	ReadUserByUuid(ctx context.Context, uuid string) (*models.User, error)
	UpdateUser()
	DeleteUser()
}

type MongoUserRepository struct {
	*mongo.Collection
}

func (repo *MongoUserRepository) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	result, err := repo.InsertOne(ctx, user)
	if err != nil || result.InsertedID == nil {
		return nil, fmt.Errorf("can't insert user: %s", user.Username)
	}
	return &user, nil
}

func (repo *MongoUserRepository) ReadUserByUuid(ctx context.Context, uuid string) (*models.User, error) {
	var user models.User
	filter := bson.M{
		"uuid": uuid,
	}

	foundUser := repo.FindOne(ctx, filter)
	err := foundUser.Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("can't find user by uuid %s", uuid)
	}
	return &user, nil
}
