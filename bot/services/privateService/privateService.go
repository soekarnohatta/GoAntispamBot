/*
Package "services" is a package that provides services to be used by other funcs.
This package should has all services for the bot.
*/
package privateService

import (
	"go.mongodb.org/mongo-driver/bson"

	"GoAntispamBot/bot/model"
	"GoAntispamBot/bot/providers/mongoProvider"
)

func UpdateNotification(userID int, notif bool) {
	// Initiate & fills the empty struct
	var notifStruct = model.Private{}
	notifStruct.UserID = userID
	notifStruct.Notif = notif

	// Start updating...
	go mongoProvider.Update("notif", notifStruct.UserID, bson.M{"$set": notifStruct}, true)
}

func RemoveNotification(userID int) {
	// Initiate & fills the empty struct
	var notifStruct = model.Private{}
	notifStruct.UserID = userID

	// Start removing...
	go mongoProvider.Remove("notif", notifStruct.UserID)
}

func FindNotification(userID int) (notif string) {
	// Initiate & fills the empty struct
	var notifStruct = model.Private{}
	notifStruct.UserID = userID

	// Start search...
	res := mongoProvider.FindOne("notif", notifStruct.UserID)
	if res != nil {
		notif = string(res.Lookup("notification").Value)
	}
	return
}
