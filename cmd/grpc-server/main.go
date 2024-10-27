package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/Mobo140/microservices/chat/internal/config"
	"github.com/Mobo140/microservices/chat/internal/config/env"
	chatRepository "github.com/Mobo140/microservices/chat/internal/repository/chat"
	messageRepository "github.com/Mobo140/microservices/chat/internal/repository/message"
	chatService "github.com/Mobo140/microservices/chat/internal/service/chat"
	chatAPI "github.com/Mobo140/microservices/chat/internal/transport/chat_api"
	desc "github.com/Mobo140/microservices/chat/pkg/chat_v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to load grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to load pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer pool.Close()

	chatRepo := chatRepository.NewRepository(pool)
	messageRepo := messageRepository.NewRepository(pool)
	chatAPIService := chatService.NewService(chatRepo, messageRepo)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatV1Server(s, chatAPI.NewImplementation(chatAPIService))

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
