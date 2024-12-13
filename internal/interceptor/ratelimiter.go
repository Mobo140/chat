package interceptor

import (
	"context"

	rateLimiter "github.com/Mobo140/microservices/chat/internal/ratelimiter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type rateLimiterInterceptor struct {
	rateLimiter *rateLimiter.TokenBucketLimiter
}

func NewRateLimiterInterceptor(
	rateLimiter *rateLimiter.TokenBucketLimiter,
) *rateLimiterInterceptor { //nolint:revive // it's ok
	return &rateLimiterInterceptor{rateLimiter: rateLimiter}
}

func (r *rateLimiterInterceptor) Unary(ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if !r.rateLimiter.Allow() {
		return nil, status.Error(codes.ResourceExhausted, "too many requests")
	}

	return handler(ctx, req)
}
