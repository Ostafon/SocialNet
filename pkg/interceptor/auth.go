package interceptor

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"socialnet/pkg/contextx"
)

// ExtractUserInterceptor — добавляет user_id и username в контекст
func ExtractUserInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			fmt.Println("📥 Incoming metadata:", md)

			if ids := md.Get("user-id"); len(ids) > 0 {
				ctx = context.WithValue(ctx, contextx.UserIDKey, ids[0])
			}
			if ids := md.Get("x-user-id"); len(ids) > 0 {
				ctx = context.WithValue(ctx, contextx.UserIDKey, ids[0])
			}
			if names := md.Get("username"); len(names) > 0 {
				ctx = context.WithValue(ctx, contextx.UsernameKey, names[0])
			}
		} else {
			fmt.Println("⚠️ No metadata found in context")
		}

		return handler(ctx, req)
	}
}
