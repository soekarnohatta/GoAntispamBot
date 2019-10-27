package sql

import (
	"strconv"
)

type Picture struct {
	ChatId   string `gorm:"primary_key"`
	Option   string `gorm:"not null"`
	Action   string `gorm:"not null"`
	Deletion string `gorm:"not null"`
	Text     string `gorm:"not null"`
}

func UpdatePicture(chatid int, option string, action string, text string, del string) error {
	tx := SESSION.Begin()

	// upsert picture
	pic := &Picture{Option: option, Action: action, Text: text, Deletion: del}
	tx.Where(Picture{ChatId: strconv.Itoa(chatid)}).Assign(Picture{Option: option, Action: action,
		Text: text, Deletion: del}).FirstOrCreate(pic)
	tx.Commit()
	return tx.Error
}

func DelPicture(chatid int) bool {
	filter := &Picture{ChatId: strconv.Itoa(chatid)}

	if SESSION.Delete(filter).RowsAffected == 0 {
		return false
	}
	return true
}

func GetPicture(chatid int) *Picture {
	opt := &Picture{ChatId: strconv.Itoa(chatid)}

	if SESSION.First(opt).RowsAffected == 0 {
		return nil
	}
	return opt
}
