package chat

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mobo140/chat/internal/model"
	"github.com/Mobo140/chat/internal/repository"
	"github.com/Mobo140/chat/internal/service"
	"github.com/Mobo140/platform_common/pkg/db"
)

var _ service.ChatService = (*serv)(nil)

const (
	unknownChat = -1
)

type serv struct {
	chatRepository    repository.ChatRepository
	messageRepository repository.MessageRepository
	logRepository     repository.LogRepository
	txManager         db.TxManager
}

func NewService(
	chatRepository repository.ChatRepository,
	messageRepository repository.MessageRepository,
	logRepository repository.LogRepository,
	txManager db.TxManager,
) *serv { //nolint:revive // it's ok
	return &serv{
		chatRepository:    chatRepository,
		messageRepository: messageRepository,
		logRepository:     logRepository,
		txManager:         txManager,
	}
}

func (s *serv) Create(ctx context.Context, info *model.ChatInfo) (int64, error) {
	var id int64
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var errTx error

		id, errTx = s.chatRepository.Create(ctx, info)
		if errTx != nil {
			return errTx
		}

		logEntry := model.LogEntry{
			ChatID:   id,
			Activity: fmt.Sprintf("Create chat: usernames:%s", strings.Join(info.Usernames, ", ")),
		}

		errTx = s.logRepository.Create(ctx, &logEntry)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return unknownChat, err
	}

	return id, nil
}

func (s *serv) Get(ctx context.Context, id int64) (*model.Chat, error) {
	var chat *model.Chat
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var errTx error

		chat, errTx = s.chatRepository.Get(ctx, id)
		if errTx != nil {
			return errTx
		}

		logEntry := model.LogEntry{
			ChatID:   id,
			Activity: fmt.Sprintf("Get chat: Id:%d, Usernames:%s", id, strings.Join(chat.Info.Usernames, ", ")),
		}

		errTx = s.logRepository.Create(ctx, &logEntry)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *serv) Delete(ctx context.Context, id int64) error {
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var errTx error

		errTx = s.chatRepository.Delete(ctx, id)
		if errTx != nil {
			return errTx
		}

		logEntry := model.LogEntry{
			ChatID:   id,
			Activity: fmt.Sprintf("Delete chat: ID=%d", id),
		}

		errTx = s.logRepository.Create(ctx, &logEntry)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *serv) SendMessage(ctx context.Context, message *model.SendMessage) error {
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var errTx error

		errTx = s.messageRepository.SendMessage(ctx, message)
		if errTx != nil {
			return errTx
		}

		logEntry := model.LogEntry{
			ChatID: message.ChatID,
			Activity: fmt.Sprintf(
				"Send message to chat: ChatID:%d, From:%s, Text:%s, CreatedAt:%s",
				message.ChatID,
				message.Message.From,
				message.Message.Text,
				message.Message.CreatedAt,
			),
		}

		errTx = s.logRepository.Create(ctx, &logEntry)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
