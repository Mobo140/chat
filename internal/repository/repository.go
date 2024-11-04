package repository

import (
	"context"

	"github.com/Mobo140/microservices/chat/internal/model"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *model.ChatInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.Chat, error)
	// Update(ctx context.Context, info *model.UpdateInfo) error
	Delete(ctx context.Context, id int64) error
}

type MessageRepository interface {
	SendMessage(ctx context.Context, message *model.Message) error
	//GetMessagesByChatID()
}

type LogRepository interface {
	Create(ctx context.Context, logEntry *model.LogEntry) error
}
