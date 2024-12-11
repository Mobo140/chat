package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TimeoutUnaryServerInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		ch := make(chan struct{})
		var result interface{}
		var handlerErr error

		go func() {
			result, handlerErr = handler(ctx, req)
			close(ch)
		}()

		select {
		case <-ctx.Done():
			return nil, status.Errorf(codes.DeadlineExceeded, "request timeout: %v", ctx.Err())
		case <-ch:
			return result, handlerErr
		}

	}
}