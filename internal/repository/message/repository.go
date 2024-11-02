package model

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/Mobo140/microservices/chat/internal/client/db"
	"github.com/Mobo140/microservices/chat/internal/model"
)

const (
	tableName       = "message"
	chatIDColumn    = "chat_id"
	fromUserColumn  = "from_user"
	textColumn      = "text"
	timestampColumn = "timestamp"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) *repo {
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

	q := db.Query{
		QueryRow: query,
		Name:     "send_message_repository",
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		log.Fatalf("failed to insert message: %v", err)
	}

	log.Printf("inserted message: %v", message)

	return nil
}
