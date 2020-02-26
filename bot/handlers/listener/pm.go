package listener

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"github.com/sirupsen/logrus"
	"regexp"

	"github.com/jumatberkah/antispambot/bot"
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/helpers/function"
)

func handleStart(b ext.Bot, u *gotgbot.Update) error {
	query := u.CallbackQuery
	pattern, _ := regexp.Compile(`start\((.+?)\)`)

	if pattern.MatchString(query.Data) {
		module := pattern.FindStringSubmatch(query.Data)[1]
		chat := u.EffectiveChat

		switch module {
		case "help":
			markup := ext.InlineKeyboardMarkup{InlineKeyboard: &btnList}
			replyText := fmt.Sprintf(
				"*%v Version* `%v`\n"+
					"by *PolyDev\n\n*",
				b.FirstName,
				bot.BotConfig.BotVer,
			)
			replyText += function.GetString(chat.Id, "handlers/helpers/help.go:helptxt")
			reply := b.NewSendableMessage(chat.Id, replyText)
			reply.ReplyMarkup = &markup
			reply.ParseMode = parsemode.Markdown
			_, err := query.Message.Delete()
			err_handler.HandleErr(err)
			_, err = reply.Send()
			err_handler.HandleErr(err)
		case "language":
			var btnLang = function.BuildKeyboardf(
				"data/keyboard/language.json",
				2,
				map[string]string{"1": fmt.Sprint(chat.Id)},
			)
			_, err := query.Message.Delete()
			err_handler.HandleErr(err)
			newMsg := b.NewSendableMessage(chat.Id, "*Available Language(s):*")
			newMsg.ParseMode = parsemode.Markdown
			newMsg.ReplyMarkup = &ext.InlineKeyboardMarkup{InlineKeyboard: &btnLang}
			_, err = newMsg.Send()
			return err
		}
	}
	return gotgbot.ContinueGroups{}
}

func LoadStartListener(u *gotgbot.Updater) {
	defer logrus.Info("PM Listener Loaded...")
	u.Dispatcher.AddHandler(handlers.NewCallback("start", handleStart))
}
