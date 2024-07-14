package main

import (
	"context"
	"fmt"

	"github.com/Sincerelyzl/larb-on-me/common/database"
	"github.com/Sincerelyzl/larb-on-me/common/models"
	"github.com/Sincerelyzl/larb-on-me/user/repository"
)

func main() {

	ctx := context.Background()
	// Connect to Database Server
	client, err := database.NewConnection(ctx, "mongodb://localhost:27018/")
	if err != nil {
		panic(err)
	}

	userCollection := client.Database("user_service").Collection("user")
	userRepo := repository.NewMongoUserRepository(userCollection)

	user := models.User{
		Username: "Saksarak",
		Password: "12345",
	}
	createdUser, err := userRepo.CreateUser(ctx, user)
	if err != nil {
		panic(err)
	}
	fmt.Println(createdUser)
}
