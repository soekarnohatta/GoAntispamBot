package sql

import (
	"strconv"
)

type NewUser struct {
	ChatId string `gorm:"primary_key"`
	UserId string `gorm:"not null"`
	Date   string `gorm:"not null"`
}

func UpdateNewUser(chatid int, userid string, date string) error {
	tx := SESSION.Begin()

	set := &NewUser{ChatId: strconv.Itoa(chatid), UserId: userid, Date: date}
	tx.Where(NewUser{ChatId: strconv.Itoa(chatid)}).Assign(NewUser{ChatId: strconv.Itoa(chatid)}).FirstOrCreate(set)
	ret := tx.Commit().Error
	return ret
}

func DelNewUser(userid int) bool {
	tx := SESSION.Begin()
	filter := &NewUser{UserId: strconv.Itoa(userid)}

	if tx.First(filter).RowsAffected == 0 {
		tx.Rollback()
		return false
	}

	tx.Delete(filter)
	tx.Commit()

	return true
}

func GetNewUser(userid int) *NewUser {
	ver := &NewUser{UserId: strconv.Itoa(userid)}

	if SESSION.First(ver).RowsAffected == 0 {
		return nil
	}
	return ver
}
