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
	"socialnet/pkg/interceptor"
	"socialnet/pkg/logger"
	pb "socialnet/services/notification/gen"
	"socialnet/services/notification/internal/handlers"
	"socialnet/services/notification/internal/model"
	"socialnet/services/notification/internal/repos"
	"socialnet/services/notification/internal/service"
)

func main() {
	_ = godotenv.Load("services/notification/cmd/notif-service/.env")

	dsn := os.Getenv("NOTIFICATION_DB")
	port := os.Getenv("NOTIFICATION_SERVICE_PORT")
	if port == "" {
		port = ":50057"
	}

	logger.Init("NotificationService")
	defer logger.Sync()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ failed to connect to DB: %v", err)
	}
	_ = db.AutoMigrate(&model.Notification{})

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("❌ Redis connection failed: %v", err)
	}

	repo := repos.NewNotificationRepo(db)
	svc := service.NewNotificationService(repo, rdb)
	handler := handlers.NewNotificationHandler(svc)

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
	pb.RegisterNotificationServiceServer(grpcServer, handler)

	log.Println("🚀 NotificationService started on", port)
	go StartGrpcWebWrapper(grpcServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("❌ failed to serve: %v", err)
	}
}
