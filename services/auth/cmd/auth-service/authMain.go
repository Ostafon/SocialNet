package main

import (
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"socialnet/pkg/interceptor"
	authpb "socialnet/services/auth/gen"
	"socialnet/services/auth/internal/handlers"
	"socialnet/services/auth/internal/model"
	"socialnet/services/auth/internal/repos"
	"socialnet/services/auth/internal/service"
)

func main() {
	dsn := os.Getenv("AUTH_BD")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=authdb port=5432 sslmode=disable"
	}

	// 🔹 Подключаемся к БД
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(" failed to connect to database: %v", err)
	}

	// 🔹 Автомиграции
	if err := db.AutoMigrate(&model.User{}, &model.RefreshToken{}, &model.PasswordReset{}); err != nil {
		log.Fatalf(" failed to migrate database: %v", err)
	}

	// 🔹 Репозиторий, сервис, хендлер
	repo := repos.NewAuthRepo(db)
	authService := service.NewAuthService(repo)
	authHandler := handlers.NewAuthHandler(authService)

	// 🔹 gRPC сервер
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf(" failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.ExtractUserInterceptor()),
	)
	authpb.RegisterAuthServiceServer(grpcServer, authHandler)

	log.Println(" AuthService started on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf(" failed to serve: %v", err)
	}
}
