package main

import (
	"context"
	"log"
	"time"

	desc "github.com/Mobo140/microservices/chat/pkg/chat_v1"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:8082"
	chatID  = 1
)

func main() {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	c := desc.NewChatV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &desc.GetRequest{Id: chatID})
	if err != nil {
		log.Fatalf("failed to create chat:%v", err)
	}

	log.Printf(color.RedString("Chat info:\n"), color.GreenString("%+v", r.GetChat()))
}
