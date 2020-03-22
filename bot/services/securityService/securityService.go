/*
Package "services" is a package that provides services to be used by other funcs.
This package should has all services for the bot.
*/
package securityService

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"

	"GoAntispamBot/bot/model"
	"GoAntispamBot/bot/providers/errorProvider"
	"GoAntispamBot/bot/providers/mongoProvider"
)

func UpdateGlobalBan(sec model.GlobalBan) {
	// Start updating...
	go mongoProvider.Update("security", sec.UserID, bson.M{"$set": sec}, true)
}

func RemoveGlobalBan(userID int) {
	// Initiate & fills the empty struct
	var secStruct = model.GlobalBan{}
	secStruct.UserID = userID

	// Start removing...
	go mongoProvider.Remove("security", secStruct.UserID)
}

func FindGlobalBan(userID int) (*model.GlobalBan, error) {
	// Initiate & fills the empty struct
	var secStruct = &model.GlobalBan{}
	secStruct.UserID = userID

	// Start search...
	a := mongoProvider.FindOne("security", secStruct.UserID)
	if a != nil {
		_ = bson.Unmarshal(a, secStruct)
		return secStruct, nil
	}
	return nil, errors.New(errorProvider.UserInvalid)
}
