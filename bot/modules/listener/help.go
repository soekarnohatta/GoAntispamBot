package listener

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
	"github.com/sirupsen/logrus"
	"regexp"
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

func handleHelp(b ext.Bot, u *gotgbot.Update) error {
	query := u.CallbackQuery
	pattern, _ := regexp.Compile(`help\((.+?)\)`)

	if pattern.MatchString(query.Data) {
		module := pattern.FindStringSubmatch(query.Data)[1]
		chat := u.EffectiveChat
		dummy := function.GetString(chat.Id, "modules/helpers/help.go:helptxt")
		msg := b.NewSendableEditMessageText(chat.Id, u.EffectiveMessage.MessageId, dummy)
		msg.ParseMode = parsemode.Html
		backButton := [][]ext.InlineKeyboardButton{{ext.InlineKeyboardButton{
			Text:         "Back",
			CallbackData: "help(back)",
		}}}
		backKeyboard := ext.InlineKeyboardMarkup{InlineKeyboard: &backButton}
		msg.ReplyMarkup = &backKeyboard
		if module != "back" {
			msg.Text = function.GetString(chat.Id, "modules/helpers/help.go:"+module)
		} else if module == "back" {
			markup := initHelpButtons()
			msg.ReplyMarkup = &markup
			msg.Text = function.GetString(chat.Id, "modules/helpers/help.go:helptxt")
		}

		_, err := msg.Send()
		err_handler.HandleErr(err)
	}
	return gotgbot.ContinueGroups{}
}

func LoadHelpListener(u *gotgbot.Updater) {
	defer logrus.Info("Help Listener Loaded...")
	u.Dispatcher.AddHandler(handlers.NewCallback("help", handleHelp))
}
