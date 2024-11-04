package service

import (
	"context"

	"github.com/Mobo140/microservices/chat/internal/model"
)

type ChatService interface {
	Create(ctx context.Context, chat *model.ChatInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.Chat, error)
	Delete(ctx context.Context, id int64) error
	// Update(ctx context.Context, info *model.UpdateInfo) error
	SendMessage(ctx context.Context, message *model.SendMessage) error
	// GetMessagesByChatID()
}
