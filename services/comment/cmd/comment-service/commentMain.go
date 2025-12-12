package main

import (
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
	"socialnet/pkg/config"
	"socialnet/pkg/interceptor"
	"socialnet/pkg/logger"
	"socialnet/services/comment/internal/model"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	pb "socialnet/services/comment/gen"
	"socialnet/services/comment/internal/handlers"
	"socialnet/services/comment/internal/repos"
	"socialnet/services/comment/internal/service"
)

func main() {
	err := godotenv.Load("services/user/cmd/comment-service/.env")
	dsn := os.Getenv("COMMENT_DB")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=commentdb port=5432 sslmode=disable"
	}
	logger.Init("CommentService")
	defer logger.Sync()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	if err := db.AutoMigrate(&model.Comment{}); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	clients := &config.GRPCClients{}
	defer clients.CloseAll()

	repo := repos.NewCommentRepo(db)
	service := service.NewCommentService(repo, clients)
	handler := handlers.NewCommentHandler(service)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.ExtractUserInterceptor(),
			interceptor.LoggingInterceptor(),
		),
	)

	pb.RegisterCommentServiceServer(grpcServer, handler)

	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("✅ CommentService started on :50054")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
