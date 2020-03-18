package model

type Backup struct {
	BotID   int        `json:"bot_id"`
	Data    DataBackup `json:"data"`
	Version int        `json:"version"`
}

type DataBackup struct {
	Language     []Lang         `json:"lang"`
	Notification []Private      `json:"notif"`
	Gban         []GlobalBan    `json:"gban"`
	Setting      []GroupSetting `json:"setting"`
	CLog         []ChatLog      `json:"chat_log"`
	Ulog         []UserLog      `json:"user_log"`
}
