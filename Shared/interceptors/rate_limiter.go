package interceptors

import (
	"context"

	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RateLimiter(logger *zap.Logger, limiter *rate.Limiter) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		if !limiter.Allow() {
			logger.Warn("the number of requests has exceeded the limit. The request has been rejected.")
			return nil, status.Errorf(codes.ResourceExhausted, "the number of requests has exceeded the limit")
		}
		return handler(ctx, req)
	}
}
