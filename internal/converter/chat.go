package converter

import (
	"github.com/Mobo140/microservices/chat/internal/model"
	desc "github.com/Mobo140/microservices/chat/pkg/chat_v1"
)

func ToChatInfoFromDesc(info *desc.ChatInfo) *model.ChatInfo {
	return &model.ChatInfo{
		Usernames: info.Usernames,
	}
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
