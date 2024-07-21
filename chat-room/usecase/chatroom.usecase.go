package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Sincerelyzl/larb-on-me/common/constants"
	"github.com/Sincerelyzl/larb-on-me/common/middleware"
	"github.com/Sincerelyzl/larb-on-me/common/models"
	"github.com/Sincerelyzl/larb-on-me/common/utils"
	"github.com/go-resty/resty/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (uc *chatRoomUsecase) CreateChatRoom(ctx context.Context, lomToken string, createByUerUuid string, chatroom models.CreateChatRoomRequest) (*models.ChatRoom, error) {

	// Generate new uuid version 7
	uuidV7, err := utils.NewUuidV7()
	if err != nil {
		return nil, err
	}
	// Convert uuid to string
	uuidV7StringChatRoom, err := utils.UuidV7ToString(uuidV7)
	if err != nil {
		return nil, err
	}

	var joinCode string
	for {
		joinCode, err = utils.GenerateRandomString(6)
		if err != nil {
			continue
		}

		existChatRoom, err := uc.chatRoomRepo.ReadChatRoomByJoinCode(ctx, joinCode)
		if err != nil && err != mongo.ErrNoDocuments {
			continue
		}
		if existChatRoom == nil {
			break
		}
	}

	// prepare chatroom owner uuid
	ownerUuid, err := utils.UuidV7FromString(createByUerUuid)
	if err != nil {
		return nil, err
	}

	// Create new chatroom
	newChatRoom := models.ChatRoom{}

	newChatRoom.Uuid = uuidV7
	newChatRoom.OwnerUuid = ownerUuid
	newChatRoom.UsersUuid = []primitive.Binary{
		ownerUuid,
	}
	newChatRoom.MessagesUuid = []primitive.Binary{}
	newChatRoom.Name = chatroom.Name
	newChatRoom.JoinCode = joinCode
	newChatRoom.CreatedAt = utils.GetNowUTCTime()
	newChatRoom.UpdatedAt = utils.GetNowUTCTime()
	newChatRoom.DeletedAt = nil

	// Save chatroom to database
	createdChatRoom, err := uc.chatRoomRepo.CreateChatRoom(ctx, newChatRoom)
	if err != nil {
		return nil, err
	}

	// @TODO : call user-service to update user model
	// @START
	userServices, err := uc.registry.Discover(ctx, "user-service")
	if err != nil {
		return nil, err
	}
	if len(userServices) == 0 {
		return nil, fmt.Errorf(constants.ErrServiceUnavailable, "user-service")
	}
	userService := userServices[0]
	userServiceClient := resty.New()
	userServiceClient.SetDebug(true)
	userServiceClient.SetRetryCount(3)
	userServiceClient.SetRetryWaitTime(2 * time.Second)
	userServiceClient.SetHeader(middleware.LOMCookieAuthPrefix, lomToken)
	body := models.UserAddChatRoomRequest{
		Uuid: uuidV7StringChatRoom,
	}
	res, err := userServiceClient.R().SetBody(body).Patch(fmt.Sprintf("http://%s/v1/user/add.chatroom.uuid", userService))
	if err != nil {
		return nil, err
	}
	if !res.IsSuccess() {
		return nil, fmt.Errorf("failed to update user model")
	}
	// @END

	return createdChatRoom, nil
}

