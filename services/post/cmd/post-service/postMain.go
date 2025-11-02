package main

import (
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"socialnet/pkg/interceptor"
	pb "socialnet/services/post/gen"
	"socialnet/services/post/internal/handlers"
	"socialnet/services/post/internal/repos"
	"socialnet/services/post/internal/service"
)

func main() {

	dsn := os.Getenv("POST_BD")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=postdb port=5432 sslmode=disable"
	}

	// 🔹 Подключаемся к БД
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(" failed to connect to database: %v", err)
	}

	// 🔹 Автомиграции
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf(" failed to migrate database: %v", err)
	}

	// 🔹 Репозиторий, сервис, хендлер
	repo := repos.NewPostRepo(db)
	postService := service.NewPostService(repo)
	postHandler := handlers.NewPostHandler(postService)

	// 🔹 gRPC сервер
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf(" failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.ExtractUserInterceptor()),
	)
	pb.RegisterPostServiceServer(grpcServer, postHandler)

	log.Println(" PostService started on :50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf(" failed to serve: %v", err)
	}

}
