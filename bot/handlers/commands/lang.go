package commands

import (
	"github.com/PaulSonOfLars/goloc"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"

	"GoAntispamBot/bot/helpers/chatStatus"
	"GoAntispamBot/bot/helpers/trans"
	"GoAntispamBot/bot/services"
)

func (r CommandHandler) SetLang(b ext.Bot, u *gotgbot.Update, args []string) error {
	r.TelegramProvider.Init(u)
	msg := u.EffectiveMessage

	if !chatStatus.RequireAdmin(msg.From.Id, r.TelegramProvider) {
		return nil
	}

	if args != nil {
		if !goloc.IsLangSupported(args[0]) {
			go r.TelegramProvider.SendText(
				trans.GetString(msg.Chat.Id, "error/langnotsupp"),
				msg.Chat.Id,
				0,
				nil,
			)
			return nil
		}

		services.UpdateLang(msg.Chat.Id, args[0])
		return nil
	}
	return nil
}
