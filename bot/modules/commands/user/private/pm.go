package private

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"github.com/sirupsen/logrus"

	"github.com/jumatberkah/antispambot/bot"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
)

var startButtons = [][]ext.InlineKeyboardButton{
	{ext.InlineKeyboardButton{
		Text:         "📝 Help",
		CallbackData: "start(help)",
	}},
	{ext.InlineKeyboardButton{
		Text:         "🇦🇺 Language",
		CallbackData: "start(language)",
	}},
	{ext.InlineKeyboardButton{
		Text: "🔗 Add Me To Your Groups",
		Url:  fmt.Sprintf("https://t.me/%v?startgroup=new", ext.Bot{}.UserName),
	}}}

var infoButtons = [][]ext.InlineKeyboardButton{
	{ext.InlineKeyboardButton{
		Text: "📡 Help",
		Url:  fmt.Sprintf("https://t.me/%v?start=help", ext.Bot{}.UserName),
	}}}

func start(b ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat

	if chat.Type == "private" {
		if len(args) != 0 {
			switch args[0] {
			case "help":
				btnList := function.BuildKeyboardf(
					"data/keyboard/help.json",
					2,
					map[string]string{"1": b.UserName},
				)

				markup := ext.InlineKeyboardMarkup{InlineKeyboard: &btnList}
				replyText := fmt.Sprintf(
					"*%v Version* `%v`\n"+
						"by *PolyDev\n\n*",
					b.FirstName,
					bot.BotConfig.BotVer,
				)

				replyText += function.GetString(chat.Id, "modules/helpers/help.go:helptxt")
				reply := b.NewSendableMessage(chat.Id, replyText)
				reply.ReplyMarkup = &markup
				reply.ReplyToMessageId = msg.MessageId
				reply.ParseMode = parsemode.Markdown
				_, err := reply.Send()
				return err
			default:
				txtStart := function.GetStringf(
					chat.Id,
					"modules/private/pm.go:start",
					map[string]string{"1": bot.BotConfig.BotVer, "2": b.FirstName},
				)

				replyMsg := b.NewSendableMessage(chat.Id, txtStart)
				replyMsg.ParseMode = "Markdown"
				replyMsg.ReplyMarkup = &ext.InlineKeyboardMarkup{InlineKeyboard: &startButtons}
				replyMsg.ReplyToMessageId = msg.MessageId
				_, err := replyMsg.Send()
				return err
			}
		}

		txtStart := function.GetStringf(
			chat.Id,
			"modules/private/pm.go:start",
			map[string]string{"1": bot.BotConfig.BotVer, "2": b.FirstName},
		)

		replyMsg := b.NewSendableMessage(chat.Id, txtStart)
		replyMsg.ParseMode = "Markdown"
		replyMsg.ReplyMarkup = &ext.InlineKeyboardMarkup{InlineKeyboard: &startButtons}
		replyMsg.ReplyToMessageId = msg.MessageId
		_, err := replyMsg.Send()
		return err
	} else {
		replyText := function.GetString(chat.Id, "modules/helpers/help.go:noprivate")
		reply := b.NewSendableMessage(chat.Id, replyText)
		reply.ReplyMarkup = &ext.InlineKeyboardMarkup{InlineKeyboard: &infoButtons}
		reply.ReplyToMessageId = msg.MessageId
		reply.ParseMode = parsemode.Markdown
		_, err := reply.Send()
		return err
	}
}

func LoadPm(u *gotgbot.Updater) {
	defer logrus.Info("PM Module Loaded...")
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("start", []rune{'/', '.'}, start))
}
