package sql

import (
	"errors"
	"strings"
)

type User struct {
	UserId   int `gorm:"primary_key"`
	UserName string
	Name     string `gorm:"not null"`
}

type Chat struct {
	ChatId    string `gorm:"primary_key"`
	ChatTitle string
	ChatType  string `gorm:"not null"`
	ChatLink  string
}

func UpdateUser(userid int, username string, name string) error {
	username = strings.ToLower(username)
	tx := SESSION.Begin()

	user := &User{UserId: userid, UserName: username, Name: name}
	tx.Where(User{UserId: userid}).Assign(User{UserName: username, Name: name}).FirstOrCreate(user)
	ret := tx.Commit().Error
	return ret
}

func UpdateChat(chatid string, chattitle string, chattype string, clink string) error {
	if chatid == "" {
		return errors.New("Chat Title Should Not Nil")
	}

	tx := SESSION.Begin()

	chat := &Chat{ChatId: chatid, ChatTitle: chattitle, ChatType: chattype, ChatLink: clink}
	tx.Where(Chat{ChatId: chatid}).Assign(Chat{ChatTitle: chattitle, ChatType: chattype,
		ChatLink: clink}).FirstOrCreate(chat)
	ret := tx.Commit().Error
	return ret
}

func DelUser(userid int) bool {
	tx := SESSION.Begin()
	user := &User{UserId: userid}

	if tx.First(user).RowsAffected == 0 {
		tx.Rollback()
		return false
	}

	tx.Delete(user)
	tx.Commit()

	return true
}

func DelChat(chatid string) bool {
	tx := SESSION.Begin()
	chat := &Chat{ChatId: chatid}

	if tx.First(chat).RowsAffected == 0 {
		tx.Rollback()
		return false
	}

	tx.Delete(chat)
	tx.Commit()

	return true
}

func GetUserIdByName(username string) *User {
	username = strings.ToLower(username)
	user := new(User)
	SESSION.Where("user_name = ?", username).First(user)
	return user
}

func GetAllChat() []Chat {
	var list []Chat
	SESSION.Find(&list)
	return list
}

func GetAllUser() []User {
	var list []User
	SESSION.Find(&list)
	return list
}
