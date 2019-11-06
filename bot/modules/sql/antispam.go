package sql

import (
	"strconv"
)

type Antispam struct {
	ChatId   string `gorm:"primary_key"`
	Option   string `gorm:"not null"`
	Arabs    string `gorm:"not null"`
	Deletion string `gorm:"not null"`
	Forward  string `gorm:"not null"`
	Link     string `gorm:"not null"`
	NonLatin string `gorm:"not null"`
}

func UpdateAntispam(chatid int, option string, arabs string, del string, fwd string, link string, nonlatin string) error {
	tx := SESSION.Begin()

	antispam := &Antispam{ChatId: strconv.Itoa(chatid), Option: option,
		Arabs: arabs, Deletion: del, Forward: fwd, Link: link, NonLatin: nonlatin}
	tx.Where(Antispam{ChatId: strconv.Itoa(chatid)}).Assign(Antispam{Option: option,
		Arabs: arabs, Deletion: del, Forward: fwd, Link: link, NonLatin: nonlatin}).FirstOrCreate(antispam)
	tx.Commit()
	return tx.Error
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
