package main

import (
	"context"
	"flag"
	"log"

	"github.com/Mobo140/microservices/chat/internal/app"
)

var (
	configPath  string
	logLevel    string
	accessToken string
)

func setupFlags() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
	flag.StringVar(&logLevel, "l", "info", "log level")
	flag.Parse()
}

func main() {
	setupFlags()

	ctx := context.Background()

	a, err := app.NewApp(ctx, configPath, logLevel)
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
