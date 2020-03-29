/*
Package "services" is a package that provides services to be used by other funcs.
This package should has all services for the bot.
*/
package settingsService

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"

	"GoAntispamBot/bot/helpers/errHandler"
	"GoAntispamBot/bot/model"
	"GoAntispamBot/bot/providers/mongoProvider"
)

func UpdateGroupSetting(setting model.GroupSetting) {
	// Start updating...
	go mongoProvider.Update("setting", setting.ChatID, bson.M{"$set": setting}, true)
}

func UpdateSetting(chatID int, set string) error {
	find := FindGroupSetting(chatID)
	if find == nil {
		UpdateGroupSetting(model.GroupSetting{
			ChatID:         chatID,
			Gban:           true,
			Username:       true,
			ProfilePicture: true,
			Time:           0,
		})
	}

	if find != nil {
		switch set {
		case "gban":
			find.Gban = true
			UpdateGroupSetting(*find)
			return nil
		case "username":
			find.Username = true
			UpdateGroupSetting(*find)
			return nil
		case "profilepicture":
			find.ProfilePicture = true
			UpdateGroupSetting(*find)
			return nil
		default:
			return errors.New("setting type is not correct")
		}
	}
	return errors.New("setting type is not correct")
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
