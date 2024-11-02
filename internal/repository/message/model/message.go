package model

import "time"

type Message struct {
	ChatID int64       `db:"chat_id"`
	Info   MessageInfo `db:""`
}
type MessageInfo struct {
	From      string    `db:"from_user"`
	Text      string    `db:"text"`
	Timestamp time.Time `db:"timestamp"`
}
