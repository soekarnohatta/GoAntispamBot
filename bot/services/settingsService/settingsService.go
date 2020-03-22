/*
Package "services" is a package that provides services to be used by other funcs.
This package should has all services for the bot.
*/
package settingsService

import (
	"go.mongodb.org/mongo-driver/bson"

	"GoAntispamBot/bot/helpers/errHandler"
	"GoAntispamBot/bot/model"
	"GoAntispamBot/bot/providers/mongoProvider"
)

func UpdateGroupSetting(setting model.GroupSetting) {
	// Start updating...
	go mongoProvider.Update("setting", setting.ChatID, bson.M{"$set": setting}, true)
}

func RemoveGroupSetting(chatID int) {
	// Initiate & fills the empty struct
	var settingStruct = model.GroupSetting{}
	settingStruct.ChatID = chatID

	// Start removing...
	go mongoProvider.Remove("setting", settingStruct.ChatID)
}

func FindGroupSetting(chatID int) *model.GroupSetting {
	// Initiate & fills the empty struct
	var settingStruct = &model.GroupSetting{}
	settingStruct.ChatID = chatID

	// Start search...
	a := mongoProvider.FindOne("setting", settingStruct.ChatID)
	if a != nil {
		err := bson.Unmarshal(a, settingStruct)
		errHandler.Error(err)
		return settingStruct
	}

	// else if the record is not found, update manually.
	settingStruct.Username = true
	settingStruct.ProfilePicture = true
	settingStruct.Gban = true
	UpdateGroupSetting(*settingStruct)
	return settingStruct
}
