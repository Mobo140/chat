package model

import "time"

type LogEntry struct {
	ChatID    int64
	Activity  string
	CreatedAt time.Time
}
