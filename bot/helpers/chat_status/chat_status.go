package chat_status

import (
	"encoding/json"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/go-redis/redis"
	"github.com/jumatberkah/antispambot/bot"
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/helpers/function"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
	"strconv"
	"time"
)

type cache struct {
	Adminid []string `json:"admin"`
}

func CanDelete(chat *ext.Chat, botId int) bool {
	k, err := chat.GetMember(botId)
	err_handler.HandleErr(err)
	return k.CanDeleteMessages
}

func IsOwner(userId int) bool {
	for _, user := range bot.BotConfig.SudoUsers {
		if user == strconv.Itoa(userId) {
			return true
		}
	}
	return false
}

func admincache(chat *ext.Chat) {
	adm, _ := chat.GetAdministrators()
	admins := make([]string, 0)
	for _, a := range adm {
		admins = append(admins, strconv.Itoa(a.User.Id))
	}
	one := &cache{admins}
	jsonad, _ := json.Marshal(one)
	sql.REDIS.Set(fmt.Sprintf("admin_%v", chat.Id), jsonad, time.Hour*2)
}

func IsUserAdmin(chat *ext.Chat, user_id int) bool {
	if chat.Type == "private" {
		return true
	}
	if function.Contains(bot.BotConfig.SudoUsers, strconv.Itoa(user_id)) {
		return true
	}

	admins, err := sql.REDIS.Get(fmt.Sprintf("admin_%v", chat.Id)).Result()
	if err != nil {
		if err == redis.Nil {
			admincache(chat)
		}
	}

	var p cache
	_ = json.Unmarshal([]byte(admins), &p)
	if function.Contains(p.Adminid, strconv.Itoa(user_id)) {
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

func RequireBotAdmin(chat *ext.Chat, msg *ext.Message) bool {
	if !IsBotAdmin(chat, nil) {
		_, err := msg.ReplyText("Anda harus menjadikan saya administrator untuk melakukannya.")
		err_handler.HandleErr(err)
		return false
	}
	return true
}

func RequireUserAdmin(chat *ext.Chat, msg *ext.Message, userId int, member *ext.ChatMember) bool {
	if !IsUserAdmin(chat, userId) {
		_, err := msg.ReplyText("Anda harus menjadi administrator untuk melakukannya.")
		err_handler.HandleErr(err)
		return false
	}
	return true
}

func RequirePrivate(chat *ext.Chat, msg *ext.Message) bool {
	if chat.Type != "private" {
		_, err := msg.ReplyText("Penggunaan dibatasi pada pesan privat!")
		err_handler.HandleErr(err)
		return false
	}
	return true
}

func RequireSupergroup(chat *ext.Chat, msg *ext.Message) bool {
	if chat.Type != "supergroup" {
		_, err := msg.ReplyText("Penggunaan dibatasi pada supergrup!")
		err_handler.HandleErr(err)
		return false
	}
	return true
}

func IsUserInChat(chat *ext.Chat, userId int) bool {
	member, err := chat.GetMember(userId)
	err_handler.HandleErr(err)
	if member.Status == "left" || member.Status == "kicked" {
		return false
	} else {
		return true
	}
}

func CanRestrict(bot ext.Bot, chat *ext.Chat) bool {
	botChatMember, err := chat.GetMember(bot.Id)
	err_handler.HandleErr(err)
	if !botChatMember.CanRestrictMembers {
		_, err := bot.SendMessage(chat.Id, "Saya tidak bisa membatasi seseorang. "+
			"Pastikan saya admin dan memiliki semua izin")
		err_handler.HandleErr(err)
		return false
	}
	return true
}
