package logs

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/Mobo140/microservices/chat/internal/client/db"
	"github.com/Mobo140/microservices/chat/internal/model"
	"github.com/Mobo140/microservices/chat/internal/repository"
)

var _ repository.LogRepository = (*logRepo)(nil)

const (
	tableName       = "logs"
	chatColumn      = "chat_id"
	activityColumn  = "activity"
	createdAtColumn = "created_at"
)

type logRepo struct {
	db db.Client
}

func NewRepository(db db.Client) *logRepo {
	return &logRepo{db: db}
}

func (l *logRepo) Create(ctx context.Context, logEntry *model.LogEntry) error {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(chatColumn, activityColumn, createdAtColumn).
		Values(logEntry.ChatID, logEntry.Activity, logEntry.CreatedAt)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "create_log",
		QueryRow: query,
	}

	_, err = l.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}