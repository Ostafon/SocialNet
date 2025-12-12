package main

import (
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"socialnet/pkg/config"
	"socialnet/pkg/interceptor"
	"socialnet/pkg/logger"
	pb "socialnet/services/post/gen"
	"socialnet/services/post/internal/handlers"
	"socialnet/services/post/internal/model"
	"socialnet/services/post/internal/repos"
	"socialnet/services/post/internal/service"
)

func main() {
	err := godotenv.Load("services/post/cmd/post-service/.env")
	dsn := os.Getenv("POST_BD")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=postdb port=5432 sslmode=disable"
	}

	clients := &config.GRPCClients{}
	defer clients.CloseAll()

	logger.Init("PostService")
	defer logger.Sync()

	// 🔹 Подключаемся к БД
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(" failed to connect to database: %v", err)
	}

	// 🔹 Автомиграции
	if err := db.AutoMigrate(&model.Post{}); err != nil {
		log.Fatalf(" failed to migrate database: %v", err)
	}

	// 🔹 Репозиторий, сервис, хендлер
	repo := repos.NewPostRepo(db)
	postService := service.NewPostService(repo, clients)
	postHandler := handlers.NewPostHandler(postService)

	// 🔹 gRPC сервер
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf(" failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.ExtractUserInterceptor(),
			interceptor.LoggingInterceptor(),
		),
	)
	pb.RegisterPostServiceServer(grpcServer, postHandler)

	log.Println(" PostService started on :50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf(" failed to serve: %v", err)
	}

}
