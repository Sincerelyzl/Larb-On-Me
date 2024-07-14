package usecase

import (
	"context"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"github.com/Sincerelyzl/larb-on-me/user/repository"
)

type UserUsecase interface {
	CreateNewUser(ctx context.Context, username string, password string) (*models.User, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}
