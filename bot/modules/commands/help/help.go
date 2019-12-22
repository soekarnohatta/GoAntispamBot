package help

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
	"github.com/sirupsen/logrus"
)

func initHelpButtons() ext.InlineKeyboardMarkup {
	helpButtons := [][]ext.InlineKeyboardButton{make([]ext.InlineKeyboardButton, 2), make([]ext.InlineKeyboardButton, 2),
		make([]ext.InlineKeyboardButton, 2), make([]ext.InlineKeyboardButton, 1)}

	// First column
	helpButtons[0][0] = ext.InlineKeyboardButton{
		Text:         "Sudo",
		CallbackData: fmt.Sprintf("help(%v)", "sudo"),
	}
	helpButtons[1][0] = ext.InlineKeyboardButton{
		Text:         "Username",
		CallbackData: fmt.Sprintf("help(%v)", "username"),
	}
	helpButtons[2][0] = ext.InlineKeyboardButton{
		Text:         "Picture",
		CallbackData: fmt.Sprintf("help(%v)", "picture"),
	}
	helpButtons[3][0] = ext.InlineKeyboardButton{
		Text:         "Notification",
		CallbackData: fmt.Sprintf("help(%v)", "notif"),
	}

	// Second column
	helpButtons[0][1] = ext.InlineKeyboardButton{
		Text:         "Anti Spam",
		CallbackData: fmt.Sprintf("help(%v)", "aspam"),
	}
	helpButtons[1][1] = ext.InlineKeyboardButton{
		Text:         "Verify",
		CallbackData: fmt.Sprintf("help(%v)", "verify"),
	}
	helpButtons[2][1] = ext.InlineKeyboardButton{
		Text:         "Privacy Policy",
		CallbackData: fmt.Sprintf("help(%v)", "ppolicy"),
	}
	markup := ext.InlineKeyboardMarkup{InlineKeyboard: &helpButtons}
	return markup
}

func help(b ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	if chat.Type != "private" {
		infoButtons := [][]ext.InlineKeyboardButton{make([]ext.InlineKeyboardButton, 1)}
		infoButtons[0][0] = ext.InlineKeyboardButton{
			Text: "📡 Help",
			Url:  fmt.Sprintf("https://t.me/%v?start=start", b.UserName),
		}
		replyText := function.GetString(chat.Id, "modules/private/pm.go:start")
		reply := b.NewSendableMessage(chat.Id, replyText)
		reply.ReplyMarkup = &ext.InlineKeyboardMarkup{&infoButtons}
		reply.ReplyToMessageId = msg.MessageId
		reply.ParseMode = parsemode.Html
		_, err := reply.Send()
		return err
	} else {
		markup := initHelpButtons()
		replyText := function.GetString(chat.Id, "modules/helpers/help.go:helptxt")
		reply := b.NewSendableMessage(chat.Id, replyText)
		reply.ReplyMarkup = &markup
		reply.ReplyToMessageId = msg.MessageId
		reply.ParseMode = parsemode.Html
		_, err := reply.Send()
		return err
	}
}

func LoadHelp(u *gotgbot.Updater) {
	defer logrus.Info("Help Module Loaded...")
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("help", []rune{'/', '.'}, help))
}
