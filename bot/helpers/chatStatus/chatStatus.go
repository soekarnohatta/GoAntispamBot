package chatStatus

import (
	"github.com/PaulSonOfLars/gotgbot/ext"

	"GoAntispamBot/bot"
)

func RequireSudo(userID int) {
	if !isSudo(userID) {
		// do something...
	}
}

func RequireAdmin(userID int) {
	if !isAdmin(userID) {
		// do something...
	}
}

func RequireSuperGroup(chat *ext.Chat) {
	if chat.Type != "supergroup" {

	}
}

func RequirePrivate(chat *ext.Chat) {
	if chat.Type != "private" {

	}
}

func isSudo(userID int) bool {
	return contains(bot.BotConfig.SudoUsers, userID)
}

func isAdmin(userID int) bool {
	return true
}

func contains(key []int, val int) bool {
	if key != nil && val != 0 {
		for _, res := range key {
			if res == val {
				return true
			}
		}
	}
	return false
}
