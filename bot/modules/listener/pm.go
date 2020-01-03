package listener

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"github.com/jumatberkah/antispambot/bot"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
	"github.com/sirupsen/logrus"
	"regexp"
)

func handleStart(b ext.Bot, u *gotgbot.Update) error {
	query := u.CallbackQuery
	pattern, _ := regexp.Compile(`start\((.+?)\)`)

	if pattern.MatchString(query.Data) {
		module := pattern.FindStringSubmatch(query.Data)[1]
		chat := u.EffectiveChat

		switch module {
		case "help":
			markup := InitHelpButtons()
			replyText := fmt.Sprintf("*%v Version* `%v`\n"+
				"by *PolyDev\n\n*", b.FirstName, bot.BotConfig.BotVer)
			replyText += function.GetString(chat.Id, "modules/helpers/help.go:helptxt")
			reply := b.NewSendableEditMessageText(chat.Id, query.Message.MessageId, replyText)
			reply.ReplyMarkup = &markup
			reply.ParseMode = parsemode.Markdown
			_, err := reply.Send()
			err_handler.HandleErr(err)
		}
	}
	return gotgbot.ContinueGroups{}
}

func LoadStartListener(u *gotgbot.Updater) {
	defer logrus.Info("PM Listener Loaded...")
	u.Dispatcher.AddHandler(handlers.NewCallback("start", handleStart))
}
