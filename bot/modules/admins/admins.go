package admins

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/chat_status"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/extraction"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/logger"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func gban(b ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage

	if chat_status.RequireOwner(msg, msg.From.Id) == false {
		return gotgbot.EndGroups{}
	}

	userid, reason := extraction.ExtractUserAndText(msg, args)
	if userid == 0 {
		_, err := b.SendMessageHTML(msg.Chat.Id, function.GetString(msg.Chat.Id, "modules/admins/admins.go:27"))
		err_handler.HandleErr(err)
		return err
	}

	if reason == "" {
		reason = "None"
	}

	ban := sql.GetUserSpam(userid)
	if ban != nil {
		if ban.Reason == reason {
			_, err := b.SendMessageHTML(msg.Chat.Id, function.GetString(msg.Chat.Id, "modules/admins/admins.go:38"))
			err_handler.HandleErr(err)
			return err
		}

		_, err := b.SendMessageHTML(msg.Chat.Id, function.GetStringf(msg.Chat.Id, "modules/admins/admins.go:43",
			map[string]string{"1": strconv.Itoa(userid), "2": ban.Reason, "3": reason}))
		err_handler.HandleErr(err)
		err = sql.UpdateUserSpam(userid, reason)
		err_handler.HandleTgErr(b, u, err)
		err = logger.SendBanLog(b, userid, reason, u)
		err_handler.HandleErr(err)
		return err
	}

	_, err := b.SendMessageHTML(msg.Chat.Id, function.GetStringf(msg.Chat.Id, "modules/admins/admins.go:54",
		map[string]string{"1": strconv.Itoa(userid)}))
	err_handler.HandleErr(err)

	err = sql.UpdateUserSpam(userid, reason)
	err_handler.HandleTgErr(b, u, err)

	_, err = b.SendMessageHTML(msg.Chat.Id, function.GetStringf(msg.Chat.Id, "modules/admins/admins.go:62",
		map[string]string{"1": strconv.Itoa(userid), "2": reason}))
	err_handler.HandleErr(err)
	err = logger.SendBanLog(b, userid, reason, u)
	err_handler.HandleErr(err)
	return err
}

func ungban(b ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage

	if chat_status.RequireOwner(msg, msg.From.Id) == false {
		return nil
	}

	userid, _ := extraction.ExtractUserAndText(msg, args)
	if userid == 0 {
		_, err := msg.ReplyHTML(function.GetString(msg.Chat.Id, "modules/admins/admins.go:27"))
		return err
	}

	ban := sql.GetUserSpam(userid)
	if ban != nil {
		_, err := msg.ReplyHTMLf(function.GetStringf(msg.Chat.Id, "modules/admins/admins.go:88",
			map[string]string{"1": strconv.Itoa(userid)}))
		err_handler.HandleErr(err)

		go func() {
			group := sql.GetAllChat
			banerr := []string{"Bad Request: USER_ID_INVALID", "Bad Request: USER_NOT_PARTICIPANT" +
				"Bad Request: chat member status can't be changed in private chats"}
			sql.DelUserSpam(userid)
			for _, a := range group() {
				cid, _ := strconv.Atoi(a.ChatId)
				_, err = b.UnbanChatMember(cid, userid)
				if err != nil {
					if function.Contains(banerr, err.Error()) == true {
						return
					} else if err.Error() == "Forbidden: bot was kicked from the supergroup chat" {
						sql.DelChat(a.ChatId)
						return
					}
				}
			}
		}()

		_, err = msg.ReplyHTMLf(function.GetStringf(msg.Chat.Id, "modules/admins/admins.go:111",
			map[string]string{"1": strconv.Itoa(userid)}))
		err_handler.HandleErr(err)
		return err
	}
	_, err := msg.ReplyHTML(function.GetString(msg.Chat.Id, "modules/admins/admins.go:116"))
	err_handler.HandleErr(err)
	return err
}

func stats(_ ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage

	if chat_status.RequireOwner(msg, msg.From.Id) == false {
		return gotgbot.EndGroups{}
	}

	teks := fmt.Sprintf("<b>Statistics</b>\n"+
		"Total Users: %v\nTotal Chats: %v\nTotal Spammers: %v", len(sql.GetAllUser()),
		len(sql.GetAllChat()), len(sql.GetAllSpamUser()))

	_, err := msg.ReplyHTML(teks)
	err_handler.HandleErr(err)
	return err
}

func broadcast(b ext.Bot, u *gotgbot.Update) error {
	var err error
	msg := u.EffectiveMessage

	if chat_status.RequireOwner(msg, msg.From.Id) == false {
		return gotgbot.EndGroups{}
	}

	group := sql.GetAllChat
	errnum := 0

	for _, a := range group() {
		cid, _ := strconv.Atoi(a.ChatId)
		_, err = b.SendMessageHTML(cid, strings.Split(msg.Text, "/broadcast")[1])
		if err != nil {
			if err.Error() == "Forbidden: bot was kicked from the supergroup chat" {
				sql.DelChat(a.ChatId)
				errnum++
			} else {
				err_handler.HandleErr(err)
				errnum++
			}
		}
	}

	_, err = msg.ReplyHTMLf("<b>Message Has Been Broadcasted</b>, <code>%v</code> <b>Has Failed</b>\n", errnum)
	err_handler.HandleErr(err)
	return err
}

func LoadAdmins(u *gotgbot.Updater) {
	defer logrus.Info("Admins Module Loaded...")
	go u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("gban", []rune{'/', '.'}, gban))
	go u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("ungban", []rune{'/', '.'}, ungban))
	go u.Dispatcher.AddHandler(handlers.NewPrefixCommand("stats", []rune{'/', '.'}, stats))
	go u.Dispatcher.AddHandler(handlers.NewPrefixCommand("broadcast", []rune{'/', '.'}, broadcast))
}
