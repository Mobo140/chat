package app

import (
	"context"
	"log"

	descAccess "github.com/Mobo140/auth/pkg/access_v1"
	"github.com/Mobo140/chat/internal/client"
	accessClient "github.com/Mobo140/chat/internal/client/access"
	"github.com/Mobo140/chat/internal/config"
	"github.com/Mobo140/chat/internal/config/env"
	"github.com/Mobo140/chat/internal/repository"
	chatRepository "github.com/Mobo140/chat/internal/repository/chat"
	logRepository "github.com/Mobo140/chat/internal/repository/logs"
	messageRepository "github.com/Mobo140/chat/internal/repository/message"
	"github.com/Mobo140/chat/internal/service"
	chatService "github.com/Mobo140/chat/internal/service/chat"
	"google.golang.org/grpc"

	"github.com/Mobo140/chat/internal/transport/handlers/chat"
	chatHandler "github.com/Mobo140/chat/internal/transport/handlers/chat"
	"github.com/Mobo140/platform_common/pkg/closer"
	"github.com/Mobo140/platform_common/pkg/db"
	"github.com/Mobo140/platform_common/pkg/db/pg"
	transaction "github.com/Mobo140/platform_common/pkg/db/transaction"
)

type serviceProvider struct {
	chatRepository    repository.ChatRepository
	messageRepository repository.MessageRepository
	logRepository     repository.LogRepository

	grpcConfig         config.GRPCConfig
	httpConfig         config.HTTPConfig
	accessClientConfig config.AccessClientConfig
	jaegerConfig       config.JaegerConfig
	pgConfig           config.PGConfig
	swaggerConfig      config.SwaggerConfig
	txManager          db.TxManager
	dbClient           db.Client

	chatService  service.ChatService
	accessClient client.AccessServiceClient

	chatImplementation *chat.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) ChatHandler(ctx context.Context, conn *grpc.ClientConn) *chat.Implementation {
	if s.chatImplementation == nil {
		s.chatImplementation = chatHandler.NewImplementation(s.ChatAPIService(ctx), s.AccessClient(conn))
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

func (s *serviceProvider) AccessClient(conn *grpc.ClientConn) client.AccessServiceClient {
	if s.accessClient == nil {
		s.accessClient = accessClient.NewAccessClient(descAccess.NewAccessV1Client(conn))
	}

	return s.accessClient
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
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
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

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to initialize http config: %v", err)
		}
		s.httpConfig = cfg
	}

	return s.httpConfig
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

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := env.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to initialize swagger config: %v", err)
		}
		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

func (s *serviceProvider) AccessClientConfig() config.AccessClientConfig {
	if s.accessClientConfig == nil {
		cfg, err := env.NewAccessClientConfig()
		if err != nil {
			log.Fatalf("failed to initialize access client config: %v", err)
		}
		s.accessClientConfig = cfg
	}

	return s.accessClientConfig
}

func (s *serviceProvider) JaegerConfig() config.JaegerConfig {
	if s.jaegerConfig == nil {
		cfg, err := env.NewJaegerConfig()
		if err != nil {
			log.Fatalf("failed to initialize jaeger config: %v", err)
		}
		s.jaegerConfig = cfg
	}

	return s.jaegerConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.NewClient(ctx, s.PGConfig().DSN())
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
