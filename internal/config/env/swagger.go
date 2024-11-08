package env

import (
	"errors"
	"net"
	"os"
)

const (
	swaggerHost = "SWAGGER_HOST"
	swaggerPort = "SWAGGER_PORT"
)

type swaggerConfig struct {
	host string
	port string
}

func NewSwaggerConfig() (*swaggerConfig, error) {
	host := os.Getenv(swaggerHost)
	if len(host) == 0 {
		return nil, errors.New("swagger host not found")
	}

	port := os.Getenv(swaggerPort)
	if len(port) == 0 {
		return nil, errors.New("swagger port not found")
	}

	return &swaggerConfig{
		host: host,
		port: port,
	}, nil
}

func(c *swaggerConfig) Address() string {
	return net.JoinHostPort(c.host, c.port)
}
