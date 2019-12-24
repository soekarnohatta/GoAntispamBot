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

func InitHelpButtons() ext.InlineKeyboardMarkup {
	helpButtons := [][]ext.InlineKeyboardButton{make([]ext.InlineKeyboardButton, 2), make([]ext.InlineKeyboardButton, 2),
		make([]ext.InlineKeyboardButton, 2), make([]ext.InlineKeyboardButton, 2)}

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
	helpButtons[3][1] = ext.InlineKeyboardButton{
		Text:         "Misc",
		CallbackData: fmt.Sprintf("help(%v)", "misc"),
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
		replyText := fmt.Sprintf("*%v Version* `3.191223.Stable`\n"+
			"by *Cruzer\n\n*", b.FirstName)

		replyText += function.GetString(chat.Id, "modules/helpers/help.go:helptxt")
		msg := b.NewSendableEditMessageText(chat.Id, u.EffectiveMessage.MessageId, replyText)
		msg.ParseMode = parsemode.Markdown
		backButton := [][]ext.InlineKeyboardButton{{ext.InlineKeyboardButton{
			Text:         "Back",
			CallbackData: "help(back)",
		}}}
		backKeyboard := ext.InlineKeyboardMarkup{InlineKeyboard: &backButton}
		msg.ReplyMarkup = &backKeyboard
		if module != "back" {
			replyTxt := fmt.Sprintf("*%v Version* `%v`\n"+
				"by *Cruzer\n\n*", b.FirstName, bot.BotConfig.BotVer)
			replyTxt += function.GetString(chat.Id, "modules/helpers/help.go:"+module)
			msg.Text = replyTxt
		} else if module == "back" {
			markup := InitHelpButtons()
			msg.ReplyMarkup = &markup
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
