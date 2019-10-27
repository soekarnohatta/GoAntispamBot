package sql

import (
	"strconv"
)

type UserSpam struct {
	UserId int    `gorm:"primary_key"`
	Reason string `gorm:"not null"`
}

type ChatSpam struct {
	ChatId string `gorm:"primary_key"`
	Reason string `gorm:"not null"`
}

type EnforceGban struct {
	ChatId string `gorm:"primary_key"`
	Option string `gorm:"not null"`
}

// User
func UpdateUserSpam(userid int, reason string) error {
	tx := SESSION.Begin()

	// upsert spam user
	user := &UserSpam{UserId: userid, Reason: reason}
	tx.Where(UserSpam{UserId: userid}).Assign(UserSpam{Reason: reason}).FirstOrCreate(user)
	tx.Commit()
	return tx.Error
}

func DelUserSpam(userid int) bool {
	tx := SESSION.Begin()
	filter := &UserSpam{UserId: userid}

	if tx.First(filter).RowsAffected == 0 {
		tx.Rollback()
		return false
	}

	tx.Delete(filter)
	tx.Commit()
	return true
}

func GetUserSpam(userid int) *UserSpam {
	spam := &UserSpam{UserId: userid}

	if SESSION.First(spam).RowsAffected == 0 {
		return nil
	}
	return spam
}

// Chat
func UpdateChatSpam(chatid int, reason string) error {
	tx := SESSION.Begin()

	// upsert spam user
	cht := &ChatSpam{ChatId: strconv.Itoa(chatid), Reason: reason}
	tx.Where(ChatSpam{ChatId: strconv.Itoa(chatid)}).Assign(ChatSpam{Reason: reason}).FirstOrCreate(cht)
	tx.Commit()
	return tx.Error
}

func DelChatSpam(chatid int) bool {
	tx := SESSION.Begin()
	filter := &ChatSpam{ChatId: strconv.Itoa(chatid)}

	if tx.First(filter).RowsAffected == 0 {
		tx.Rollback()
		return false
	}

	tx.Delete(filter)
	tx.Commit()
	return true
}

// Enforce Gban
func UpdateEnforceGban(chatid int, option string) error {
	tx := SESSION.Begin()

	// upsert gban enforcing
	chat := &EnforceGban{ChatId: strconv.Itoa(chatid), Option: option}
	tx.Where(EnforceGban{ChatId: strconv.Itoa(chatid)}).Assign(EnforceGban{Option: option}).FirstOrCreate(chat)
	tx.Commit()
	return tx.Error
}

func DelEnforceGban(chatid int) bool {
	tx := SESSION.Begin()
	filter := &EnforceGban{ChatId: strconv.Itoa(chatid)}

	if tx.First(filter).RowsAffected == 0 {
		tx.Rollback()
		return false
	}

	tx.Delete(filter)
	tx.Commit()
	return true

}

// Get Function
func GetChatSpam(chatid int) *ChatSpam {
	spam := &ChatSpam{ChatId: strconv.Itoa(chatid)}

	if SESSION.First(spam).RowsAffected == 0 {
		return nil
	}
	return spam
}

func GetEnforceGban(chatid int) *EnforceGban {
	spam := &EnforceGban{ChatId: strconv.Itoa(chatid)}

	if SESSION.First(spam).RowsAffected == 0 {
		return nil
	}
	return spam
}

func GetAllSpamUser() []UserSpam {
	var list []UserSpam
	SESSION.Find(&list)
	return list
}

func GetAllSpamChat() []ChatSpam {
	var list []ChatSpam
	SESSION.Find(&list)
	return list
}
