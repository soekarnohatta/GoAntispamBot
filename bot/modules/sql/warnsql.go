package sql

import (
	"github.com/lib/pq"
)

type Warns struct {
	UserId   string         `gorm:"primary_key"`
	ChatId   string         `gorm:"primary_key"`
	NumWarns int            `gorm:"default:2"`
	Reasons  pq.StringArray `gorm:"type:varchar(64)[]"`
}

type WarnSettings struct {
	ChatId    string `gorm:"primary_key"`
	WarnLimit int    `gorm:"default:3"`
}

func WarnUser(userId string, chatId string, reason string) (int, []string) {
	warnedUser := &Warns{UserId: userId, ChatId: chatId}
	tx := SESSION.Begin()
	tx.FirstOrInit(warnedUser)

	// Increment warns
	warnedUser.NumWarns++

	// Add reason if it exists
	if reason != "" {
		if len(reason) >= 64 {
			reason = reason[:63]
		}
		warnedUser.Reasons = append(warnedUser.Reasons, reason)
	}

	// Upsert warn
	tx.Save(warnedUser)
	tx.Commit()

	return warnedUser.NumWarns, warnedUser.Reasons
}

func RemoveWarn(userId string, chatId string) bool {
	removed := false
	warnedUser := &Warns{UserId: userId, ChatId: chatId}
	tx := SESSION.Begin()

	tx.FirstOrInit(warnedUser)

	// only remove if user has warns
	if warnedUser.NumWarns > 0 {
		warnedUser.NumWarns -= 1
		tx.Save(warnedUser)
		removed = true
	}
	tx.Commit()

	return removed
}

func ResetWarns(userId string, chatId string) {
	warnedUser := &Warns{UserId: userId, ChatId: chatId}
	tx := SESSION.Begin()

	tx.FirstOrInit(warnedUser)

	// resetting all warn fields
	warnedUser.NumWarns = 0
	warnedUser.Reasons = make([]string, 0)
	tx.Save(warnedUser)
	tx.Commit()
}

func GetWarns(userId string, chatId string) (int, []string) {
	user := &Warns{UserId: userId, ChatId: chatId}
	SESSION.FirstOrInit(user)
	return user.NumWarns, user.Reasons
}

func SetWarnLimit(chatId string, warnLimit int) {
	warnSetting := &WarnSettings{ChatId: chatId}
	tx := SESSION.Begin()
	// init record if it doesn't exist
	tx.FirstOrInit(warnSetting)
	warnSetting.WarnLimit = warnLimit
	// upsert record
	tx.Save(warnSetting)
}

func GetWarnSetting(chatId string) int {
	warnSetting := &WarnSettings{ChatId: chatId}
	SESSION.FirstOrCreate(warnSetting)
	return warnSetting.WarnLimit
}
