package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ExtractUserInterceptor — вытаскивает user-id и username из metadata (приходит от API Gateway)
func ExtractUserInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if ids := md.Get("user-id"); len(ids) > 0 {
				ctx = context.WithValue(ctx, "user_id", ids[0])
			}
			if names := md.Get("username"); len(names) > 0 {
				ctx = context.WithValue(ctx, "username", names[0])
			}
		}

		return handler(ctx, req)
	}
}
