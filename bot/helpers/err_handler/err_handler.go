package err_handler

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	log "github.com/sirupsen/logrus"
)

type CommandCallback func()

func HandleErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func FatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func HandleTgErr(b ext.Bot, u *gotgbot.Update, err error) {
	if err != nil {
		var msg = u.EffectiveMessage
		_, err = msg.ReplyText(err.Error())
		HandleErr(err)
	}
}

func HandleCbErr(b ext.Bot, u *gotgbot.Update, err error) {
	if err != nil {
		var msg = u.CallbackQuery
		_, err = b.AnswerCallbackQueryText(msg.Id, err.Error(), true)
	}
}
