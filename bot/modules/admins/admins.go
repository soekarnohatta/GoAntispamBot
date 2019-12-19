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
	"html"
	"strconv"
	"strings"
)

func gbanUser(b ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage

	if chat_status.RequireOwner(msg, msg.From.Id) == false {
		return nil
	}

	userId, reason := extraction.ExtractUserAndText(msg, args)
	if userId == 0 {
		_, err := b.SendMessageHTML(msg.Chat.Id, function.GetString(msg.Chat.Id, "modules/admins/admins.go:27"))
		return err
	}

	if reason == "" {
		reason = "No Reason Has Been Specified"
	}

	ban := sql.GetUserSpam(userId)
	if ban != nil {
		if ban.Reason == reason {
			_, err := b.SendMessageHTML(msg.Chat.Id, function.GetString(msg.Chat.Id, "modules/admins/admins.go:38"))
			return err
		}

		_, err := b.SendMessageHTML(msg.Chat.Id, function.GetStringf(msg.Chat.Id, "modules/admins/admins.go:43",
			map[string]string{"1": strconv.Itoa(userId), "2": ban.Reason, "3": reason}))
		err_handler.HandleErr(err)
		go sql.UpdateUserSpam(userId, reason)
		err = logger.SendBanLog(b, userId, reason, u)
		return err
	}

	_, err := b.SendMessageHTML(msg.Chat.Id, function.GetStringf(msg.Chat.Id, "modules/admins/admins.go:54",
		map[string]string{"1": strconv.Itoa(userId)}))
	err_handler.HandleErr(err)
	go sql.UpdateUserSpam(userId, reason)
	_, err = b.SendMessageHTML(msg.Chat.Id, function.GetStringf(msg.Chat.Id, "modules/admins/admins.go:62",
		map[string]string{"1": strconv.Itoa(userId), "2": reason}))
	err_handler.HandleErr(err)
	err = logger.SendBanLog(b, userId, reason, u)
	return err
}

func unGbanUser(b ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage

	if chat_status.RequireOwner(msg, msg.From.Id) == false {
		return nil
	}

	userId, _ := extraction.ExtractUserAndText(msg, args)

	if userId == 0 {
		_, err := msg.ReplyHTML(function.GetString(msg.Chat.Id, "modules/admins/admins.go:27"))
		return err
	}

	ban := sql.GetUserSpam(userId)

	if ban != nil {
		_, err := msg.ReplyHTMLf(function.GetStringf(msg.Chat.Id, "modules/admins/admins.go:88",
			map[string]string{"1": strconv.Itoa(userId)}))
		err_handler.HandleErr(err)

		go func() {
			group := sql.GetAllChat
			banerr := []string{"Bad Request: USER_ID_INVALID", "Bad Request: USER_NOT_PARTICIPANT" +
				"Bad Request: chat member status can't be changed in private chats"}
			go sql.DelUserSpam(userId)
			for _, a := range group() {
				cid, _ := strconv.Atoi(a.ChatId)
				_, err = b.UnbanChatMember(cid, userId)
				if err != nil {
					if function.Contains(banerr, err.Error()) == true {
						continue
					} else if err.Error() == "Forbidden: bot was kicked from the supergroup chat" {
						sql.DelChat(a.ChatId)
						continue
					}
				}
			}
		}()

		_, err = msg.ReplyHTMLf(function.GetStringf(msg.Chat.Id, "modules/admins/admins.go:111",
			map[string]string{"1": strconv.Itoa(userId)}))
		return err
	}
	_, err := msg.ReplyHTML(function.GetString(msg.Chat.Id, "modules/admins/admins.go:116"))
	return err
}

func stats(_ ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage

	if chat_status.RequireOwner(msg, msg.From.Id) == false {
		return nil
	}

	replyText := fmt.Sprintf(
		"<b>Statistics</b>\n"+
			"Total User(s): %v\n"+
			"Total Chat(s): %v\n"+
			"Total Spammer(s): %v", len(sql.GetAllUser()),
		len(sql.GetAllChat()), len(sql.GetAllSpamUser()))

	_, err := msg.ReplyHTML(replyText)
	return err
}

func broadcast(b ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage

	if chat_status.RequireOwner(msg, msg.From.Id) == false {
		return nil
	}

	group := sql.GetAllChat
	errNum := 0

	for _, a := range group() {
		cid, _ := strconv.Atoi(a.ChatId)
		_, err := b.SendMessageHTML(cid, html.EscapeString(strings.Split(msg.Text, "/broadcast")[1]))
		if err != nil {
			if err.Error() == "Forbidden: bot was kicked from the supergroup chat" {
				sql.DelChat(a.ChatId)
				errNum++
				continue
			} else {
				err_handler.HandleErr(err)
				errNum++
				continue
			}
		}
	}

	_, err := msg.ReplyHTMLf("<b>Message Has Been Broadcasted</b>,"+
		"<code>%v</code> <b>Has Failed</b>\n", errNum)
	return err
}

func LoadAdmins(u *gotgbot.Updater) {
	defer logrus.Info("Admins Module Loaded...")
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("gban", []rune{'/', '.'}, gbanUser))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("ungban", []rune{'/', '.'}, unGbanUser))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("stats", []rune{'/', '.'}, stats))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("broadcast", []rune{'/', '.'}, broadcast))
}
