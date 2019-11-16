package private

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
)

func start(_ ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat

	txtStart := function.GetString(chat.Id, "modules/private/pm.go:15")

	if chat.Type == "supergroup" {
		_, err := msg.Delete()
		err_handler.HandleErr(err)
		return gotgbot.EndGroups{}
	}
	_, err := msg.ReplyHTML(txtStart)
	err_handler.HandleErr(err)
	return gotgbot.EndGroups{}
}

// LoadPm -> Register handlers
func LoadPm(u *gotgbot.Updater) {
	go u.Dispatcher.AddHandler(handlers.NewPrefixCommand("start", []rune{'/', '.'}, start))
}
