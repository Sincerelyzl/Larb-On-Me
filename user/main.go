package main

import (
	"fmt"

	"github.com/Sincerelyzl/larb-on-me/common/middleware"
)

type test2 struct {
	Message    string `json:"message"`
	Connection string `json:"connection"`
	Content    string `json:"content"`
}
type test struct {
	Username string  `json:"username"`
	Role     string  `json:"role"`
	Name     string  `json:"name"`
	Position string  `json:"position"`
	Date     string  `json:"date"`
	Age      int     `json:"age"`
	Mores    []test2 `json:"mores"`
}

func main() {

	// create context.
	// ctx := context.Background()

	userInput := test{
		Username: "test",
		Role:     "admin",
		Name:     "test",
		Position: "test",
		Date:     "test",
		Age:      20,
		Mores: []test2{
			{
				Message:    "test",
				Connection: "test",
				Content:    "test",
			},
			{
				Message:    "test",
				Connection: "test",
				Content:    "test",
			},
		},
	}
	encrypted, err := middleware.GenerateLOMKeys(userInput)
	if err != nil {
		panic(err)
	}
	fmt.Println("\nencrypyed:" + encrypted)

	claimsUser := test{}
	err = middleware.ClaimsLOM(encrypted, &claimsUser)
	if err != nil {
		panic(err)
	}

	fmt.Println("\ndecrypted:" + fmt.Sprintf("%v", claimsUser))
	// create claims.

	// // Connect to database server.
	// client, err := database.NewConnection(ctx, "mongodb://localhost:27018/")
	// if err != nil {
	// 	panic(err)
	// }

	// // Get all collection need to use.
	// userCollection := client.Database("user_service").Collection("users")

	// // Create all repository to use.
	// userRepo := repository.NewMongoUserRepository(userCollection)

	// // Create all usecase to use.
	// userUseCase := usecase.NewUserUsecase(userRepo)

	// // Create all handler to use.
	// userHandler := handler.NewUserHandler(userUseCase)

	// // Create all http server to use.
	// userHttpServer := httpserver.NewHTTPServer(userHandler)

	// // Run http server.
	// fmt.Println("User service is running on port 3008")
	// if err := userHttpServer.Run("3008"); err != nil {
	// 	panic(err)
	// }

}
