package sql

import (
	"strconv"
)

type Username struct {
	ChatId   string `gorm:"primary_key"`
	Option   string `gorm:"not null"`
	Action   string `gorm:"not null"`
	Deletion string `gorm:"not null"`
	Text     string `gorm:"not null"`
}

func UpdateUsername(chatid int, option string, action string, text string, del string) error {
	tx := SESSION.Begin()

	// upsert spam user
	username := &Username{ChatId: strconv.Itoa(chatid), Option: option,
		Action: action, Text: text, Deletion: del}
	tx.Where(Username{ChatId: strconv.Itoa(chatid)}).Assign(Username{Option: option,
		Action: action, Text: text, Deletion: del}).FirstOrCreate(username)
	tx.Commit()
	return tx.Error
}

func DelUsername(chatid int) bool {
	tx := SESSION.Begin()
	filter := &Username{ChatId: strconv.Itoa(chatid)}

	if tx.First(filter).RowsAffected == 0 {
		tx.Rollback()
		return false
	}

	tx.Delete(filter)
	tx.Commit()
	return true
}

func GetUsername(chatid int) *Username {
	opt := &Username{ChatId: strconv.Itoa(chatid)}

	if SESSION.First(opt).RowsAffected == 0 {
		return nil
	}
	return opt
}
