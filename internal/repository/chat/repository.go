package chat

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Mobo140/microservices/chat/internal/model"
	"github.com/Mobo140/microservices/chat/internal/repository"
	"github.com/Mobo140/microservices/chat/internal/repository/chat/converter"
	modelRepo "github.com/Mobo140/microservices/chat/internal/repository/chat/model"
	"github.com/Mobo140/platform_common/pkg/db"
)

var _ repository.ChatRepository = (*chatRepo)(nil)

const (
	tableName       = "chat"
	usernamesColumn = "usernames"
	idColumn        = "id"
)

type chatRepo struct {
	db db.Client
}

func NewRepository(db db.Client) *chatRepo { //nolint:revive // it's ok
	return &chatRepo{db: db}
}

func (r *chatRepo) Create(ctx context.Context, info *model.ChatInfo) (int64, error) {
	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(usernamesColumn).
		Values(info.Usernames).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %v", err)
	}

	q := db.Query{
		QueryRow: query,
		Name:     "chat_repository.create",
	}

	var chatID int64

	err = r.db.DB().ScanOneContext(ctx, &chatID, q, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to insert chat: %v", err)
	}

	return chatID, nil
}

func (r *chatRepo) Get(ctx context.Context, id int64) (*model.Chat, error) {
	builderSelect := sq.Select(idColumn, usernamesColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	var chat modelRepo.Chat

	q := db.Query{
		QueryRow: query,
		Name:     "chat_repository.get",
	}

	err = r.db.DB().ScanOneContext(ctx, &chat, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to select chat: %v", err)
	}

	return converter.ToChatFromRepo(&chat), nil
}

// func (r *repo) Update(ctx context.Context, *model.UpdateInfo) error {

// 	builderSelect := sq.Update(tableName).
// 			Set(usernamesColumn, squirrel.Expr("array_append("+usernamesColumn+", ?)", username)).
// 			PlaceholderFormat(squirrel.Dollar).
// 			Where(squirrel.And{
// 				squirrel.Eq{idColumn: id},
// 				squirrel.Expr("NOT (? = ANY("+usernamesColumn+"))", username),
// 			})

// 	query, args, err := builderSelect.ToSql()
// 	if err != nil {
// 		log.Fatalf("failed to build query: %v", err)
// 	}

// }

func (r *chatRepo) Delete(ctx context.Context, id int64) error {
	builderDelete := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %v", err)
	}

	q := db.Query{
		QueryRow: query,
		Name:     "chat_repository.delete",
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to delete chat: %v", err)
	}

	return nil
}
