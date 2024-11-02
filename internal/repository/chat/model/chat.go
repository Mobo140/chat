package model

type Chat struct {
	ID   int64    `db:"id"`
	Info ChatInfo `db:""`
}

type ChatInfo struct {
	Usernames []string `db:"usernames"`
}

// type UpdateInfo struct {
// 	ID       int64
// 	UserName string
// }
