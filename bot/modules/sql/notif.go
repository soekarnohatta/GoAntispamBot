package sql

import (
	"strconv"
)

type Notification struct {
	ChatId       string `gorm:"primary_key"`
	Notification string `gorm:"not null"`
}

func UpdateNotification(chatid int, notification string) {
	tx := SESSION.Begin()

	set := &Notification{ChatId: strconv.Itoa(chatid), Notification: notification}
	tx.Save(set)
	tx.Commit()

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
