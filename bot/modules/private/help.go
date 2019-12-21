package private

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
	"fmt"
)

func initHelpButtons() ext.InlineKeyboardMarkup {
	helpButtons := [][]ext.InlineKeyboardButton{make([]ext.InlineKeyboardButton, 2), make([]ext.InlineKeyboardButton, 2),
		make([]ext.InlineKeyboardButton, 1), make([]ext.InlineKeyboardButton, 1), make([]ext.InlineKeyboardButton, 1)}

	// First column
	helpButtons[0][0] = ext.InlineKeyboardButton{
		Text:         "Sudo",
		CallbackData: fmt.Sprintf("help(%v)", "admin"),
	}
	helpButtons[1][0] = ext.InlineKeyboardButton{
		Text:         "Username",
		CallbackData: fmt.Sprintf("help(%v)", "bans"),
	}
	helpButtons[2][0] = ext.InlineKeyboardButton{
		Text:         "Picture",
		CallbackData: fmt.Sprintf("help(%v)", "blacklist"),
	}
	helpButtons[3][0] = ext.InlineKeyboardButton{
		Text:         "Verify",
		CallbackData: fmt.Sprintf("help(%v)", "deleting"),
	}
	helpButtons[4][0] = ext.InlineKeyboardButton{
		Text:         "Notification",
		CallbackData: fmt.Sprintf("help(%v)", "feds"),
	}

	// Second column
	helpButtons[0][1] = ext.InlineKeyboardButton{
		Text:         "Anti Spam",
		CallbackData: fmt.Sprintf("help(%v)", "misc"),
	}
	helpButtons[1][1] = ext.InlineKeyboardButton{
		Text:         "Privacy Policy",
		CallbackData: fmt.Sprintf("help(%v)", "muting"),
	}
	markup := ext.InlineKeyboardMarkup{InlineKeyboard: &helpButtons}
	return markup
}