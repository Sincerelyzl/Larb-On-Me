package usecase

import (
	"context"
	"fmt"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"github.com/Sincerelyzl/larb-on-me/common/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (uc *userUsecase) CreateNewUser(ctx context.Context, username string, password string) (*models.User, error) {

	// Generate new uuid version 7
	uuidV7, err := utils.NewUuidV7()
	if err != nil {
		return nil, err
	}

	// Create empty chatroom slice store uuid
	chatRoomsUuid := []primitive.Binary{}

	// create user object
	var user models.User

	// initial default value for user
	user.Uuid = uuidV7
	user.Username = username
	user.Password = password
	user.ChatRoomsUuid = chatRoomsUuid

	// find have user exist in database by username by count it
	foundCount, err := uc.userRepo.CountUserByUsername(ctx, user.Username)
	if err != nil {
		return nil, err
	}

	// user with username already exist
	if foundCount > 0 {
		return nil, fmt.Errorf("user with username %s already exist", user.Username)
	}

	// insert user to database
	newUser, err := uc.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
