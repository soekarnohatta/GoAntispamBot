package welcomeService

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"GoAntispamBot/bot/helpers/errHandler"
	"GoAntispamBot/bot/model"
	"GoAntispamBot/bot/providers/mongoProvider"
)

var ctx = context.Background()

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
	go mongoProvider.Update("welcome", chatID, bson.M{"$set": mu}, true)
}

// HasUserClickedButton Has the user clicked button to unmute themselves
func HasUserClickedButton(userID int, chatID int) bool {
	mu := &model.MutedUser{UserId: userID, ChatId: chatID}
	SESSION.FirstOrInit(mu)
	return mu.ButtonClicked
}

// IsUserHuman Is the user a human
func IsUserHuman(userID, chatID string) bool {
	mu := &MutedUser{UserId: userID, ChatId: chatID}
	return SESSION.First(mu).RowsAffected != 0
}

// SetWelcPref Set whether to welcome or not
func SetWelcPref(chatID string, pref bool) {
	w := &Welcome{ChatId: chatID}
	tx := SESSION.Begin()
	tx.FirstOrCreate(w)
	w.ShouldWelcome = pref
	tx.Save(w)
	tx.Commit()
}

// SetCustomWelcome Set the custom welcome string
func SetCustomWelcome(chatID string, welcome string, buttons []WelcomeButton, welcType int) {
	w := &Welcome{ChatId: chatID}
	if buttons == nil {
		buttons = make([]WelcomeButton, 0)
	}

	tx := SESSION.Begin()
	prevButtons := make([]WelcomeButton, 0)
	tx.Where(&WelcomeButton{ChatId: chatID}).Find(&prevButtons)
	for _, btn := range prevButtons {
		tx.Delete(&btn)
	}

	for _, btn := range buttons {
		tx.Save(&btn)
	}

	tx.FirstOrCreate(w)
	w.CustomWelcome = welcome
	w.WelcomeType = welcType
	tx.Save(w)
	tx.Commit()
}

// GetDelPref Get Whether to delete service messages or not
func GetDelPref(chatID string) bool {
	return GetWelcomePrefs(chatID).DelJoined
}

// SetDelPref Set whether to delete service messages or not
func SetDelPref(chatID string, pref bool) {
	w := &Welcome{ChatId: chatID}
	tx := SESSION.Begin()
	tx.FirstOrCreate(w)
	w.DelJoined = pref
	tx.Save(w)
	tx.Commit()
}

// SetMutePref Set whether to mute users when they join or not
func SetMutePref(chatID string, pref bool) {
	w := &Welcome{ChatId: chatID}
	tx := SESSION.Begin()
	tx.FirstOrCreate(w)
	w.ShouldMute = pref
	tx.Save(w)
	tx.Commit()
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
