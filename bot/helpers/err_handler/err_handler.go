package err_handler

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/sirupsen/logrus"
)

func HandleErr(err error) {
	if err != nil {
		logrus.Error(err)
	}
}

func FatalError(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}

func HandleTgErr(b ext.Bot, u *gotgbot.Update, err error) {
	if err != nil {
		msg := u.EffectiveMessage
		rep := b.NewSendableMessage(msg.Chat.Id, err.Error())
		rep.ReplyToMessageId = msg.MessageId
		_, err := rep.Send()
		if err != nil {
			if err.Error() == "Bad Request: reply message not found" {
				rep.ReplyToMessageId = 0
				_, _ = rep.Send()
			}
		}
	}
}
