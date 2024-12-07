package access

import (
	"context"
	"flag"
	"fmt"
	"log"

	descAccess "github.com/Mobo140/auth/pkg/access_v1"
	"github.com/Mobo140/microservices/chat/internal/config"
	"github.com/Mobo140/microservices/chat/internal/config/env"
	"github.com/Mobo140/microservices/chat/internal/model"
	"github.com/Mobo140/platform_common/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type client struct {
	accessClient descAccess.AccessV1Client
}

func New(accessClient descAccess.AccessV1Client) *client {
	return &client{
		accessClient: accessClient,
	}
}

func (c *client) Check(ctx context.Context, enpoint string) error {
	_, err := c.accessClient.Check(ctx, &descAccess.CheckRequest{
		EndpointAddress: enpoint,
	})

	if err != nil {
		logger.Error("Access denied: ", zap.Error(err))

		return err
	}

	logger.Info("Access granted")

	return nil
}

var (
	configPath  string
	accessToken string
)

func setupFlags() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
	flag.StringVar(&accessToken, "a", "", "access token")
	flag.Parse()
}

func main() {
	setupFlags()

	// при вызове метода мы передаем в заголовке access токен вызываем метод Get например 
	ctx := context.Background()
	md := metadata.New(map[string]string{"Authorization": "Bearer " + accessToken})
	ctx = metadata.NewOutgoingContext(ctx, md)

	creds, err := credentials.NewClientTLSFromFile("../../../auth.pem", "")
	if err != nil {
		log.Fatalf("could not process the credentials: %v", err)
	}

	if err := config.Load(configPath); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	cfg, err := env.NewAccessClientConfig()
	if err != nil {
		log.Fatalf("failed to create auth client config: %v", err)
	}

	conn, err := grpc.NewClient(cfg.Address(), grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}

	cl := descAccess.NewAccessV1Client(conn)

	_, err = cl.Check(ctx, &descAccess.CheckRequest{
		EndpointAddress: model.EndpointPath,
	})
	if err != nil {
		log.Fatalf(err.Error()) //nolint:govet // it's ok
	}

	fmt.Println("Access granted")
}
