package main

import (
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"socialnet/pkg/interceptor"
	"socialnet/pkg/logger"
	userpb "socialnet/services/user/gen"
	"socialnet/services/user/internal/handlers"
	"socialnet/services/user/internal/model"
	"socialnet/services/user/internal/repos"
	"socialnet/services/user/internal/service"
)

func main() {
	dsn := os.Getenv("USER_BD")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=userdb port=5432 sslmode=disable"
	}

	logger.Init("UserService")
	defer logger.Sync()

	// 🔹 Подключаемся к БД
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(" failed to connect to database: %v", err)
	}

	// 🔹 Автомиграции
	if err := db.AutoMigrate(&model.Follow{}, &model.User{}); err != nil {
		log.Fatalf(" failed to migrate database: %v", err)
	}

	// 🔹 Репозиторий, сервис, хендлер
	repo := repos.NewUserRepo(db)
	userService := service.NewUserService(repo)
	userHandler := handlers.NewUserHandler(userService)

	// 🔹 gRPC сервер
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf(" failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.ExtractUserInterceptor(),
			interceptor.LoggingInterceptor(),
		),
	)
	userpb.RegisterUserServiceServer(grpcServer, userHandler)

	log.Println(" UserService started on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf(" failed to serve: %v", err)
	}

}
