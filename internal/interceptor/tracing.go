package interceptor

import (
	"context"
	"errors"

	"github.com/Mobo140/platform_common/pkg/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const traceIDKey = "x-trace-id"

func ServerTracingInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, info.FullMethod)
	defer span.Finish()

	spanContext, ok := span.Context().(jaeger.SpanContext)
	if ok {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			err := errors.New("metadata is not provided")
			logger.Error("Failed to get metadata from context: ",
				zap.Any("request", req),
				zap.Error(err),
			)

			return nil, err
		}

		for key, values := range md {
			for _, value := range values {
				ctx = metadata.AppendToOutgoingContext(ctx, key, value)
			}
		}
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(traceIDKey, spanContext.TraceID().String()))

		header := metadata.New(map[string]string{traceIDKey: spanContext.TraceID().String()})

		err := grpc.SendHeader(ctx, header)
		if err != nil {
			return nil, err
		}
	}

	res, err := handler(ctx, req)
	if err != nil {
		ext.Error.Set(span, true)
		span.SetTag("err", err.Error())
	} // else {
	//	Ответ может быть большим, поэтому лучше его не добавлять его в тэг.
	//	span.SetTag("res", res).
	// }

	return res, err
}
