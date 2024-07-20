package repository

import (
	"context"
	"time"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
	ReadUserByUuid(ctx context.Context, uuid string) (*models.User, error)
	ReadUserByUsername(ctx context.Context, username string) (*models.User, error)
	ReadUsers(ctx context.Context, pagination *models.Pagination) ([]*models.User, error)
	CountUserByUsername(ctx context.Context, username string) (int64, error)
	UpdateUserByUuid(ctx context.Context, uuid string, updateUser models.User) (*models.User, error)
	UpdateUser(ctx context.Context, filter, updateUser models.User) (*models.User, error)
	DeleteUserByUuid(ctx context.Context, uuid string) (*time.Time, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(collection *mongo.Collection) UserRepository {
	return &userRepository{
		collection: collection,
	}
}
