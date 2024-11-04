package model

import "time"

type LogEntry struct {
	UserID    int64     `db:"user_id"`
	Action    string    `db:"action"`
	Timestamp time.Time `db:"timestamp"`
}
