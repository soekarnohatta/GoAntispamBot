package language

import (
	"fmt"
	"github.com/PaulSonOfLars/goloc"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/caching"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/chat_status"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
	"github.com/sirupsen/logrus"
)

func setlang(b ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat_status.RequireUserAdmin(chat, msg, user.Id) == false {
		return nil
	}

	if len(args) == 0 {
		_, err := msg.ReplyText("Please insert the language code so that i can change your language")
		err_handler.HandleErr(err)
		return err
	}

	if !goloc.IsLangSupported(args[0]) {
		_, err := msg.ReplyText(function.GetString(chat.Id, "modules/language/language.go:58"))
		err_handler.HandleErr(err)
		return err
	}

	_, err := caching.REDIS.Set(fmt.Sprintf("lang_%v", chat.Id), args[0], 7200).Result()
	if err != nil {
		err = sql.UpdateLang(chat.Id, args[0])
		err_handler.HandleTgErr(b, u, err)
		_, err = msg.ReplyText(function.GetStringf(chat.Id, "modules/language/language.go:51",
			map[string]string{"1": args[0]}))
		return err
	}
	_, err = msg.ReplyText(function.GetStringf(chat.Id, "modules/language/language.go:51",
		map[string]string{"1": args[0]}))
	err_handler.HandleErr(err)
	return err
}

func LoadLang(u *gotgbot.Updater) {
	defer logrus.Info("Lang Module Loaded...")
	go u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("setlang", []rune{'/', '.'}, setlang))
}
