package access

import (
	"context"

	descAccess "github.com/Mobo140/auth/pkg/access_v1"
	"github.com/Mobo140/platform_common/pkg/logger"
	"go.uber.org/zap"
)

type client struct {
	accessClient descAccess.AccessV1Client
}

func New(accessClient descAccess.AccessV1Client) *client {
	return &client{
		accessClient: accessClient,
	}
}

func (c *client) Check(ctx context.Context, endpoint string) error {
	_, err := c.accessClient.Check(ctx, &descAccess.CheckRequest{
		EndpointAddress: endpoint,
	})

	if err != nil {
		logger.Error("Access denied: ", zap.Error(err))

		return err
	}

	logger.Info("Access granted")

	return nil
}

//  при вызове метода мы передаем в заголовке access токен вызываем метод Get например
// ctx := context.Background()
// md := metadata.New(map[string]string{"Authorization": "Bearer " + accessToken})
// ctx = metadata.NewOutgoingContext(ctx, md)
