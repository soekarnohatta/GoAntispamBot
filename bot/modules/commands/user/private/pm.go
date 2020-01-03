package private

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"github.com/jumatberkah/antispambot/bot"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
	"github.com/jumatberkah/antispambot/bot/modules/listener"
	"github.com/sirupsen/logrus"
)

func start(b ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat

	if len(args) != 0 {
		switch args[0] {
		case "help":
			markup := listener.InitHelpButtons()
			replyText := fmt.Sprintf("*%v Version* `%v`\n"+
				"by *PolyDev\n\n*", b.FirstName, bot.BotConfig.BotVer)
			replyText += function.GetString(chat.Id, "modules/helpers/help.go:helptxt")
			reply := b.NewSendableMessage(chat.Id, replyText)
			reply.ReplyMarkup = &markup
			reply.ReplyToMessageId = msg.MessageId
			reply.ParseMode = parsemode.Markdown
			_, err := reply.Send()
			return err
		default:
			if chat.Type == "private" {
				startButtons := [][]ext.InlineKeyboardButton{make([]ext.InlineKeyboardButton, 2)}
				startButtons[0][0] = ext.InlineKeyboardButton{
					Text:         "📝 Help",
					CallbackData: fmt.Sprintf("start(%v)", "help"),
				}
				startButtons[0][1] = ext.InlineKeyboardButton{
					Text: "🔗 Add Me To Your Groups",
					Url:  "https://t.me/PolyesterBot?startgroup=new",
				}

				txtStart := function.GetStringf(chat.Id, "modules/private/pm.go:start",
					map[string]string{"1": bot.BotConfig.BotVer})
				replyMsg := b.NewSendableMessage(chat.Id, txtStart)
				replyMsg.ParseMode = "Markdown"
				replyMsg.ReplyMarkup = &ext.InlineKeyboardMarkup{&startButtons}
				replyMsg.ReplyToMessageId = msg.MessageId
				_, err := replyMsg.Send()
				return err
			} else {
				infoButtons := [][]ext.InlineKeyboardButton{make([]ext.InlineKeyboardButton, 1)}
				infoButtons[0][0] = ext.InlineKeyboardButton{
					Text: "📡 Help",
					Url:  fmt.Sprintf("https://t.me/%v?start=help", b.UserName),
				}
				replyText := function.GetString(chat.Id, "modules/helpers/help.go:noprivate")
				reply := b.NewSendableMessage(chat.Id, replyText)
				reply.ReplyMarkup = &ext.InlineKeyboardMarkup{&infoButtons}
				reply.ReplyToMessageId = msg.MessageId
				reply.ParseMode = parsemode.Markdown
				_, err := reply.Send()
				return err
			}
		}
	}

	if chat.Type == "private" {
		startButtons := [][]ext.InlineKeyboardButton{make([]ext.InlineKeyboardButton, 2)}
		startButtons[0][0] = ext.InlineKeyboardButton{
			Text:         "📝 Help",
			CallbackData: fmt.Sprintf("start(%v)", "help"),
		}
		startButtons[0][1] = ext.InlineKeyboardButton{
			Text: "🔗 Add Me To Your Groups",
			Url:  "https://t.me/PolyesterBot?startgroup=new",
		}

		txtStart := function.GetStringf(chat.Id, "modules/private/pm.go:start",
			map[string]string{"1": bot.BotConfig.BotVer, "2": b.FirstName})
		replyMsg := b.NewSendableMessage(chat.Id, txtStart)
		replyMsg.ParseMode = "Markdown"
		replyMsg.ReplyMarkup = &ext.InlineKeyboardMarkup{&startButtons}
		replyMsg.ReplyToMessageId = msg.MessageId
		_, err := replyMsg.Send()
		return err
	} else {
		infoButtons := [][]ext.InlineKeyboardButton{make([]ext.InlineKeyboardButton, 1)}
		infoButtons[0][0] = ext.InlineKeyboardButton{
			Text: "📡 Help",
			Url:  fmt.Sprintf("https://t.me/%v?start=help", b.UserName),
		}
		replyText := function.GetString(chat.Id, "modules/helpers/help.go:noprivate")
		reply := b.NewSendableMessage(chat.Id, replyText)
		reply.ReplyMarkup = &ext.InlineKeyboardMarkup{&infoButtons}
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
