package sql

import (
	"strconv"
)

type Antispam struct {
	ChatId   string `gorm:"primary_key"`
	Arabs    string `gorm:"not null"`
	Deletion string `gorm:"not null"`
	Forward  string `gorm:"not null"`
	Link     string `gorm:"not null"`
}

func UpdateAntispam(chatid int, arabs string, del string, fwd string, link string) {
	tx := SESSION.Begin()

	antispam := &Antispam{ChatId: strconv.Itoa(chatid), Arabs: arabs, Deletion: del, Forward: fwd, Link: link}
	tx.Where(Antispam{ChatId: strconv.Itoa(chatid)}).Assign(Antispam{Arabs: arabs,
		Deletion: del, Forward: fwd, Link: link}).FirstOrCreate(antispam)
	tx.Commit()
}

func DelAntispam(chatid int) bool {
	tx := SESSION.Begin()
	filter := &Antispam{ChatId: strconv.Itoa(chatid)}

	if tx.First(filter).RowsAffected == 0 {
		tx.Rollback()
		return false
	}

	tx.Delete(filter)
	tx.Commit()

	return true
}

func GetAntispam(chatid int) *Antispam {
	opt := &Antispam{ChatId: strconv.Itoa(chatid)}

	if SESSION.First(opt).RowsAffected == 0 {
		return nil
	}
	return opt
}
