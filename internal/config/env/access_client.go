package env

import (
	"errors"
	"net"
	"os"
)

const (
	authClientHost = "ACCESS_CLIENT_HOST"
	authClientPort = "ACCESS_CLIENT_PORT"
)

type accessClientConfig struct {
	host string
	port string
}

func NewAccessClientConfig() (*accessClientConfig, error) {
	host := os.Getenv(authClientHost)
	if len(host) == 0 {
		return nil, errors.New("host for auth client not set")
	}

	port := os.Getenv(authClientPort)
	if len(port) == 0 {
		return nil, errors.New("port for auth client not set")
	}

	return &accessClientConfig{
		host: host,
		port: port,
	}, nil
}

func (c *accessClientConfig) Address() string {
	return net.JoinHostPort(c.host, c.port)
}