func (uc *chatRoomUsecase) JoinChatRoomByJoinCode(ctx context.Context, joinUser models.User, joinCode models.JoinChatRoomRequest) (*models.ChatRoom, error) {
	// Find chatroom by join code
	chatRoom, err := uc.chatRoomRepo.ReadChatRoomByJoinCode(ctx, joinCode.JoinCode)
	if err != nil {
		return nil, err
	}
	// Check if user already in chatroom
	for _, memberUserUuid := range chatRoom.UsersUuid {
		if memberUserUuid.Equal(joinUser.Uuid) {
			return nil, errors.New("you already in this chatroom")
		}
	}
	// Add user to chatroom
	chatRoom.UsersUuid = append(chatRoom.UsersUuid, joinUser.Uuid)

	chatRoomUuidString, err := utils.UuidV7ToString(chatRoom.Uuid)
	if err != nil {
		return nil, err
	}

	// Add chatroom uuid to user model
	joinUser.ChatRoomsUuid = append(joinUser.ChatRoomsUuid, chatRoom.Uuid)

	// @TODO : call user-service to update user model

	//save chatroom to database
	chatRoom, err = uc.chatRoomRepo.UpdateChatRoomByUuid(ctx, chatRoomUuidString, *chatRoom)
	if err != nil {
		return nil, err
	}

	return chatRoom, nil
}

func (uc *chatRoomUsecase) LeaveChatRoom(ctx context.Context, leaveUser models.User, chatRoomUuid string) (*models.ChatRoom, error) {
	// Find chatroom by uuid
	leaveingChatRoom, err := uc.chatRoomRepo.ReadChatRoomByUuid(ctx, chatRoomUuid)
	if err != nil {
		return nil, err
	}

	// Check if user in chatroom
	for i, memberUserUuid := range leaveingChatRoom.UsersUuid {
		if memberUserUuid.Equal(leaveUser.Uuid) {
			leaveingChatRoom.UsersUuid = append(leaveingChatRoom.UsersUuid[:i], leaveingChatRoom.UsersUuid[i+1:]...)
			break
		}
	}
	// Delete chatroom uuid from user model
	for i, chatRoomUuid := range leaveUser.ChatRoomsUuid {
		if chatRoomUuid.Equal(leaveingChatRoom.Uuid) {
			leaveUser.ChatRoomsUuid = append(leaveUser.ChatRoomsUuid[:i], leaveUser.ChatRoomsUuid[i+1:]...)
			break
		}
	}

	// @TODO : call user-service to update user model

	//save chatroom to database
	leaveingChatRoom, err = uc.chatRoomRepo.UpdateChatRoomByUuid(ctx, chatRoomUuid, *leaveingChatRoom)
	if err != nil {
		return nil, err
	}

	return leaveingChatRoom, nil
}

func (uc *chatRoomUsecase) DeleteChatRoom(ctx context.Context, chatRoomOwner models.User, chatRoomUuid string) error {
	// Find chatroom by uuid
	deleteChatRoom, err := uc.chatRoomRepo.ReadChatRoomByUuid(ctx, chatRoomUuid)
	if err != nil {
		return err
	}

	// Check if user is owner of chatroom
	if !deleteChatRoom.OwnerUuid.Equal(chatRoomOwner.Uuid) {
		return errors.New("you are not the owner of this chatroom")
	}

	// delete chatroom from database
	deleteAt := utils.GetNowUTCTime()
	deleteChatRoom.DeletedAt = &deleteAt

	// update chatroom to database
	_, err = uc.chatRoomRepo.UpdateChatRoomByUuid(ctx, chatRoomUuid, *deleteChatRoom)
	if err != nil {
		return err
	}

	return nil
}

func (uc *chatRoomUsecase) GetChatRoomsByUser(ctx context.Context, user models.User) ([]models.ChatRoom, error) {

	chatRooms := []models.ChatRoom{}
	// Find chatroom by uuid
	for _, chatRoomUuid := range user.ChatRoomsUuid {
		// Convert uuid to string
		chatRoomUuidString, err := utils.UuidV7ToString(chatRoomUuid)
		if err != nil {
			return nil, err
		}
		// Find chatroom by uuid
		chatRoom, err := uc.chatRoomRepo.ReadChatRoomByUuid(ctx, chatRoomUuidString)
		if err != nil {
			return nil, err
		}

		chatRooms = append(chatRooms, *chatRoom)
	}

	return chatRooms, nil
}
