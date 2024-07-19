package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sincerelyzl/larb-on-me/common/database"
	"github.com/Sincerelyzl/larb-on-me/user/handler"
	"github.com/Sincerelyzl/larb-on-me/user/httpserver"
	"github.com/Sincerelyzl/larb-on-me/user/repository"
	"github.com/Sincerelyzl/larb-on-me/user/usecase"
)

func main() {

	// create context.
	ctx := context.Background()

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
	userCollection := client.Database("user_service").Collection("users")

	// Create all repository to use.
	userRepo := repository.NewMongoUserRepository(userCollection)

	// Create all usecase to use.
	userUseCase := usecase.NewUserUsecase(userRepo)

	// Create all handler to use.
	userHandler := handler.NewUserHandler(userUseCase)

	// Create all http server to use.
	userHttpServer := httpserver.NewHTTPServer(userHandler)

	// create http server.
	server := &http.Server{
		Addr:    ":3008",
		Handler: userHttpServer.Router,
	}

	// Handle signal.
	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	// Run http server.
	go func() {
		fmt.Println("user-service is running on port " + server.Addr)
		if err := userHttpServer.Run(server.Addr); err != nil {
			panic(err)
		}
	}()

	// Wait for signal.
	<-done

	// Shutdown http server.
	fmt.Println("user-service is shutting down")
	ctxWithTimeout, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("error while shutting down Server. Initiating force shutdown...")
		log.Fatalln(err)
	}
	log.Println("server exiting")
}
