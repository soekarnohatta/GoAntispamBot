package sql

import (
	"strconv"
	"strings"
)

type User struct {
	UserId    int `gorm:"primary_key"`
	UserName  string
	FirstName string `gorm:"not null"`
	LastName  string
}

type Chat struct {
	ChatId    string `gorm:"primary_key"`
	ChatTitle string
	ChatType  string `gorm:"not null"`
	ChatLink  string
}

func UpdateUser(userId int, userName string, firstName string, lastName string) {
	username := strings.ToLower(userName)
	tx := SESSION.Begin()

	user := &User{UserId: userId, UserName: username, FirstName: firstName, LastName: lastName}
	tx.Save(user)
	tx.Commit()
}

func UpdateChat(chatid string, chattitle string, chattype string, clink string) {
	if chatid == "" {
		return
	}

	tx := SESSION.Begin()

	chat := &Chat{ChatId: chatid, ChatTitle: chattitle, ChatType: chattype, ChatLink: clink}
	tx.Save(chat)
	tx.Commit()
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

func GetUser(userId int) *User {
	ver := &User{UserId: userId}

	if SESSION.First(ver).RowsAffected == 0 {
		return nil
	}
	return ver
}

func GetChat(chatId int) *Chat {
	ver := &Chat{ChatId: strconv.Itoa(chatId)}

	if SESSION.First(ver).RowsAffected == 0 {
		return nil
	}
	return ver
}
