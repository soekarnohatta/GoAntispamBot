package sql

import (
	"strconv"
)

type Verify struct {
	ChatId   string `gorm:"primary_key"`
	Option   string `gorm:"not null"`
	Deletion string `gorm:"not null"`
	Text     string `gorm:"not null"`
}

func UpdateVerify(chatId int, option string, text string, del string) error {
	tx := SESSION.Begin()

	set := &Verify{ChatId: strconv.Itoa(chatId), Option: option, Text: text, Deletion: del}
	tx.Where(Verify{ChatId: strconv.Itoa(chatId)}).Assign(Verify{Option: option,
		Text: text, Deletion: del}).FirstOrCreate(set)
	ret := tx.Commit().Error
	return ret
}

func DelVerify(chatid int) bool {
	tx := SESSION.Begin()
	filter := &Verify{ChatId: strconv.Itoa(chatid)}

	if tx.First(filter).RowsAffected == 0 {
		tx.Rollback()
		return false
	}

	tx.Delete(filter)
	tx.Commit()

	return true
}

func GetVerify(chatid int) *Verify {
	ver := &Verify{ChatId: strconv.Itoa(chatid)}

	if SESSION.First(ver).RowsAffected == 0 {
		return nil
	}
	return ver
}
