package model

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Mobo140/microservices/chat/internal/model"
	"github.com/Mobo140/microservices/chat/internal/repository"
	"github.com/Mobo140/platform_common/pkg/db"
)

var _ repository.MessageRepository = (*messageRepo)(nil)

const (
	tableName       = "message"
	chatIDColumn    = "chat_id"
	fromUserColumn  = "from_user"
	textColumn      = "text"
	timestampColumn = "timestamp"
)

type messageRepo struct {
	db db.Client
}

func NewRepository(db db.Client) *messageRepo { //nolint:revive // it's ok
	return &messageRepo{db: db}
}

func (r *messageRepo) SendMessage(ctx context.Context, message *model.SendMessage) error {
	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(chatIDColumn, fromUserColumn, textColumn).
		Values(message.ChatID, message.Message.From, message.Message.Text)

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %v", err)
	}

	q := db.Query{
		QueryRow: query,
		Name:     "send_message_repository",
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to insert message: %v", err)
	}

	return nil
}
