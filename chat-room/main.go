package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	repository "github.com/Sincerelyzl/larb-on-me/chat-room/repository/chatroom"
	"github.com/Sincerelyzl/larb-on-me/discovery"
	"github.com/Sincerelyzl/larb-on-me/discovery/consul"
	"github.com/gin-gonic/gin"

	"github.com/Sincerelyzl/larb-on-me/chat-room/handler"
	"github.com/Sincerelyzl/larb-on-me/chat-room/httpserver"
	"github.com/Sincerelyzl/larb-on-me/chat-room/usecase"
	"github.com/Sincerelyzl/larb-on-me/common/database"
	"github.com/Sincerelyzl/larb-on-me/common/middleware"
	"github.com/Sincerelyzl/larb-on-me/common/utils"
)

var (
	consulAddress          = utils.EnvString("CONSUL_ADDRESS", "localhost:8500")
	serviceHost            = utils.EnvString("SERVICE_HOST", "localhost")
	servicePort            = utils.EnvString("SERVICE_PORT", ":3009")
	serviceName            = utils.EnvString("SERVICE_NAME", "chatroom-service")
	mongoURI               = utils.EnvString("MONGO_URI", "mongodb://localhost:27019/")
	mongoDatabaseName      = utils.EnvString("MONGO_DATABASE_NAME", "chatroom_service")
	collectionChatRoomName = utils.EnvString("COLLECTION_CHATROOM_NAME", "chatrooms")
)

func main() {

	// create context.
	ctx := context.Background()

	// create connection database timeout.
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
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
	chatRoomCollection := client.Database(mongoDatabaseName).Collection(collectionChatRoomName)

	// Create all repository to use.
	chatRoomRepo := repository.NewMongoChatroomRepository(chatRoomCollection)

	// Create all usecase to use.
	chatRoomUseCase := usecase.NewChatRoomUsecase(chatRoomRepo)

	// Create all handler to use.
	chatRoomHandler := handler.NewChatRoomHandler(chatRoomUseCase)

	// Create all http server to use.
	chatRoomHttpServer := httpserver.NewHTTPServer(chatRoomHandler)

	// create http server.
	server := &http.Server{
		Addr:    servicePort,
		Handler: chatRoomHttpServer.Router,
	}

	// registry consul
	registry, err := consul.NewRegistry(consulAddress, serviceName)
	if err != nil {
		panic(err)
	}
	instanceId := discovery.GenerateInstaceId(serviceName)
	if err = registry.Register(ctx, instanceId, serviceHost, servicePort); err != nil {
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
		if err := chatRoomHttpServer.Run(server.Addr); err != nil {
			panic(err)
		}
	}()
	gin.SetMode(gin.ReleaseMode)

	<-done

	// unregister consul
	if err := registry.Unregister(ctx, instanceId, serviceName); err != nil {
		panic(err)
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
