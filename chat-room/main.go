package main

import (
	"context"
	"fmt"
	"time"

	repository "github.com/Sincerelyzl/larb-on-me/chat-room/repository/chatroom"
	"github.com/Sincerelyzl/larb-on-me/discovery/consul"
	"github.com/gin-gonic/gin"

	"github.com/Sincerelyzl/larb-on-me/chat-room/handler"
	"github.com/Sincerelyzl/larb-on-me/chat-room/httpserver"
	"github.com/Sincerelyzl/larb-on-me/chat-room/usecase"
	"github.com/Sincerelyzl/larb-on-me/common/database"
)

func main() {

	// create context.
	ctx := context.Background()

	// @ TESTING CONSUL
	registry, err := consul.NewRegistry("localhost:8500", "chat-room-service")
	if err != nil {
		panic(err)
	}
	err = registry.Register(context.Background(), "chat-room-service-1", "chat-room-service", "localhost:3009")
	if err != nil {
		panic(err)
	}
	services, err := registry.Client.Agent().Services()
	if err != nil {
		panic(err)
	}
	for _, service := range services {
		fmt.Println(service.Service)
	}
	// TESTING CONSUL @

	// create connection database timeout.
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Connect to database server.
	client, err := database.NewConnection(ctxWithTimeout, "mongodb://localhost:27018/")
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	// Check connection.
	if err := client.Ping(ctx, nil); err != nil {
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
	fmt.Println("user-service is running on port 3009")
	if err := chatRoomHttpServer.Run("3009"); err != nil {
		panic(err)
	}
}
