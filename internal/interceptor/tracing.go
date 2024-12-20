package interceptor

import (
	"context"

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

	// Логируем входящие метаданные
	incomingMD, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Debug("No incoming metadata found in context")
	} else {
		logger.Debug("Incoming metadata", zap.Any("metadata", incomingMD))
		// Особенно проверяем authorization header
		if auth := incomingMD.Get("authorization"); len(auth) > 0 {
			logger.Debug("Found authorization header", zap.Strings("auth", auth))
		} else {
			logger.Debug("No authorization header found in metadata")
		}
	}

	// Копируем все метаданные в новый контекст
	ctx = metadata.NewOutgoingContext(ctx, incomingMD)

	spanContext, ok := span.Context().(jaeger.SpanContext)
	if ok {
		// Добавляем trace ID к существующим метаданным
		ctx = metadata.AppendToOutgoingContext(ctx, traceIDKey, spanContext.TraceID().String())

		header := metadata.New(map[string]string{traceIDKey: spanContext.TraceID().String()})
		if err := grpc.SendHeader(ctx, header); err != nil {
			logger.Error("Failed to send header", zap.Error(err))
			return nil, err
		}
	}

	// Логируем исходящие метаданные
	if outgoingMD, ok := metadata.FromOutgoingContext(ctx); ok {
		logger.Debug("Outgoing metadata", zap.Any("metadata", outgoingMD))
	}

	res, err := handler(ctx, req)
	if err != nil {
		ext.Error.Set(span, true)
		span.SetTag("err", err.Error())
		logger.Error("Handler error", zap.Error(err))
	}

	return res, err
}
