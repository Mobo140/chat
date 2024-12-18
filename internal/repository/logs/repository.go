package logs

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/Mobo140/chat/internal/model"
	"github.com/Mobo140/chat/internal/repository"
	"github.com/Mobo140/platform_common/pkg/db"
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

func NewRepository(db db.Client) *logRepo { //nolint:revive // it's ok
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
