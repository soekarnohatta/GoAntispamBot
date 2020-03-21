/*
Package "providers" is a package that provides required things reqired by the bot
to be used by other funcs.
This package should has all providers for the bot.
*/
package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"GoAntispamBot/bot/helpers/errHandler"
	"GoAntispamBot/bot/helpers/trans"
	"GoAntispamBot/bot/model"
	"GoAntispamBot/bot/services"
)

var myClient = &http.Client{Timeout: 5 * time.Second}

func IsGlobalBan(userID int) bool {
	_, err := services.FindGlobalBan(userID) // Find a user from the DB.
	if err != nil {
		// Theoritically if there's an error it would be likely the user is not a
		// spammer because the DB returned record not found error.
		return false
	}
	return true
}

func IsCASBan(userID int) bool {
	// Request data to CAS API.
	cas := fmt.Sprintf("https://api.cas.chat/check?user_id=%v", userID)
	r, err := myClient.Get(cas)
	if err != nil {
		return false
	}
	defer r.Body.Close()

	// Deserialize it...
	var ban model.CasBan
	err = json.NewDecoder(r.Body).Decode(&ban)
	errHandler.Error(err)
	return ban.Ok
}

func FilterSpamUser(telegramProvider TelegramProvider) {
	msg := telegramProvider.Message

	if IsGlobalBan(msg.From.Id) {
		doBanSpammer(telegramProvider)
		return
	} else if IsCASBan(msg.From.Id) {
		doBanSpammer(telegramProvider)
		return
	}
}

func doBanSpammer(telegramProvider TelegramProvider) {
	msg := telegramProvider.Message
	go telegramProvider.KickMember(msg.From.Id, msg.Chat.Id, -1)
	go telegramProvider.SendText(
		trans.GetStringf(telegramProvider.Message.Chat.Id, "actions/spammer", map[string]string{"1": msg.Text}),
		0,
		0,
		nil,
	)
}
