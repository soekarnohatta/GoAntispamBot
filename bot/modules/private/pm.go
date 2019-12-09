package private

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
	"github.com/sirupsen/logrus"
)

func start(_ ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat

	txtStart := function.GetString(chat.Id, "modules/private/pm.go:15")

	if chat.Type == "supergroup" {
		_, err := msg.Delete()
		return err
	}

	_, err := msg.ReplyHTML(txtStart)
	return err
}

func LoadPm(u *gotgbot.Updater) {
	defer logrus.Info("PM Module Loaded...")
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("start", []rune{'/', '.'}, start))
}
