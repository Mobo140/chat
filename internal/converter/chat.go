package converter

import (
	"errors"

	"github.com/Mobo140/microservices/chat/internal/model"
	desc "github.com/Mobo140/microservices/chat/pkg/chat_v1"
)

func ToChatInfoFromDesc(info *desc.ChatInfo) (*model.ChatInfo, error) {
	if info == nil {
		return nil, errors.New("chatInfo is empty")
	}

	return &model.ChatInfo{
		Usernames: info.Usernames,
	}, nil
}

func ToChatFromService(chat *model.Chat) *desc.Chat {
	return &desc.Chat{
		Id:   chat.ID,
		Info: ToChatInfoFromService(chat.Info),
	}
}

func ToChatInfoFromService(chatInfo model.ChatInfo) *desc.ChatInfo {
	return &desc.ChatInfo{
		Usernames: chatInfo.Usernames,
	}
}
