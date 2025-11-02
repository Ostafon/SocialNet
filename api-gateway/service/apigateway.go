package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
	"socialnet/pkg/utils"

	midl "socialnet/api-gateway/middlewares"
	authpb "socialnet/services/auth/gen"
	chatpb "socialnet/services/chat/gen"
	commentpb "socialnet/services/comment/gen"
	likepb "socialnet/services/like/gen"
	notificationpb "socialnet/services/notification/gen"
	postpb "socialnet/services/post/gen"
	searchpb "socialnet/services/search/gen"
	userpb "socialnet/services/user/gen"
)

func main() {
	ctx := context.Background()
	mux := runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
			md := metadata.MD{}
			if v := r.Context().Value(utils.UserIDKey); v != nil {
				if userId, ok := v.(string); ok {
					md.Set("user-id", userId)
				}
			}
			return md
		}),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// 🔹 Регистрация всех сервисов
	if err := authpb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts); err != nil {
		log.Fatalf("failed to register auth service: %v", err)
	}
	if err := userpb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost:50052", opts); err != nil {
		log.Fatalf("failed to register user service: %v", err)
	}
	if err := postpb.RegisterPostServiceHandlerFromEndpoint(ctx, mux, "localhost:50053", opts); err != nil {
		log.Fatalf("failed to register post service: %v", err)
	}
	if err := commentpb.RegisterCommentServiceHandlerFromEndpoint(ctx, mux, "localhost:50054", opts); err != nil {
		log.Fatalf("failed to register comment service: %v", err)
	}
	if err := likepb.RegisterLikeServiceHandlerFromEndpoint(ctx, mux, "localhost:50055", opts); err != nil {
		log.Fatalf("failed to register like service: %v", err)
	}
	if err := chatpb.RegisterChatServiceHandlerFromEndpoint(ctx, mux, "localhost:50056", opts); err != nil {
		log.Fatalf("failed to register chat service: %v", err)
	}
	if err := notificationpb.RegisterNotificationServiceHandlerFromEndpoint(ctx, mux, "localhost:50057", opts); err != nil {
		log.Fatalf("failed to register notification service: %v", err)
	}
	if err := searchpb.RegisterSearchServiceHandlerFromEndpoint(ctx, mux, "localhost:50058", opts); err != nil {
		log.Fatalf("failed to register search service: %v", err)
	}

	handler := midl.CorsMiddleware(midl.AuthMiddleware(mux))

	log.Println("API Gateway listening on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
