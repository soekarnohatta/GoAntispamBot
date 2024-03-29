package chat_status

import (
	"encoding/json"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"strconv"

	"github.com/jumatberkah/antispambot/bot"
	"github.com/jumatberkah/antispambot/bot/helpers/caching"
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/helpers/function"
)

type adminCache struct {
	Admin []string `json:"admin"`
}

type restrictCache struct {
	Restrict bool
}

type deleteCache struct {
	Delete bool
}

func IsOwner(userId int) bool {
	if function.Contains(bot.BotConfig.SudoUsers, strconv.Itoa(userId)) {
		return true
	}
	return false
}

func IsUserAdmin(chat *ext.Chat, userId int) bool {
	if chat.Type == "private" {
		return true
	}

	if IsOwner(userId) {
		return true
	}

	admins, err := caching.CACHE.Get(fmt.Sprintf("admin_%v", chat.Id))

	if err != nil {
		AdminCache(chat)
		admins, err = caching.CACHE.Get(fmt.Sprintf("admin_%v", chat.Id))
		var aCache adminCache
		_ = json.Unmarshal(admins, &aCache)
		if function.Contains(aCache.Admin, strconv.Itoa(userId)) {
			return true
		}
		return false
	}

	var aCache adminCache
	_ = json.Unmarshal(admins, &aCache)
	if function.Contains(aCache.Admin, strconv.Itoa(userId)) {
		return true
	}
	return false
}

func IsBotAdmin(chat *ext.Chat, member *ext.ChatMember) bool {
	if chat.Type == "private" {
		return true
	}

	if member == nil {
		mem, err := chat.GetMember(chat.Bot.Id)
		err_handler.HandleErr(err)
		if mem == nil {
			return false
		}
		member = mem
	}

	if member.Status == "administrator" || member.Status == "creator" {
		return true
	} else {
		return false
	}
}

func RequireUserAdmin(chat *ext.Chat, msg *ext.Message, userId int) bool {
	if !IsUserAdmin(chat, userId) {
		_, err := msg.ReplyHTML(function.GetString(chat.Id, "handlers/helpers/chat_status.go:73"))
		err_handler.HandleErr(err)
		return false
	}
	return true
}

func RequireOwner(msg *ext.Message, userId int) bool {
	if !IsOwner(userId) {
		_, err := msg.ReplyHTML(function.GetString(msg.Chat.Id, "handlers/helpers/chat_status.go:82"))
		err_handler.HandleErr(err)
		return false
	}
	return true
}

func RequirePrivate(chat *ext.Chat, msg *ext.Message) bool {
	if chat.Type != "private" {
		_, err := msg.ReplyHTML(function.GetString(msg.Chat.Id, "handlers/helpers/chat_status.go:91"))
		err_handler.HandleErr(err)
		return false
	}
	return true
}

func RequireSupergroup(chat *ext.Chat, msg *ext.Message) bool {
	if chat.Type != "supergroup" {
		_, err := msg.ReplyHTML(function.GetString(msg.Chat.Id, "handlers/helpers/chat_status.go:100"))
		err_handler.HandleErr(err)
		return false
	}
	return true
}

func CanRestrict(bot ext.Bot, chat *ext.Chat) bool {
	restrict, err := caching.CACHE.Get(fmt.Sprintf("restrict_%v", chat.Id))
	if err != nil {
		doRestrictCache(chat, bot)
		botChatMember, err := chat.GetMember(bot.Id)
		err_handler.HandleErr(err)
		if botChatMember != nil && !botChatMember.CanRestrictMembers {
			_, err := bot.SendMessage(chat.Id, function.GetString(chat.Id, "handlers/helpers/chat_status.go:111"))
			err_handler.HandleErr(err)
			return false
		}
		return true
	}

	var rCache restrictCache
	_ = json.Unmarshal(restrict, &rCache)
	return rCache.Restrict
}

func CanDelete(bot ext.Bot, chat *ext.Chat) bool {
	restrict, err := caching.CACHE.Get(fmt.Sprintf("delete_%v", chat.Id))

	if err != nil {
		doDeleteCache(chat, bot)
		botChatMember, err := chat.GetMember(bot.Id)
		err_handler.HandleErr(err)
		if !botChatMember.CanDeleteMessages {
			_, err := bot.SendMessage(chat.Id, function.GetString(chat.Id, "handlers/helpers/chat_status.go:122"))
			err_handler.HandleErr(err)
			return false
		}
		return true
	}

	var rCache deleteCache
	_ = json.Unmarshal(restrict, &rCache)
	return rCache.Delete
}

func AdminCache(chat *ext.Chat) {
	listAdmins, _ := chat.GetAdministrators()
	admins := make([]string, 0)

	for _, user := range listAdmins {
		admins = append(admins, strconv.Itoa(user.User.Id))
	}

	cacheAdmin := &adminCache{admins}
	finalCache, _ := json.Marshal(cacheAdmin)
	_ = caching.CACHE.Set(fmt.Sprintf("admin_%v", chat.Id), finalCache)
}

func doDeleteCache(chat *ext.Chat, bot ext.Bot) {
	botChatMember, err := chat.GetMember(bot.Id)
	err_handler.HandleErr(err)
	deleteCaches := false
	if botChatMember != nil && botChatMember.CanDeleteMessages {
		deleteCaches = true
	}

	cacheDelete := &deleteCache{deleteCaches}
	finalCache, _ := json.Marshal(cacheDelete)
	_ = caching.CACHE.Set(fmt.Sprintf("delete_%v", chat.Id), finalCache)
}

func doRestrictCache(chat *ext.Chat, bot ext.Bot) {
	botChatMember, err := chat.GetMember(bot.Id)
	err_handler.HandleErr(err)
	deleteRestrict := false
	if botChatMember != nil && botChatMember.CanRestrictMembers {
		deleteRestrict = true
	}

	cacheDelete := &restrictCache{deleteRestrict}
	finalCache, _ := json.Marshal(cacheDelete)
	_ = caching.CACHE.Set(fmt.Sprintf("restrict_%v", chat.Id), finalCache)
}
