package env

import (
	"errors"
	"os"

	"github.com/Mobo140/microservices/chat/internal/config"
)

var _ config.PGConfig = (*pgConfig)(nil)

const (
	dsnEnvName = "PG_DSN"
)

type pgConfig struct {
	dsn string
}

func NewPGConfig() (*pgConfig, error) { //nolint:revive // it's ok
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}

	return &pgConfig{
		dsn: dsn,
	}, nil
}

func (c *pgConfig) DSN() string {
	return c.dsn
}
