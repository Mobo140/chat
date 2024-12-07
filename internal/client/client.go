package client

import "context"

type AccessServiceClient interface {
	Check(ctx context.Context, endpoint string) error
}
