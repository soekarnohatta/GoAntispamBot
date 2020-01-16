package help

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

func help(b ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	if chat.Type != "private" {
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

}

func LoadHelp(u *gotgbot.Updater) {
	defer logrus.Info("Help Module Loaded...")
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("help", []rune{'/', '.'}, help))
}
