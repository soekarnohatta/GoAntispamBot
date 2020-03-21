/*
Package "errHandler" is a package that handles all kinds of error(s).
This package should handle all error(s).
*/
package errHandler

import (
	log "github.com/sirupsen/logrus"

	"GoAntispamBot/bot/helpers/trans"
	"GoAntispamBot/bot/providers"
)

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
func SendError(err error, telegramProvider providers.TelegramProvider) {
	if err != nil {
		go telegramProvider.SendText(
			trans.GetStringf(telegramProvider.Message.Chat.Id, "error/error", map[string]string{"1": err.Error()}),
			0,
			0,
			nil,
		)
	}
}
