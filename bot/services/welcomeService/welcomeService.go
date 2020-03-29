package welcomeService

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"GoAntispamBot/bot/helpers/errHandler"
	"GoAntispamBot/bot/model"
	"GoAntispamBot/bot/providers/mongoProvider"
)

var ctx = context.Background()

const (
	TEXT        = 0
	BUTTON_TEXT = 1
	STICKER     = 2
	DOCUMENT    = 3
	PHOTO       = 4
	AUDIO       = 5
	VOICE       = 6
	VIDEO       = 7
)

func GetWelcomePrefs(chatID int) *model.Welcome {
	var welcStruct = &model.Welcome{}
	welcStruct.ChatId = chatID

	// Start search...
	a := mongoProvider.FindOne("welcome", welcStruct.ChatId)
	if a != nil {
		err := bson.Unmarshal(a, welcStruct)
		errHandler.Error(err)
		return welcStruct
	}

	// else if the record is not found, update manually.
	welcStruct.CleanWelcome = 1
	welcStruct.CustomWelcome = ""
	welcStruct.DelJoined = true
	welcStruct.ShouldMute = true
	welcStruct.ShouldWelcome = true
	//UpdateGroupSetting(*welcStruct)
	return welcStruct
}

// GetWelcomeButtons Get the buttons for the welcome message
func GetWelcomeButtons(chatID int) []model.WelcomeButton {
	return findButton(chatID)
}

// SetCleanWelcome Set whether to clean old welcome messages or not
func SetCleanWelcome(chatID int, cw int) {
	w := GetWelcomePrefs(chatID)
	if w != nil {
		w.CleanWelcome = cw
		go mongoProvider.Update("welcome", chatID, bson.M{"$set": w}, true)
	}
}

// GetCleanWelcome Get whether to clean old welcome messages or not
func GetCleanWelcome(chatID int) int {
	return GetWelcomePrefs(chatID).CleanWelcome
}

// UserClickedButton Mark the user as a human
func UserClickedButton(userID int, chatID int) {
	mu := &model.MutedUser{UserId: userID, ChatId: chatID, ButtonClicked: true}
	go mongoProvider.Update("welcomebutton", chatID, bson.M{"$set": mu}, true)
}

// HasUserClickedButton Has the user clicked button to unmute themselves
func HasUserClickedButton(userID int, chatID int) bool {
	mu := &model.MutedUser{UserId: userID, ChatId: chatID}
	a := mongoProvider.FindOne("welcomebutton", mu)
	if a != nil {
		err := bson.Unmarshal(a, mu)
		errHandler.Error(err)
		return mu.ButtonClicked
	}
	return mu.ButtonClicked
}

// IsUserHuman Is the user a human
func IsUserHuman(userID int, chatID int) bool {
	mu := &model.MutedUser{UserId: userID, ChatId: chatID}
	return mongoProvider.FindOne("welcomebutton", mu) != nil
}

// SetWelcPref Set whether to welcome or not
func SetWelcPref(chatID int, pref bool) {
	w := &model.Welcome{ChatId: chatID}
	w.ShouldWelcome = pref
	go mongoProvider.Update("welcome", chatID, bson.M{"$set": w}, true)

}

// SetCustomWelcome Set the custom welcome string
func SetCustomWelcome(chatID int, welcome string, buttons []model.WelcomeButton, welcType int) {
	w := &model.Welcome{ChatId: chatID}
	if buttons == nil {
		buttons = make([]model.WelcomeButton, 0)
	}

	prevButtons := findButton(chatID)
	for _, btn := range prevButtons {
		go mongoProvider.Remove("welcomebutton", btn)
	}

	for _, btn := range buttons {
		go mongoProvider.Update("welcomebutton", chatID, btn, true)
	}

	w.CustomWelcome = welcome
	w.WelcomeType = welcType
	go mongoProvider.Update("welcome", chatID, w, true)
}

// GetDelPref Get Whether to delete service messages or not
func GetDelPref(chatID int) bool {
	return GetWelcomePrefs(chatID).DelJoined
}

// SetDelPref Set whether to delete service messages or not
func SetDelPref(chatID int, pref bool) {
	w := GetWelcomePrefs(chatID)
	w.DelJoined = pref
	go mongoProvider.Update("welcome", chatID, w, true)
}

// SetMutePref Set whether to mute users when they join or not
func SetMutePref(chatID int, pref bool) {
	w := GetWelcomePrefs(chatID)
	w.ShouldMute = pref
	go mongoProvider.Update("welcome", chatID, w, true)
}

//----------------------------------------------------------------------------------------------------------------------

func findButton(chatID int) []model.WelcomeButton {
	var buttonStruct = model.WelcomeButton{}
	var res []model.WelcomeButton
	buttonStruct.ChatId = chatID

	db, err := mongoProvider.Connect() // Init connection
	errHandler.Fatal(err)

	// Start searching...
	csr, _ := db.Collection("welcomebutton").Find(ctx, buttonStruct)
	for csr.Next(ctx) {
		err := csr.Decode(buttonStruct)
		errHandler.Fatal(err)
		res = append(res, buttonStruct)
	}

	errHandler.Error(csr.Err())
	csr.Close(ctx)
	return res
}
