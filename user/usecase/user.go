package usecase

import (
	"context"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"github.com/Sincerelyzl/larb-on-me/user/repository"
)

type UserUsecase interface {
	CreateNewUser(ctx context.Context, username string, password string) (*models.User, error)
	Login(ctx context.Context, username string, password string) (*models.User, error)
	ChangePassword(ctx context.Context, uuid string, oldPassword string, newPassword string) error
	DeleteUser(ctx context.Context, uuid string) (*models.User, error)
	GetUsers(ctx context.Context, page int64) ([]*models.User, error)
	AddChatRoomUUID(ctx context.Context, userUuid string, chatRoomUUID string) error
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}
