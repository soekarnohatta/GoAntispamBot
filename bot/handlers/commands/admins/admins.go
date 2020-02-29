package admins

import (
	"encoding/json"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jumatberkah/antispambot/bot"
	"github.com/jumatberkah/antispambot/bot/helpers/chat_status"
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/helpers/extraction"
	"github.com/jumatberkah/antispambot/bot/helpers/function"
	"github.com/jumatberkah/antispambot/bot/helpers/telegramProvider"
	"github.com/jumatberkah/antispambot/bot/sql"
)

var (
	banerr = []string{
		"Bad Request: USER_ID_INVALID",
		"Bad Request: USER_NOT_PARTICIPANT",
		"Bad Request: chat member status can't be changed in private chats"}
	_requestProvider = new(telegramProvider.RequestProvider)
)

func gbanUser(b ext.Bot, u *gotgbot.Update, args []string) error {
	_requestProvider.Init(u)
	msg := u.EffectiveMessage
	chat := u.EffectiveChat

	// Permission Check
	if !chat_status.RequireOwner(msg, msg.From.Id) {
		return nil
	}

	userID, reason := extraction.ExtractUserAndText(msg, args)
	if userID == 0 {
		_requestProvider.ReplyHTML(function.GetString(
			chat.Id,
			"handlers/admins/admins.go:27"),
		)
		return nil
	} else if function.Contains(bot.BotConfig.SudoUsers, fmt.Sprint(userID)) || userID == b.Id {
		_requestProvider.ReplyHTML(function.GetString(
			chat.Id,
			"handlers/admins/admins.go:33"),
		)
		return nil
	}

	if reason == "" {
		reason = "No Reason Has Been Specified"
	}

	timeAdd, _ := strconv.Atoi(fmt.Sprint(time.Now().Unix()))
	ban := sql.GetUserSpam(userID)
	if ban != nil {
		if ban.Reason == reason {
			_requestProvider.ReplyHTML(function.GetString(
				chat.Id,
				"handlers/admins/admins.go:38"),
			)
			return nil
		}

		_requestProvider.ReplyHTML(function.GetStringf(
			chat.Id,
			"handlers/admins/admins.go:43",
			map[string]string{
				"1": strconv.Itoa(userID),
				"2": ban.Reason,
				"3": reason}),
		)

		sql.UpdateUserSpam(
			userID,
			reason,
			fmt.Sprint(msg.From.Id),
			timeAdd,
		)

		err := function.SendBanLog(b, userID, reason, u)
		err_handler.HandleErr(err)
		return err
	}

	_requestProvider.ReplyHTML(function.GetStringf(
		chat.Id,
		"handlers/admins/admins.go:54",
		map[string]string{"1": strconv.Itoa(userID)}),
	)

	sql.UpdateUserSpam(
		userID,
		reason,
		fmt.Sprint(msg.From.Id),
		timeAdd,
	)

	_requestProvider.ReplyHTML(function.GetStringf(
		chat.Id,
		"handlers/admins/admins.go:62",
		map[string]string{"1": strconv.Itoa(userID), "2": reason}),
	)

	err := function.SendBanLog(b, userID, reason, u)
	err_handler.HandleErr(err)
	return err
}

func unGbanUser(b ext.Bot, u *gotgbot.Update, args []string) error {
	_requestProvider.Init(u)
	msg := u.EffectiveMessage
	chat := u.EffectiveChat

	// Permission Check
	if !chat_status.RequireOwner(msg, msg.From.Id) {
		return nil
	}

	userID, _ := extraction.ExtractUserAndText(msg, args)
	if userID == 0 {
		_requestProvider.ReplyHTML(function.GetString(chat.Id,
			"handlers/admins/admins.go:27"),
		)
		return nil
	} else if function.Contains(bot.BotConfig.SudoUsers, fmt.Sprint(userID)) || userID == b.Id {
		_requestProvider.ReplyHTML(function.GetString(chat.Id,
			"handlers/admins/admins.go:33"),
		)
		return nil
	}

	ban := sql.GetUserSpam(userID)
	if ban != nil {
		_requestProvider.ReplyHTML(function.GetStringf(
			chat.Id,
			"handlers/admins/admins.go:88",
			map[string]string{"1": strconv.Itoa(userID)}),
		)

		go func() {
			group := sql.GetAllChat
			sql.DelUserSpam(userID)

			for _, a := range group() {
				cID, _ := strconv.Atoi(a.ChatId)
				_, err := b.UnbanChatMember(cID, userID)
				if err != nil {
					if function.Contains(banerr, err.Error()) == true {
						err_handler.HandleErr(err)
						continue
					} else if err.Error() == "Forbidden: bot was kicked from the supergroup chat" {
						sql.DelChat(a.ChatId)
						continue
					}
				}
			}
		}()

		_requestProvider.ReplyHTML(function.GetStringf(
			chat.Id,
			"handlers/admins/admins.go:111",
			map[string]string{"1": strconv.Itoa(userID)}),
		)
		return nil
	}

	_requestProvider.ReplyHTML(function.GetString(
		chat.Id,
		"handlers/admins/admins.go:116"),
	)
	return nil
}

