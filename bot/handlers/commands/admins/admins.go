package admins

import (
	"encoding/json"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"

	"github.com/jumatberkah/antispambot/bot"
	"github.com/jumatberkah/antispambot/bot/helpers/chat_status"
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/helpers/extraction"
	"github.com/jumatberkah/antispambot/bot/helpers/function"
	"github.com/jumatberkah/antispambot/bot/helpers/logger"
	"github.com/jumatberkah/antispambot/bot/sql"
)

var banerr = []string{
	"Bad Request: USER_ID_INVALID",
	"Bad Request: USER_NOT_PARTICIPANT",
	"Bad Request: chat member status can't be changed in private chats"}

func gbanUser(b ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage

	if !chat_status.RequireOwner(msg, msg.From.Id) {
		return nil
	}

	userId, reason := extraction.ExtractUserAndText(msg, args)
	if userId == 0 {
		_, err := b.SendMessageHTML(
			msg.Chat.Id,
			function.GetString(msg.Chat.Id, "handlers/admins/admins.go:27"),
		)
		return err
	} else if function.Contains(bot.BotConfig.SudoUsers, fmt.Sprint(userId)) || userId == b.Id {
		_, err := b.SendMessageHTML(
			msg.Chat.Id,
			function.GetString(msg.Chat.Id, "handlers/admins/admins.go:33"),
		)
		return err
	}

	if reason == "" {
		reason = "No Reason Has Been Specified"
	}

	timeAdd, _ := strconv.Atoi(fmt.Sprint(time.Now().Unix()))
	ban := sql.GetUserSpam(userId)
	if ban != nil {
		if ban.Reason == reason {
			_, err := b.SendMessageHTML(
				msg.Chat.Id,
				function.GetString(msg.Chat.Id, "handlers/admins/admins.go:38"),
			)
			return err
		}

		_, err := b.SendMessageHTML(
			msg.Chat.Id,
			function.GetStringf(
				msg.Chat.Id,
				"handlers/admins/admins.go:43",
				map[string]string{
					"1": strconv.Itoa(userId),
					"2": ban.Reason,
					"3": reason},
			),
		)

		err_handler.HandleErr(err)
		sql.UpdateUserSpam(
			userId,
			reason,
			fmt.Sprint(msg.From.Id),
			timeAdd,
		)
		err = logger.SendBanLog(b, userId, reason, u)
		return err
	}

	_, err := b.SendMessageHTML(
		msg.Chat.Id,
		function.GetStringf(
			msg.Chat.Id,
			"handlers/admins/admins.go:54",
			map[string]string{"1": strconv.Itoa(userId)},
		),
	)

	err_handler.HandleErr(err)
	sql.UpdateUserSpam(
		userId,
		reason,
		fmt.Sprint(msg.From.Id),
		timeAdd,
	)

	_, err = b.SendMessageHTML(
		msg.Chat.Id,
		function.GetStringf(
			msg.Chat.Id,
			"handlers/admins/admins.go:62",
			map[string]string{"1": strconv.Itoa(userId), "2": reason}),
	)
	err_handler.HandleErr(err)

	err = logger.SendBanLog(b, userId, reason, u)
	return err
}

func unGbanUser(b ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage

	if !chat_status.RequireOwner(msg, msg.From.Id) {
		return nil
	}

	userId, _ := extraction.ExtractUserAndText(msg, args)

	if userId == 0 {
		_, err := msg.ReplyHTML(function.GetString(msg.Chat.Id, "handlers/admins/admins.go:27"))
		return err
	} else if function.Contains(bot.BotConfig.SudoUsers, fmt.Sprint(userId)) || userId == b.Id {
		_, err := msg.ReplyHTML(function.GetString(msg.Chat.Id, "handlers/admins/admins.go:33"))
		return err
	}

	ban := sql.GetUserSpam(userId)
	if ban != nil {
		_, err := msg.ReplyHTMLf(function.GetStringf(msg.Chat.Id, "handlers/admins/admins.go:88",
			map[string]string{"1": strconv.Itoa(userId)}))
		err_handler.HandleErr(err)

		go func() {
			group := sql.GetAllChat
			sql.DelUserSpam(userId)

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

		_, err = msg.ReplyHTML(
			function.GetStringf(
				msg.Chat.Id,
				"handlers/admins/admins.go:111",
				map[string]string{"1": strconv.Itoa(userId)},
			),
		)
		return err
	}
	_, err := msg.ReplyHTML(function.GetString(msg.Chat.Id, "handlers/admins/admins.go:116"))
	return err
}

func stats(_ ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage

	if !chat_status.RequireOwner(msg, msg.From.Id) {
		return nil
	}

	replyText := "*Statistics*" +
		"\nTotal User(s): `%v`" +
		"\nTotal Chat(s): `%v`" +
		"\nTotal Spammer(s): `%v`"

	_, err := msg.ReplyMarkdownf(
		replyText,
		len(sql.GetAllUser()),
		len(sql.GetAllChat()),
		len(sql.GetAllSpamUser()),
	)
	return err
}

func broadcast(b ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage

	if !chat_status.RequireOwner(msg, msg.From.Id) {
		return nil
	}

	group := sql.GetAllChat
	errNum := 0
	txtToSend := ""

	if msg.ReplyToMessage != nil {
		txtToSend = msg.ReplyToMessage.Text
	} else {
		txtToSend = strings.Split(msg.OriginalHTML(), "/broadcast")[1]
	}

	if txtToSend != "" {
		for _, a := range group() {
			cid, _ := strconv.Atoi(a.ChatId)
			_, err := b.SendMessageHTML(cid, txtToSend)
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
	} else {
		_, err := msg.ReplyHTML("<b>You must specify a message in order to broadcast something!</b>")
		return err
	}

	_, err := msg.ReplyHTMLf(
		"<b>Message Has Been Broadcasted</b>,"+
			"<code>%v</code> <b>Has Failed</b>\n",
		errNum,
	)
	return err
}

func dbg(b ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage
	if !chat_status.RequireOwner(msg, msg.From.Id) {
		return nil
	}

	if msg.ReplyToMessage != nil {
		jsonData, err := json.Marshal(msg.ReplyToMessage)
		err_handler.HandleErr(err)
		_, err = msg.ReplyText(string(jsonData))
		return err
	} else {
		jsonData, err := json.Marshal(msg)
		err_handler.HandleErr(err)
		_, err = msg.ReplyText(string(jsonData))
		return err
	}
}

func ping(_ ext.Bot, u *gotgbot.Update) error {
	_, err := u.EffectiveMessage.ReplyMarkdown("*Pong*")
	return err
}

func LoadAdmins(u *gotgbot.Updater) {
	defer logrus.Info("Admins Module Loaded...")
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("gban", []rune{'/', '.'}, gbanUser))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("ungban", []rune{'/', '.'}, unGbanUser))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("stats", []rune{'/', '.'}, stats))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("broadcast", []rune{'/', '.'}, broadcast))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("ping", []rune{'/', '.'}, ping))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("dbg", []rune{'/', '.'}, dbg))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("backup", []rune{'/', '.'}, backupDb))
}
