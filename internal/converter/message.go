package converter

import (
	"errors"

	"github.com/Mobo140/microservices/chat/internal/model"
	desc "github.com/Mobo140/microservices/chat/pkg/chat_v1"
)

func ToMessageFromDesc(message *desc.Message) (*model.Message, error) {
	if message == nil {
		return nil, errors.New("message is empty")
	}

	return &model.Message{
		From: message.From,
		Text: message.Text,
	}, nil
}
