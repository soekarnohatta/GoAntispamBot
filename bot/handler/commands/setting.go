package commands

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"

	"GoAntispamBot/bot/helpers/chatStatus"
	"GoAntispamBot/bot/providers/telegramProvider"
)

type CommandSetting struct {
	TelegramProvider telegramProvider.TelegramProvider
}

func (r CommandSetting) Setting(b ext.Bot, u *gotgbot.Update, args []string) error {
	r.TelegramProvider.Init(u)
	msg := u.EffectiveMessage

	if chatStatus.RequirePrivate(r.TelegramProvider) {
		if len(args) >= 1 {
			for _, val := range args {

			}
		}
	}
}
