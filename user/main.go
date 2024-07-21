package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sincerelyzl/larb-on-me/common/database"
	"github.com/Sincerelyzl/larb-on-me/common/middleware"
	"github.com/Sincerelyzl/larb-on-me/common/utils"
	"github.com/Sincerelyzl/larb-on-me/discovery"
	"github.com/Sincerelyzl/larb-on-me/discovery/consul"
	"github.com/Sincerelyzl/larb-on-me/user/handler"
	"github.com/Sincerelyzl/larb-on-me/user/httpserver"
	"github.com/Sincerelyzl/larb-on-me/user/repository"
	"github.com/Sincerelyzl/larb-on-me/user/usecase"

	_ "github.com/joho/godotenv/autoload"
)

var (
	consulAddress  = utils.EnvString("CONSUL_ADDRESS", "localhost:8500")
	serviceHost    = utils.EnvString("SERVICE_HOST", "localhost")
	servicePort    = utils.EnvString("SERVICE_PORT", ":3008")
	serviceName    = utils.EnvString("SERVICE_NAME", "user-service")
	mongoURI       = utils.EnvString("MONGO_URI", "mongodb://localhost:27018/")
	mongoDatabase  = utils.EnvString("MONGO_DATABASE", "user_service")
	collectionUser = utils.EnvString("COLLECTION_USER", "users")
)

func main() {

	// create context.
	ctx := context.Background()

	// create connection database timeout.
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Connect to database server.
	client, err := database.NewConnection(ctxWithTimeout, mongoURI)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	// Check connection.
	if err := client.Ping(ctx, nil); err != nil {
		panic(err)
	}

	// Get all collection need to use.
	userCollection := client.Database(mongoDatabase).Collection(collectionUser)

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
		Addr:    servicePort,
		Handler: userHttpServer.Router,
	}

	// registry consul
	registry, err := consul.NewRegistry(consulAddress, serviceName)
	if err != nil {
		panic(err)
	}
	instanceId := discovery.GenerateInstaceId(serviceName)
	if err = registry.Register(ctx, instanceId, serviceName, serviceHost+servicePort); err != nil {
		panic(err)
	}

	// Health check.
	discovery.CreateThreadHealthCheck(ctx, registry, instanceId, serviceName)

	// Handle signal.
	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	// Run http server.
	go func() {
		middleware.LogGlobal.Log.Info("service is running", "service-name", serviceName, "port", servicePort)
		if err := userHttpServer.Run(server.Addr); err != nil {
			panic(err)
		}
	}()

	// Wait for signal.
	<-done

	// unregister consul
	if err = registry.Unregister(ctx, instanceId, serviceName); err != nil {
		middleware.LogGlobal.Log.Error("unregister service", "error", err)
	}

	// Shutdown http server.
	middleware.LogGlobal.Log.Fatal("service shutting down", "service-name", serviceName, "port", servicePort)
	ctxWithTimeout, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		middleware.LogGlobal.Log.Fatal("server shutdown", "error", err)
	}
	middleware.LogGlobal.Log.Info("shutdown gracefully", "service-name", serviceName, "port", servicePort)
}
