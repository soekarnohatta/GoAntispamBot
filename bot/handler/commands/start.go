package commands

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"

	"GoAntispamBot/bot/helpers/trans"
	"GoAntispamBot/bot/providers/telegramProvider"
)

type CommandStart struct {
	TelegramProvider telegramProvider.TelegramProvider
}

func (r CommandStart) Start(b ext.Bot, u *gotgbot.Update) error {
	r.TelegramProvider.Init(u)
	msg := u.EffectiveMessage

	if msg.Chat.Type == "private" {
		go r.TelegramProvider.SendText(
			trans.GetString(msg.Chat.Id, "actions/startpm"),
			msg.Chat.Id,
			0,
			nil,
		)
		return nil
	}

	go r.TelegramProvider.SendText(
		trans.GetString(msg.Chat.Id, "actions/start"),
		msg.Chat.Id,
		0,
		nil,
	)
	return nil
}
