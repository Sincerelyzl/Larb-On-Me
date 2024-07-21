package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sincerelyzl/larb-on-me/common/database"
	"github.com/Sincerelyzl/larb-on-me/common/middleware"
	"github.com/Sincerelyzl/larb-on-me/discovery"
	"github.com/Sincerelyzl/larb-on-me/discovery/consul"
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

	// registry consul
	registry, err := consul.NewRegistry("localhost:8500", "user-service")
	if err != nil {
		panic(err)
	}
	instanceId := discovery.GenerateInstaceId("user-service")
	if err = registry.Register(ctx, instanceId, "user-service", "localhost:3008"); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceId, "user-service"); err != nil {
				middleware.LogGlobal.Log.Error("health check", "error", err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// Handle signal.
	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	// Run http server.
	go func() {
		middleware.LogGlobal.Log.Info("user-service is running", "port", "3008")
		if err := userHttpServer.Run(server.Addr); err != nil {
			panic(err)
		}
	}()

	// Wait for signal.
	<-done

	// unregister consul
	if err = registry.Unregister(ctx, instanceId, "user-service"); err != nil {
		middleware.LogGlobal.Log.Error("unregister service", "error", err)
	}

	// Shutdown http server.
	fmt.Println("user-service is shutting down")
	ctxWithTimeout, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		middleware.LogGlobal.Log.Fatal("server shutdown", "error", err)
	}
	middleware.LogGlobal.Log.Info("shutdown gracefully", "service", "user-service")
}
