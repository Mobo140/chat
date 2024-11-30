package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	descAccess "github.com/Mobo140/auth/pkg/access_v1"
	"github.com/Mobo140/microservices/chat/internal/config/env"
	"github.com/Mobo140/microservices/chat/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var accessToken = flag.String("a", "", "access token")

func main() {
	flag.Parse()

	ctx := context.Background()
	md := metadata.New(map[string]string{"Authorization": "Bearer" + *accessToken})
	ctx = metadata.NewOutgoingContext(ctx, md)

	creds, err := credentials.NewClientTLSFromFile("auth.pem", "")
	if err != nil {
		log.Fatalf("could not process the credentials: %v", err)
	}

	cfg, err := env.NewAuthGRPCConfig()
	if err != nil {
		log.Fatalf("failed to load auth grpc config: %v", err)
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
