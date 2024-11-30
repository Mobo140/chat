package env

import (
	"errors"
	"net"
	"os"

	"github.com/Mobo140/microservices/chat/internal/config"
)

var _ config.GRPCConfig = (*grpcAuthConfig)(nil)

const (
	grpcAuthHostEnvName = "AUTH_GRPC_HOST"
	grpcAuthPortEnvName = "AUTH_GRPC_PORT"
)

type grpcAuthConfig struct {
	host string
	port string
}

func NewAuthGRPCConfig() (*grpcAuthConfig, error) { //nolint:revive // it's ok
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("grpc auth host not found")
	}

	port := os.Getenv(grpcAuthPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("grpc auth port not found")
	}

	return &grpcAuthConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *grpcAuthConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
