package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"socialnet/pkg/config"
	"socialnet/pkg/interceptor"
	"socialnet/pkg/logger"
	pb "socialnet/services/chat/gen"
	"socialnet/services/chat/internal/handlers"
	"socialnet/services/chat/internal/model"
	"socialnet/services/chat/internal/repos"
	"socialnet/services/chat/internal/service"
)

func main() {
	_ = godotenv.Load("services/chat/cmd/chat-service/.env")

	dsn := os.Getenv("CHAT_DB")
	port := os.Getenv("CHAT_SERVICE_PORT")
	if port == "" {
		port = ":50056"
	}

	logger.Init("ChatService")
	defer logger.Sync()

	// Подключаемся к БД
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ failed to connect to DB: %v", err)
	}
	_ = db.AutoMigrate(&model.Chat{}, &model.Participant{}, &model.Message{})

	// Подключаемся к Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("❌ Redis connection failed: %v", err)
	}
	log.Println("✅ Connected to Redis")

	clients := &config.GRPCClients{}
	defer clients.CloseAll()

	// Инициализация слоёв
	repo := repos.NewChatRepo(db)
	svc := service.NewChatService(repo, rdb, clients)
	handler := handlers.NewChatHandler(svc)

	// Запуск gRPC сервера
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("❌ failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.ExtractUserInterceptor(),
			interceptor.LoggingInterceptor(),
		),
	)

	pb.RegisterChatServiceServer(grpcServer, handler)
	log.Println("🚀 ChatService started on", port)
	go ChatGrpcWebWrapper(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("❌ failed to serve: %v", err)
	}
}
