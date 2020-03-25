package handler

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
)

func (r Handler) PingHandler(b ext.Bot, u *gotgbot.Update) error {
	r.TelegramProvider.Init(u)
	go r.TelegramProvider.SendText(
		"<b>Pong...</b>",
		u.EffectiveChat.Id,
		0,
		nil,
	)
	return nil
}