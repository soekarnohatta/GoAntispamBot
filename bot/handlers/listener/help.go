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

var btnList = function.BuildKeyboardf("data/keyboard/help.json", 2, map[string]string{"1": "Polyesterbot"})
var backButton = [][]ext.InlineKeyboardButton{{ext.InlineKeyboardButton{
	Text:         "Back",
	CallbackData: "help(back)",
}}}

func handleHelp(b ext.Bot, u *gotgbot.Update) error {
	query := u.CallbackQuery
	pattern, _ := regexp.Compile(`help\((.+?)\)`)

	if pattern.MatchString(query.Data) {
		module := pattern.FindStringSubmatch(query.Data)[1]
		chat := u.EffectiveChat
		replyText := fmt.Sprintf("*%v Version* `%v`\n"+
			"by *PolyDev\n\n*",
			b.FirstName,
			bot.BotConfig.BotVer,
		)

		replyText += function.GetString(chat.Id, "handlers/helpers/help.go:helptxt")
		msg := b.NewSendableMessage(chat.Id, replyText)
		msg.ParseMode = parsemode.Markdown
		backKeyboard := ext.InlineKeyboardMarkup{InlineKeyboard: &backButton}
		msg.ReplyMarkup = &backKeyboard
		if module != "back" {
			replyTxt := fmt.Sprintf("*%v Version* `%v`\n"+
				"by *PolyDev\n\n*", b.FirstName, bot.BotConfig.BotVer)
			replyTxt += function.GetString(chat.Id, "handlers/helpers/help.go:"+module)
			msg.Text = replyTxt
		} else if module == "back" {
			markup := ext.InlineKeyboardMarkup{InlineKeyboard: &btnList}
			msg.ReplyMarkup = &markup
		}

		_, err := query.Message.Delete()
		err_handler.HandleErr(err)
		_, err = msg.Send()
		err_handler.HandleErr(err)
	}
	return gotgbot.ContinueGroups{}
}

func LoadHelpListener(u *gotgbot.Updater) {
	defer logrus.Info("Help Listener Loaded...")
	u.Dispatcher.AddHandler(handlers.NewCallback("help", handleHelp))
}
