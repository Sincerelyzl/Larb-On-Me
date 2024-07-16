package usecase

import (
	"context"
	"time"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"github.com/Sincerelyzl/larb-on-me/common/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (uc *chatRoomUsecase) CreateChatRoom(ctx context.Context, chatroom models.CreateChatRoomRequest) (*models.ChatRoom, error) {

	// Generate new uuid version 7
	uuidV7, err := utils.NewUuidV7()
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
		if err != nil {
			continue
		}
		if existChatRoom == nil {
			break
		}
	}

	// Create new chatroom
	newChatRoom := models.ChatRoom{}

	newChatRoom.Uuid = uuidV7
	newChatRoom.UsersUuid = []primitive.Binary{}
	newChatRoom.MessagesUuid = []primitive.Binary{}
	newChatRoom.Name = chatroom.Name
	newChatRoom.JoinCode = joinCode
	newChatRoom.CreatedAt = time.Now().UTC()
	newChatRoom.UpdatedAt = time.Now().UTC()

	// Save chatroom to database
	createdChatRoom, err := uc.chatRoomRepo.CreateChatRoom(ctx, newChatRoom)
	if err != nil {
		return nil, err
	}

	return createdChatRoom, nil
}
