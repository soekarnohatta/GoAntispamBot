package sql

import (
	"strconv"
)

type Notification struct {
	ChatId       string `gorm:"primary_key"`
	Notification string `gorm:"not null"`
}

func UpdateNotification(chatid int, notification string) error {
	tx := SESSION.Begin()

	set := &Notification{ChatId: strconv.Itoa(chatid), Notification: notification}
	tx.Where(Notification{ChatId: strconv.Itoa(chatid)}).Assign(Notification{Notification: notification}).FirstOrCreate(set)
	tx.Commit()
	return tx.Error
}

func DelNotification(chatid int) bool {
	tx := SESSION.Begin()
	filter := &Notification{ChatId: strconv.Itoa(chatid)}

	if tx.First(filter).RowsAffected == 0 {
		tx.Rollback()
		return false
	}

	tx.Delete(filter)
	tx.Commit()
	return true
}

func GetNotification(chatid int) *Notification {
	ver := &Notification{ChatId: strconv.Itoa(chatid)}

	if SESSION.First(ver).RowsAffected == 0 {
		return nil
	}
	return ver
}
