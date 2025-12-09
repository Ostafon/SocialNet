package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"socialnet/pkg/logger"
	"time"
)

func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		start := time.Now()
		resp, err := handler(ctx, req)
		elapsed := time.Since(start)

		if err != nil {
			logger.Log.Errorw("❌ gRPC error",
				"method", info.FullMethod,
				"error", err,
				"duration", elapsed)
		} else {
			logger.Log.Infow("✅ gRPC call",
				"method", info.FullMethod,
				"duration", elapsed)
		}
		return resp, err
	}
}
