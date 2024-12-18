package chat

import (
	"sync"

	desc "github.com/Mobo140/microservices/chat/pkg/chat_v1"
)

type Chat struct {
	streams map[string]desc.ChatV1_ConnectChatServer
	m       sync.RWMutex
}

func NewChat() *Chat {
	return &Chat{
		streams: make(map[string]desc.ChatV1_ConnectChatServer),
	}
}
