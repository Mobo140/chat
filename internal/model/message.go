package model

import "time"

type Message struct {
	From      string
	Text      string
	CreatedAt time.Time
}

type MessageInfo struct {
	ChatID    int64
	Message   Message
	Timestamp time.Time
}

type SendMessage struct {
	ChatID  int64
	Message Message
}
