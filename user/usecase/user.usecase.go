package usecase

import (
	"context"
	"fmt"

	"github.com/Sincerelyzl/larb-on-me/common/constants"
	"github.com/Sincerelyzl/larb-on-me/common/models"
	"github.com/Sincerelyzl/larb-on-me/common/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (uc *userUsecase) Login(ctx context.Context, username, password string) (*models.User, error) {
	// find user by username
	user, err := uc.userRepo.ReadUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	// chec if user deleted
	if user.DeletedAt != nil {
		return nil, fmt.Errorf(constants.ErrUserDeleted, user.DeletedAt.Format(constants.TimeLayout))
	}

	// compare password
	match := utils.VerifyPassword(password, user.Password)
	if !match {
		return nil, constants.ErrPasswordMismatch
	}

	return user, nil
}

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

	// encrypt password
	encryptedPassword, errEncryptedPassword := utils.HashPassword(password)
	if errEncryptedPassword != nil {
		return nil, errEncryptedPassword
	}

	user.PermissionGroup = models.DefaultPermissionGroup
	user.Password = encryptedPassword
	user.ChatRoomsUuid = chatRoomsUuid
	user.CreatedAt = utils.GetNowUTCTime()
	user.UpdatedAt = utils.GetNowUTCTime()

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

func (u *userUsecase) ChangePassword(ctx context.Context, uuid string, oldPassword string, newPassword string) error {
	// find user by uuid
	user, err := u.userRepo.ReadUserByUuid(ctx, uuid)
	if err != nil {
		return err
	}

	// encrypt new password
	encryptedNewPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// set new password to user
	user.Password = encryptedNewPassword

	// update user to database
	if _, err := u.userRepo.UpdateUserByUuid(ctx, uuid, *user); err != nil {
		return err
	}

	// return nil if success
	return nil
}

func (u *userUsecase) DeleteUser(ctx context.Context, uuid string) (*models.User, error) {
	// delete user by uuid
	err := u.userRepo.DeleteUserByUuid(ctx, uuid)
	if err != nil {
		return nil, err
	}

	// return nil if success
	return nil, nil
}
