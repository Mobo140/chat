package app

import (
	"context"
	"log"

	"github.com/Mobo140/microservices/chat/internal/client/db"
	"github.com/Mobo140/microservices/chat/internal/client/db/pg"
	transaction "github.com/Mobo140/microservices/chat/internal/client/db/transaction"
	"github.com/Mobo140/microservices/chat/internal/closer"
	"github.com/Mobo140/microservices/chat/internal/config"
	"github.com/Mobo140/microservices/chat/internal/config/env"
	"github.com/Mobo140/microservices/chat/internal/repository"
	chatRepository "github.com/Mobo140/microservices/chat/internal/repository/chat"
	logRepository "github.com/Mobo140/microservices/chat/internal/repository/logs"
	messageRepository "github.com/Mobo140/microservices/chat/internal/repository/message"
	"github.com/Mobo140/microservices/chat/internal/service"
	chatService "github.com/Mobo140/microservices/chat/internal/service/chat"
	"github.com/Mobo140/microservices/chat/internal/transport/handlers/chat"
	chatHandler "github.com/Mobo140/microservices/chat/internal/transport/handlers/chat"
)

type serviceProvider struct {
	chatRepository    repository.ChatRepository
	messageRepository repository.MessageRepository
	logRepository     repository.LogRepository

	grpcConfig config.GRPCConfig
	pgConfig   config.PGConfig
	txManager  db.TxManager
	dbClient   db.Client

	chatService service.ChatService

	chatImplementation *chat.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) ChatHandler(ctx context.Context) *chat.Implementation {
	if s.chatImplementation == nil {
		s.chatImplementation = chatHandler.NewImplementation(s.ChatAPIService(ctx))
	}

	return s.chatImplementation
}

func (s *serviceProvider) ChatAPIService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = chatService.NewService(
			s.ChatRepository(ctx),
			s.MessageRepository(ctx),
			s.LogRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.chatService
}

func (s *serviceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if s.chatRepository == nil {
		s.chatRepository = chatRepository.NewRepository(s.DBClient(ctx))
	}

	return s.chatRepository
}

func (s *serviceProvider) MessageRepository(ctx context.Context) repository.MessageRepository {
	if s.messageRepository == nil {
		s.messageRepository = messageRepository.NewRepository(s.DBClient(ctx))
	}

	return s.messageRepository
}

func (s *serviceProvider) LogRepository(ctx context.Context) repository.LogRepository {
	if s.logRepository == nil {
		s.logRepository = logRepository.NewRepository(s.DBClient(ctx))
	}

	return s.logRepository
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to initialize grpc config: %v", err)
		}
		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to initialize pg config: %v", err)
		}
		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to initialize db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %v", err)
		}
		closer.Add(cl.Close)
		s.dbClient = cl
	}

	return s.dbClient
}