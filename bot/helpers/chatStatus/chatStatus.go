package chatStatus

import (
	"encoding/json"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/go-redis/redis"

	"GoAntispamBot/bot"
	"GoAntispamBot/bot/helpers/errHandler"
	"GoAntispamBot/bot/helpers/trans"
	"GoAntispamBot/bot/providers/redisProvider"
	"GoAntispamBot/bot/providers/telegramProvider"
)

func RequireSudo(userID int, telegramProvider telegramProvider.TelegramProvider) bool {
	if !IsSudo(userID) {
		go telegramProvider.SendText(
			trans.GetString(telegramProvider.Message.Chat.Id, "error/nosudo"),
			0,
			0,
			nil,
		)
		return false
	}
	return true
}

func RequireAdmin(userID int, telegramProvider telegramProvider.TelegramProvider) bool {
	if !isAdmin(userID, telegramProvider.Message.Chat) {
		go telegramProvider.SendText(
			trans.GetString(telegramProvider.Message.Chat.Id, "error/noadmin"),
			0,
			0,
			nil,
		)
		return false
	}
	return true
}

func RequireSuperGroup(chat *ext.Chat, telegramProvider telegramProvider.TelegramProvider) bool {
	if chat.Type != "supergroup" {
		go telegramProvider.SendText(
			trans.GetString(telegramProvider.Message.Chat.Id, "error/nosupergroup"),
			0,
			0,
			nil,
		)
		return false
	}
	return true
}

func RequirePrivate(chat *ext.Chat, telegramProvider telegramProvider.TelegramProvider) bool {
	if chat.Type != "private" {
		go telegramProvider.SendText(
			trans.GetString(telegramProvider.Message.Chat.Id, "error/noprivate"),
			0,
			0,
			nil,
		)
		return false
	}
	return true
}

func IsSudo(userID int) bool {
	return contains(bot.BotConfig.SudoUsers, userID)
}

func isAdmin(userID int, chat *ext.Chat) bool {
	admins := redisProvider.Redis.Get(fmt.Sprintf("admin_%v", chat.Id))
	if admins.Err() != redis.Nil {
		doCreateAdminCache(chat)
	}

	isAdmin := make(map[string][]int)
	_ = admins.Scan(&isAdmin)
	return contains(isAdmin["admins"], userID)
}

func doCreateAdminCache(chat *ext.Chat) {
	listAdmins, _ := chat.GetAdministrators()

	if listAdmins != nil {
		admins := make([]int, 0)
		for _, user := range listAdmins {
			admins = append(admins, user.User.Id)
		}

		cacheAdmin := make(map[string][]int)
		cacheAdmin["admins"] = admins
		finalCache, _ := json.Marshal(&cacheAdmin)
		err := redisProvider.Redis.Set(fmt.Sprintf("admin_%v", chat.Id), finalCache, 600)
		errHandler.Error(err.Err())
	}
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
