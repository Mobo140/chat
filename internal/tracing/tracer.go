package tracing

import (
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"github.com/Mobo140/microservices/chat/internal/model"
)


func Init(logger *zap.Logger, serviceName string) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: model.ConstSendAllTracers,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: "localhost:6831",
		},
	}

	_, err := cfg.InitGlobalTracer(serviceName)
	if err != nil {
		logger.Fatal("failed to init tracing", zap.Error(err))
	}
}
