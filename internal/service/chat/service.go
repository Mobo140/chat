package chat

import (
	"context"

	"github.com/Mobo140/microservices/chat/internal/model"
	"github.com/Mobo140/microservices/chat/internal/repository"
	"github.com/Mobo140/microservices/chat/internal/service"
)

var _ service.ChatService = (*serv)(nil)

type serv struct {
	chatRepository    repository.ChatRepository
	messageRepository repository.MessageRepository
}

func NewService(chatRepository repository.ChatRepository, messageRepository repository.MessageRepository) *serv {
	return &serv{chatRepository: chatRepository, messageRepository: messageRepository}
}

func (s *serv) Create(ctx context.Context, info *model.ChatInfo) (int64, error) {

	id, err := s.chatRepository.Create(ctx, info)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *serv) Get(ctx context.Context, id int64) (*model.Chat, error) {

	chat, err := s.chatRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *serv) Delete(ctx context.Context, id int64) error {
	return s.chatRepository.Delete(ctx, id)
}

func (s *serv) SendMessage(ctx context.Context, message *model.Message) error {
	return s.messageRepository.SendMessage(ctx, message)
}
