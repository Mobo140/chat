package converter

import (
	model "github.com/Mobo140/chat/internal/model"
	modelRepo "github.com/Mobo140/chat/internal/repository/chat/model"
)

func ToChatFromRepo(chat *modelRepo.Chat) *model.Chat {
	return &model.Chat{
		ID:   chat.ID,
		Info: ToChatInfoFromRepo(chat.Info),
	}
}

func ToChatInfoFromRepo(chatInfo modelRepo.ChatInfo) model.ChatInfo {
	return model.ChatInfo{
		Usernames: chatInfo.Usernames,
	}
}
