package model

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/Mobo140/microservices/chat/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	tableName       = "message"
	chatIDColumn    = "chat_id"
	fromUserColumn  = "from_user"
	textColumn      = "text"
	timestampColumn = "timestamp"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *repo {
	return &repo{db: db}
}

func (r *repo) SendMessage(ctx context.Context, message *model.Message) error {

	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(chatIDColumn, fromUserColumn, textColumn).
		Values(message.ChatID, message.Info.From, message.Info.Text)

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed too build query: %v", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to insert message: %v", err)
	}

	log.Printf("inserted message: %v", message)

	return nil
}
