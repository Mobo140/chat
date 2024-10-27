package model

import "time"

type Message struct {
	ChatID int64
	Info   MessageInfo
}
type MessageInfo struct {
	From      string
	Text      string
	Timestamp time.Time
}
