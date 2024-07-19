package main

import (
	"context"
	"fmt"
	"time"

	repository "github.com/Sincerelyzl/larb-on-me/chat-room/repository/chatroom"
	"github.com/gin-gonic/gin"

	"github.com/Sincerelyzl/larb-on-me/chat-room/handler"
	"github.com/Sincerelyzl/larb-on-me/chat-room/httpserver"
	"github.com/Sincerelyzl/larb-on-me/chat-room/usecase"
	"github.com/Sincerelyzl/larb-on-me/common/database"
)

func main() {

	// create context.
	ctx := context.Background()

	// Connect to database server.
	ctxtimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	client, err := database.NewConnection(ctxtimeout, "mongodb://localhost:27019")
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	// Get all collection need to use.
	chatRoomCollection := client.Database("chatroom_service").Collection("chatrooms")

	// Create all repository to use.
	chatRoomRepo := repository.NewMongoChatroomRepository(chatRoomCollection)

	// Create all usecase to use.
	chatRoomUseCase := usecase.NewChatRoomUsecase(chatRoomRepo)

	// Create all handler to use.
	chatRoomHandler := handler.NewChatRoomHandler(chatRoomUseCase)

	// Create all http server to use.
	chatRoomHttpServer := httpserver.NewHTTPServer(chatRoomHandler)

	gin.SetMode(gin.ReleaseMode)

	// Run http server.
	fmt.Println("User service is running on port 3009")
	if err := chatRoomHttpServer.Run("3009"); err != nil {
		panic(err)
	}
}
