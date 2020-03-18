package handlers

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"

	"GoAntispamBot/bot/providers"
)

type Ping struct {
	telegramProvider *providers.TelegramProvider
}

func (r *Ping) PingHandler(b ext.Bot, u *gotgbot.Update) {
	r.telegramProvider.Init(u)
	go r.telegramProvider.SendText(
		"<b>Pong...</b>",
		u.EffectiveChat.Id,
		0,
		nil,
	)
}