func stats(_ ext.Bot, u *gotgbot.Update) error {
	_requestProvider.Init(u)
	msg := u.EffectiveMessage

	// Permission Check
	if !chat_status.RequireOwner(msg, msg.From.Id) {
		return nil
	}

	replyText := fmt.Sprintf("*Statistics*"+
		"\nTotal User(s): `%v`"+
		"\nTotal Chat(s): `%v`"+
		"\nTotal Spammer(s): `%v`",
		len(sql.GetAllUser()),
		len(sql.GetAllChat()),
		len(sql.GetAllSpamUser()),
	)

	_requestProvider.ReplyMarkdown(replyText)
	return nil
}

func unban(b ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat
	user := u.EffectiveUser
	message := u.EffectiveMessage

	// Permission checks
	if !chat_status.RequireOwner(message, message.From.Id) {
		return nil
	}

	if !chat_status.RequireSupergroup(chat, message) {
		return nil
	}

	if !chat_status.IsBotAdmin(chat, nil) && chat_status.RequireUserAdmin(chat, message, user.Id) {
		return gotgbot.EndGroups{}
	}

	userID := 0
	cID := 0
	pattern, _ := regexp.Compile(`-100\d{10}`)
	if pattern.MatchString(message.Text) {
		cID, _ = strconv.Atoi(pattern.String())
	}

	userID, _ = extraction.ExtractUserAndText(message, args)

	if userID == 0 {
		_, err := message.ReplyText("Try targeting a user next time bud.")
		return err
	}

	_, err := b.GetChatMember(cID, userID)
	if err != nil {
		_, err := message.ReplyText("This user is ded m8.")
		return err
	}

	userMember, _ := b.GetChatMember(cID, userID)
	if !userMember.CanRestrictMembers && userMember.Status != "creator" {
		_, err = message.ReplyText("You don't have permissions to unban users!")
		return err
	}

	if userID == b.Id {
		_, err := message.ReplyText("What exactly are you attempting to do?.")
		return err
	}

	_, err = chat.UnbanMember(userID)
	err_handler.HandleErr(err)
	_, err = message.ReplyText("Fine, I'll allow it, this time...")
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
			cID, _ := strconv.Atoi(a.ChatId)
			_, err := b.SendMessageHTML(cID, txtToSend)
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
	_requestProvider.Init(u)
	msg := u.EffectiveMessage

	// Permission Check
	if !chat_status.RequireOwner(msg, msg.From.Id) {
		return nil
	}

	if msg.ReplyToMessage != nil {
		//jsonData, err := json.Marshal(msg.ReplyToMessage)
		//err_handler.HandleErr(err)
		output, err := json.MarshalIndent(msg.ReplyToMessage, "", "  ")
		err_handler.HandleErr(err)
		_requestProvider.SendTextAsync(
			string(output),
			0,
			0,
			"",
			nil,
		)
		return nil
	}

	//jsonData, err := json.Marshal(msg)
	//err_handler.HandleErr(err)
	output, err := json.MarshalIndent(msg, "", "  ")
	err_handler.HandleErr(err)
	_requestProvider.SendTextAsync(
		string(output),
		0,
		0,
		"",
		nil,
	)
	return nil
}

func ping(_ ext.Bot, u *gotgbot.Update) error {
	_requestProvider.Init(u)
	_requestProvider.ReplyMarkdown("*Pong...*")
	return nil
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
