package tracing

import (
	"github.com/Mobo140/microservices/chat/internal/model"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

// вынести в платформенную библиотеку.
func Init(logger *zap.Logger, serviceName string, jaegerAddress string) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  string(model.ConstType),
			Param: model.ConstSendAllTracers,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: jaegerAddress,
		},
	}

	_, err := cfg.InitGlobalTracer(serviceName)
	if err != nil {
		logger.Fatal("failed to init tracing", zap.Error(err))
	}
}
