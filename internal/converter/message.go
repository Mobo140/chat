package converter

import (
	"github.com/Mobo140/microservices/chat/internal/model"
	desc "github.com/Mobo140/microservices/chat/pkg/chat_v1"
)

func ToMessageFromDesc(message *desc.Message) *model.MessageInfo {
	return &model.MessageInfo{
		From:      message.Form,
		Text:      message.Text,
		Timestamp: message.Timestamp.AsTime(),
	}
}
