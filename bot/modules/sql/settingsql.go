package sql

import (
	"strconv"
)

type Setting struct {
	ChatId   string `gorm:"primary_key"`
	Time     string `gorm:"not null"`
	Deletion string `gorm:"not null"`
}

func UpdateSetting(chatid int, time string, delete string) {
	tx := SESSION.Begin()

	set := &Setting{Time: time, Deletion: delete}
	tx.Where(Setting{ChatId: strconv.Itoa(chatid)}).Assign(Setting{Time: time, Deletion: delete}).FirstOrCreate(set)
	tx.Commit()

}

func DelSetting(ChatId int) bool {
	tx := SESSION.Begin()
	filter := &Setting{ChatId: strconv.Itoa(ChatId)}

	if tx.First(filter).RowsAffected == 0 {
		tx.Rollback()
		return false
	}

	tx.Delete(filter)
	tx.Commit()

	return true
}

func GetSetting(chatid int) *Setting {
	tim := &Setting{ChatId: strconv.Itoa(chatid)}

	if SESSION.First(tim).RowsAffected == 0 {
		return nil
	}
	return tim
}
