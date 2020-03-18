package lang

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"

	"GoAntispamBot/bot/services"
)

func setLang(b ext.Bot, u *gotgbot.Update, args []string) {
	if args != nil {
		services.UpdateLang(u.EffectiveChat.Id, args[0])
	}
}