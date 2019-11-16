package chat_status

import (
	"encoding/json"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/jumatberkah/antispambot/bot"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/caching"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
	"strconv"
)

type cache struct {
	Adminid []string `json:"admin"`
}

func IsOwner(userId int) bool {
	for _, user := range bot.BotConfig.SudoUsers {
		if user == strconv.Itoa(userId) {
			return true
		}
	}
	return false
}

func IsUserAdmin(chat *ext.Chat, user_id int) bool {
	if chat.Type == "private" {
		return true
	}
	if IsOwner(user_id) {
		return true
	}

	admins, err := caching.CACHE.Get(fmt.Sprintf("admin_%v", chat.Id))
	if err != nil {
		admincache(chat)
	}

	var x cache
	_ = json.Unmarshal(admins, &x)
	if function.Contains(x.Adminid, strconv.Itoa(user_id)) {
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
		_, err := msg.ReplyText(function.GetString(chat.Id, "modules/helpers/chat_status.go:73"))
		err_handler.HandleErr(err)
		return false
	}
	return true
}

func RequireOwner(msg *ext.Message, userId int) bool {
	if !IsOwner(userId) {
		_, err := msg.ReplyHTML(function.GetString(msg.Chat.Id, "modules/helpers/chat_status.go:82"))
		err_handler.HandleErr(err)
		return false
	}
	return true
}

func RequirePrivate(chat *ext.Chat, msg *ext.Message) bool {
	if chat.Type != "private" {
		_, err := msg.ReplyText(function.GetString(msg.Chat.Id, "modules/helpers/chat_status.go:91"))
		err_handler.HandleErr(err)
		return false
	}
	return true
}

func RequireSupergroup(chat *ext.Chat, msg *ext.Message) bool {
	if chat.Type != "supergroup" {
		_, err := msg.ReplyText(function.GetString(msg.Chat.Id, "modules/helpers/chat_status.go:100"))
		err_handler.HandleErr(err)
		return false
	}
	return true
}

func CanRestrict(bot ext.Bot, chat *ext.Chat) bool {
	botChatMember, err := chat.GetMember(bot.Id)
	err_handler.HandleErr(err)
	if !botChatMember.CanRestrictMembers {
		_, err := bot.SendMessage(chat.Id, function.GetString(chat.Id, "modules/helpers/chat_status.go:111"))
		err_handler.HandleErr(err)
		return false
	}
	return true
}

func CanDelete(bot ext.Bot, chat *ext.Chat) bool {
	botChatMember, err := chat.GetMember(bot.Id)
	err_handler.HandleErr(err)
	if !botChatMember.CanDeleteMessages {
		_, err := bot.SendMessage(chat.Id, function.GetString(chat.Id, "modules/helpers/chat_status.go:122"))
		err_handler.HandleErr(err)
		return false
	}
	return true
}

func admincache(chat *ext.Chat) {
	x, _ := chat.GetAdministrators()
	admins := make([]string, 0)

	for _, y := range x {
		admins = append(admins, strconv.Itoa(y.User.Id))
	}

	w := &cache{admins}
	z, _ := json.Marshal(w)
	_ = caching.CACHE.Set(fmt.Sprintf("admin_%v", chat.Id), z)
}