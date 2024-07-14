package repository

import (
	"context"
	"fmt"

	"github.com/Sincerelyzl/larb-on-me/common/models"
)

func (repo *userRepository) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	result, err := repo.userCollection.InsertOne(ctx, user)
	if err != nil || result.InsertedID == nil {
		return nil, fmt.Errorf("can't insert user: %s", user.Username)
	}
	return &user, nil
}
