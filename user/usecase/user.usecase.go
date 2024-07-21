package usecase

import (
	"context"
	"fmt"
	"log"

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

func (uc *userUsecase) GetUsers(ctx context.Context, page int64) ([]*models.User, error) {
	// set limit and offset for pagination
	limit := constants.LimitPagination
	offset := limit * (page - 1)

	// create pagination object
	pagination := &models.Pagination{
		Limit:  limit,
		Offset: offset,
	}

	// find all user
	users, err := uc.userRepo.ReadUsers(ctx, pagination)
	if err != nil {
		log.Println("error get users: ", err)
		return nil, err
	}

	return users, nil
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

	// compare old password
	match := utils.VerifyPassword(oldPassword, user.Password)
	if !match {
		return constants.ErrOldPasswordNotMatch
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
	// find user by uuid
	userShouldDelete, err := u.userRepo.ReadUserByUuid(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if userShouldDelete.DeletedAt != nil {
		return nil, fmt.Errorf(constants.ErrUserDeleted, userShouldDelete.DeletedAt.Format(constants.TimeLayout))
	}

	// delete user by uuid
	deletedAtTime, err := u.userRepo.DeleteUserByUuid(ctx, uuid)
	if err != nil {
		return nil, err
	}

	// set deleted at time to user
	userShouldDelete.DeletedAt = deletedAtTime

	// return nil if success
	return userShouldDelete, nil
}

func (u *userUsecase) AddChatRoomUUID(ctx context.Context, userUuid string, chatRoomUUID string) error {
	// find user by uuid
	user, err := u.userRepo.ReadUserByUuid(ctx, userUuid)
	if err != nil {
		return err
	}

	// convert chatroom uuid string to uuid
	chatRoomUuidV7, err := utils.UuidV7FromString(chatRoomUUID)
	if err != nil {
		return err
	}

	// check if chatroom uuid already exist in user chatroom
	for _, chatRoom := range user.ChatRoomsUuid {
		if chatRoom.Equal(chatRoomUuidV7) {
			return fmt.Errorf(constants.ErrChatRoomAlreadyExistInUserUuid, chatRoomUUID, userUuid)
		}
	}

	// add chatroom uuid to user chatroom
	user.ChatRoomsUuid = append(user.ChatRoomsUuid, chatRoomUuidV7)

	// update user to database
	if _, err := u.userRepo.UpdateUserByUuid(ctx, userUuid, *user); err != nil {
		return err
	}

	// return nil if success
	return nil
}
