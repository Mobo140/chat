package config

import "github.com/joho/godotenv"

type GRPCConfig interface {
	Address() string
}

type HTTPConfig interface {
	Address() string
}

type PGConfig interface {
	DSN() string
}

type SwaggerConfig interface {
	Address() string
}

type AccessClientConfig interface {
	Address() string
}

type JaegerConfig interface {
	Address() string
}

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
