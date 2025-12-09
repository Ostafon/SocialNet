package main

import (
	"log"
	"net"
	"os"
	"socialnet/pkg/interceptor"
	"socialnet/pkg/logger"
	"socialnet/services/like/internal/model"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	pb "socialnet/services/like/gen"
	"socialnet/services/like/internal/handlers"
	"socialnet/services/like/internal/repos"
	"socialnet/services/like/internal/service"
)

func main() {
	dsn := os.Getenv("LIKE_DB")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=likedb port=5432 sslmode=disable"
	}

	logger.Init("LikeService")
	defer logger.Sync()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	if err := db.AutoMigrate(&model.Like{}); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	repo := repos.NewLikeRepo(db)
	service := service.NewLikeService(repo)
	handler := handlers.NewLikeHandler(service)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.ExtractUserInterceptor(),
			interceptor.LoggingInterceptor(),
		),
	)

	pb.RegisterLikeServiceServer(grpcServer, handler)

	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("LikeService started on :50055")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
