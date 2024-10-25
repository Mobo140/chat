package main

import (
	"context"
	"fmt"
	"log"
	"net"

	sq "github.com/Masterminds/squirrel"
	desc "github.com/Mobo140/microservices/chat-server/pkg/chat_v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = 8082
	dbDSN    = "host=localhost port=50002 user=n1234 password=qwerty dbname=chat sslmode=disable"
)

type server struct {
	desc.UnimplementedChatV1Server
	pool *pgxpool.Pool
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {

	builderInsert := sq.Insert("chat").
		PlaceholderFormat(sq.Dollar).
		Columns("usernames").
		Values(req.Info.Usernames).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var chatID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&chatID)
	if err != nil {
		log.Fatalf("failed to insert chat: %v", err)
	}

	log.Printf("inserted chat with id: %d", chatID)

	return &desc.CreateResponse{
		Id: chatID,
	}, nil

}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {

	builderSelect := sq.Select("usernames").
		From("chat").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.Id}).
		Limit(1)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var usernames []string

	err = s.pool.QueryRow(ctx, query, args...).Scan(&usernames)
	if err != nil {
		log.Fatalf("failed to select chat: %v", err)
	}

	log.Printf("usernames: %s", usernames)

	return &desc.GetResponse{
		Info: &desc.ChatInfo{
			Usernames: usernames,
		},
	}, nil
}

func main() {

	ctx := context.Background()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.New(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatV1Server(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
