package main

import (
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"

	"socialnet/pkg/config"
	"socialnet/pkg/interceptor"
	"socialnet/pkg/logger"

	pb "socialnet/services/search/gen"
	"socialnet/services/search/internal/handlers"
	"socialnet/services/search/internal/service"
)

func main() {
	// 🔹 Загружаем переменные окружения
	err := godotenv.Load("services/search/cmd/search-service/.env")
	if err != nil {
		log.Println("⚠️ Warning: cannot load .env, using default values")
	}

	// 🔹 Подключаем gRPC клиентов
	clients := &config.GRPCClients{}

	userClient, err := clients.GetUserClient(os.Getenv("USER_SERVICE_ADDR"))
	if err != nil {
		log.Fatalf("❌ error creating user client: %v", err)
	}

	postClient, err := clients.GetPostClient(os.Getenv("POST_SERVICE_ADDR"))
	if err != nil {
		log.Fatalf("❌ error creating post client: %v", err)
	}

	defer clients.CloseAll()

	// 🔹 Инициализация логгера
	logger.Init("SearchService")
	defer logger.Sync()

	// 🔹 Сервис и хендлер
	searchService := service.NewSearchService(userClient, postClient)
	searchHandler := handlers.NewSearchHandler(searchService)

	// 🔹 gRPC сервер
	port := os.Getenv("SEARCH_SERVICE_PORT")
	if port == "" {
		port = "50058"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("❌ failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.ExtractUserInterceptor(),
			interceptor.LoggingInterceptor(),
		),
	)

	pb.RegisterSearchServiceServer(grpcServer, searchHandler)

	log.Printf("🚀 SearchService started on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("❌ failed to serve: %v", err)
	}
}
