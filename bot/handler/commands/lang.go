package commands

import (
	"github.com/PaulSonOfLars/goloc"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"

	"GoAntispamBot/bot/helpers/chatStatus"
	"GoAntispamBot/bot/helpers/trans"
	"GoAntispamBot/bot/providers/telegramProvider"
	"GoAntispamBot/bot/services/langService"
)

type CommandLang struct {
	TelegramProvider telegramProvider.TelegramProvider
}

func (r CommandLang) SetLang(b ext.Bot, u *gotgbot.Update, args []string) error {
	r.TelegramProvider.Init(u)
	msg := u.EffectiveMessage

	if !chatStatus.RequireAdmin(msg.From.Id, r.TelegramProvider) {
		return nil
	}

	if args != nil {
		if !goloc.IsLangSupported(args[0]) {
			go r.TelegramProvider.ReplyText(trans.GetString(msg.Chat.Id, "error/langnotsupp"))
			return nil
		}

		langService.UpdateLang(msg.Chat.Id, args[0])
		go r.TelegramProvider.ReplyText(trans.GetString(msg.Chat.Id, "actions/changelang"))
		return nil
	}
	return nil
}
