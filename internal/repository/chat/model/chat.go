package model

type Chat struct {
	ID   int64
	Info ChatInfo
}

type ChatInfo struct {
	Usernames []string
}

// type UpdateInfo struct {
// 	ID       int64
// 	UserName string
// }
