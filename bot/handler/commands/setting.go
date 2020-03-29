package commands

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"

	"GoAntispamBot/bot/helpers/chatStatus"
	"GoAntispamBot/bot/helpers/trans"
	"GoAntispamBot/bot/providers/telegramProvider"
	"GoAntispamBot/bot/services/settingsService"
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
				err := settingsService.UpdateSetting(msg.Chat.Id, val)
				if err != nil {
					go r.TelegramProvider.SendText(
						trans.GetString(msg.Chat.Id, "error/setinvalid"),
						msg.Chat.Id,
						0,
						nil,
					)
					return nil
				}
			}
		}
		go r.TelegramProvider.SendText(
			trans.GetString(msg.Chat.Id, "error/setinvalid"),
			msg.Chat.Id,
			0,
			nil,
		)
		return nil
	}
	return nil
}
