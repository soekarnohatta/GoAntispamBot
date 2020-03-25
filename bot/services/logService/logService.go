package logService

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"

	"GoAntispamBot/bot/model"
	"GoAntispamBot/bot/providers/errorProvider"
	"GoAntispamBot/bot/providers/mongoProvider"
)

// UpdateUser function will update an user information.
func UpdateUser(val model.UserLog) {
	// Start updating...
	go mongoProvider.Update("user_log", val.UserID, bson.M{"$set": val}, true)
}

// RemoveUser function will remove an user information.
func RemoveUser(userID int) {
	// Initiate & fills the empty struct
	var userStruct = model.UserLog{}
	userStruct.UserID = userID

	// Start removing...
	go mongoProvider.Remove("user_log", userStruct.UserID)
}

// FindUser function will find an user information.
func FindUser(userID int) (*model.UserLog, error) {
	// Initiate & fills the empty struct
	var userStruct = &model.UserLog{}
	userStruct.UserID = userID

	// Start search...
	a := mongoProvider.FindOne("user_log", userStruct.UserID)
	if a != nil {
		_ = bson.Unmarshal(a, userStruct)
		return userStruct, nil
	}
	return nil, errors.New(errorProvider.UserInvalid)
}

// UpdateChat function will update a chat information.
func UpdateChat(val model.ChatLog) {
	// Start updating...
	go mongoProvider.Update("chat_log", val.ChatID, bson.M{"$set": val}, true)
}

// RemoveChat function will remove a chat information.
func RemoveChat(chatID int) {
	// Initiate & fills the empty struct
	var chatStruct = model.ChatLog{}
	chatStruct.ChatID = chatID

	// Start removing...
	go mongoProvider.Remove("chat_log", chatStruct.ChatID)
}

// FindChat function will find a chat information.
func FindChat(chatID int) (*model.ChatLog, error) {
	// Initiate & fills the empty struct
	var chatStruct = &model.ChatLog{}
	chatStruct.ChatID = chatID

	// Start search...
	a := mongoProvider.FindOne("user_log", chatStruct.ChatID)
	if a != nil {
		_ = bson.Unmarshal(a, chatStruct)
		return chatStruct, nil
	}
	return nil, errors.New(errorProvider.UserInvalid)
}