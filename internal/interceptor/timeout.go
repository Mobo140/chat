package interceptor

import (
	"context"
	"sync"
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

		if ctx.Err() != nil {
			return nil, status.Errorf(codes.Canceled, "request canceled: %v", ctx.Err())
		}

		ch := make(chan struct{})
		var result interface{}
		var handlerErr error
		var mu sync.Mutex

		go func() {
			res, err := handler(ctx, req)

			mu.Lock()
			result = res
			handlerErr = err
			mu.Unlock()

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
