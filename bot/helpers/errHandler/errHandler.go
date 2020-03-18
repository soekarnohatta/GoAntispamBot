/*
Package "errHandler" is a package that handles all kinds of error(s).
This package should handle all error(s).
*/
package errHandler

import (
	"github.com/PaulSonOfLars/gotgbot"
	log "github.com/sirupsen/logrus"

	"GoAntispamBot/bot/helpers/trans"
	"GoAntispamBot/bot/providers"
)

var telegramProvider = providers.TelegramProvider{}

// Error function returns nothing as it only handles error and log it.
func Error(err error) {
	if err != nil {
		log.Error(err)
	}
}

// Fatal function returns nothing as it only handles error and log it.
func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// SendError function will send an error message to the chat.
func SendError(err error, u *gotgbot.Update) {
	if err != nil {
		telegramProvider.Init(u)
		go telegramProvider.SendText(
			trans.GetStringf(u.EffectiveChat.Id, "error/error", map[string]string{"1": err.Error()}),
			u.EffectiveChat.Id,
			0,
			nil,
		)
	}
}
