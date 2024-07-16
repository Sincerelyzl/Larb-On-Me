package main

import (
	"context"
	"fmt"

	"github.com/Sincerelyzl/larb-on-me/common/database"
	"github.com/Sincerelyzl/larb-on-me/user/handler"
	"github.com/Sincerelyzl/larb-on-me/user/httpserver"
	"github.com/Sincerelyzl/larb-on-me/user/repository"
	"github.com/Sincerelyzl/larb-on-me/user/usecase"
)

func main() {

	// create context.
	ctx := context.Background()

	// Connect to database server.
	client, err := database.NewConnection(ctx, "mongodb://localhost:27018/")
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	// Get all collection need to use.
	userCollection := client.Database("user_service").Collection("users")

	// Create all repository to use.
	userRepo := repository.NewMongoUserRepository(userCollection)

	// Create all usecase to use.
	userUseCase := usecase.NewUserUsecase(userRepo)

	// Create all handler to use.
	userHandler := handler.NewUserHandler(userUseCase)

	// Create all http server to use.
	userHttpServer := httpserver.NewHTTPServer(userHandler)

	// Run http server.
	fmt.Println("User service is running on port 3008")
	if err := userHttpServer.Run("3008"); err != nil {
		panic(err)
	}
}
