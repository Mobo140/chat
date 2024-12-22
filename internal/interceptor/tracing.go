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
    // Получаем входящие метаданные
    incomingMD, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        logger.Debug("No incoming metadata found in context")
        incomingMD = metadata.New(nil)
    }

    // Создаем carrier для извлечения контекста трейсинга
    carrier := opentracing.TextMapCarrier{}
    for k, vals := range incomingMD {
        if len(vals) > 0 {
            carrier[k] = vals[0]
        }
    }

    // Пытаемся извлечь родительский контекст
    parentSpanContext, err := opentracing.GlobalTracer().Extract(
        opentracing.TextMap,
        carrier,
    )

    var span opentracing.Span
    if err != nil {
        span = opentracing.StartSpan(info.FullMethod)
    } else {
        span = opentracing.StartSpan(
            info.FullMethod,
            opentracing.ChildOf(parentSpanContext),
        )
    }
    defer span.Finish()

    // ВАЖНО: Сохраняем оригинальные метаданные в контексте
    ctx = metadata.NewIncomingContext(ctx, incomingMD)
    ctx = opentracing.ContextWithSpan(ctx, span)

    // Добавляем trace ID к существующим метаданным
    if spanContext, ok := span.Context().(jaeger.SpanContext); ok {
        traceID := spanContext.TraceID().String()
        
        // Создаем новые метаданные, сохраняя существующие
        mdCopy := metadata.Join(incomingMD, metadata.New(map[string]string{
            traceIDKey: traceID,
        }))
        
        ctx = metadata.NewOutgoingContext(ctx, mdCopy)

        header := metadata.New(map[string]string{traceIDKey: traceID})
        if err := grpc.SendHeader(ctx, header); err != nil {
            logger.Error("Failed to send header", zap.Error(err))
        }
    }

    // Логируем для отладки
    if auth := incomingMD.Get("authorization"); len(auth) > 0 {
        logger.Debug("Authorization header present")
    } else {
        logger.Debug("No authorization header found")
    }

    res, err := handler(ctx, req)
    if err != nil {
        ext.Error.Set(span, true)
        span.SetTag("error", true)
        span.SetTag("error.message", err.Error())
        logger.Error("Handler error", zap.Error(err))
    }

    return res, err
}