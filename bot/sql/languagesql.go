package sql

import (
	"strconv"
)

type Lang struct {
	ChatId string `gorm:"primary_key"`
	Lang   string `gorm:"not null"`
}

func UpdateLang(chatid int, lang string) {
	tx := SESSION.Begin()

	set := &Lang{ChatId: strconv.Itoa(chatid), Lang: lang}
	tx.Where(Lang{ChatId: strconv.Itoa(chatid)}).Assign(Lang{Lang: lang}).FirstOrCreate(set)
	tx.Commit()
}

func DelLang(chatid int) bool {
	tx := SESSION.Begin()
	filter := &Lang{ChatId: strconv.Itoa(chatid)}

	if tx.First(filter).RowsAffected == 0 {
		tx.Rollback()
		return false
	}

	tx.Delete(filter)
	tx.Commit()

	return true
}

func GetLang(chatid int) *Lang {
	ver := &Lang{ChatId: strconv.Itoa(chatid)}

	if SESSION.First(ver).RowsAffected == 0 {
		return nil
	}
	return ver
}
