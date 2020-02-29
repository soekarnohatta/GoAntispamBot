package listener

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/handlers/Filters"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"github.com/sirupsen/logrus"
	"html"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jumatberkah/antispambot/bot/helpers/caching"
	"github.com/jumatberkah/antispambot/bot/helpers/chat_status"
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/helpers/extraction"
	"github.com/jumatberkah/antispambot/bot/helpers/function"
	"github.com/jumatberkah/antispambot/bot/sql"
)

func usernameScan(b ext.Bot, u *gotgbot.Update) error {
	db := sql.GetUsername(u.EffectiveChat.Id)

	if db == nil {
		return nil
	}

	if db.Option != "true" {
		return gotgbot.ContinueGroups{}
	}

	user := u.EffectiveUser
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	if chat.Type != "supergroup" {
		return nil
	}

	if chat_status.IsUserAdmin(chat, user.Id) {
		return nil
	}

	if !chat_status.CanRestrict(b, chat) {
		return nil
	}

	if user.Username != "" {
		return gotgbot.ContinueGroups{}
	}

	replyText := function.GetStringf(
		chat.Id,
		"handlers/listener/listener.go:45",
		map[string]string{
			"1": strconv.Itoa(user.Id),
			"2": html.EscapeString(user.FirstName),
			"3": db.Action,
			"4": strconv.Itoa(user.Id)},
	)
	markup := &ext.InlineKeyboardMarkup{}

	if db.Action != "warn" {
		reply := b.NewSendableMessage(chat.Id, replyText)
		reply.ParseMode = parsemode.Html
		reply.ReplyToMessageId = msg.MessageId

		switch db.Action {
		case "mute":
			restrictSend := b.NewSendableRestrictChatMember(chat.Id, user.Id)
			restrictSend.UntilDate = extraction.ExtractTime(
				b,
				msg,
				sql.GetSetting(chat.Id).Time,
			)
			_, err := restrictSend.Send()
			err_handler.HandleErr(err)

			keyb := function.BuildKeyboardf(
				"data/keyboard/usermonitor.json",
				1,
				map[string]string{
					"1": "Cara Pasang Username",
					"2": function.GetString(chat.Id, "handlers/listener/listener.go:51"),
					"3": fmt.Sprintf("umute_%v_%v", user.Id, chat.Id),
					"4": "https://www.telegra.ph/Cara-Membuat-Username-di-Klien-Telegram-01-28",
				})
			markup = &ext.InlineKeyboardMarkup{InlineKeyboard: &keyb}
		case "kick":
			_, err := b.UnbanChatMember(chat.Id, user.Id)
			err_handler.HandleErr(err)
			markup = nil
		case "ban":
			restrictSend := b.NewSendableKickChatMember(chat.Id, user.Id)
			restrictSend.UntilDate = -1
			_, err := restrictSend.Send()
			err_handler.HandleErr(err)

			keyb := function.BuildKeyboardf(
				"data/keyboard/usermonitor.json",
				1,
				map[string]string{
					"1": "Cara Pasang Username",
					"2": function.GetString(chat.Id, "handlers/listener/listener.go:56"),
					"3": fmt.Sprintf("uba_%v_%v", user.Id, chat.Id),
					"4": "https://www.telegra.ph/Cara-Membuat-Username-di-Klien-Telegram-01-28",
				})
			markup = &ext.InlineKeyboardMarkup{InlineKeyboard: &keyb}
		}

		reply.ReplyMarkup = markup
		_, err := reply.Send()
		if err != nil {
			if err.Error() == "Bad Request: reply message not found" {
				reply.ReplyToMessageId = 0
				_, err := reply.Send()
				err_handler.HandleErr(err)
			}
		}

		notif := sql.GetNotification(user.Id)
		if notif != nil {
			if notif.Notification == "true" {
				txt := function.GetStringf(user.Id,
					"unamep",
					map[string]string{
						"1": strconv.Itoa(user.Id),
						"2": html.EscapeString(user.FirstName),
						"3": db.Action,
						"4": strconv.Itoa(user.Id),
						"5": html.EscapeString(chat.Title)},
				)
				reply.Text = txt
				reply.ReplyToMessageId = 0
				reply.ChatId = user.Id
				_, err := reply.Send()
				err_handler.HandleErr(err)
			}
		}
	} else {
		limit := sql.GetWarnSetting(strconv.Itoa(chat.Id))
		warns, _ := sql.WarnUser(
			strconv.Itoa(user.Id),
			strconv.Itoa(chat.Id),
			"Username",
		)

		if warns >= limit {
			sql.ResetWarns(strconv.Itoa(user.Id), strconv.Itoa(chat.Id))
			replyText = function.GetStringf(
				msg.Chat.Id,
				"handlers/warn2",
				map[string]string{
					"1": strconv.Itoa(user.Id),
					"2": html.EscapeString(user.FirstName),
					"3": strconv.Itoa(warns),
					"4": strconv.Itoa(limit),
					"5": strconv.Itoa(user.Id),
				})
			_, err := chat.UnbanMember(user.Id)
			err_handler.HandleErr(err)
		} else {
			keyb := function.BuildKeyboardf(
				"data/keyboard/usermonitor.json",
				1,
				map[string]string{
					"1": "Cara Pasang Username",
					"2": function.GetString(chat.Id, "rmwarn"),
					"3": fmt.Sprintf("rmWarn(%v)", user.Id),
					"4": "https://www.telegra.ph/Cara-Membuat-Username-di-Klien-Telegram-01-28",
				})
			markup = &ext.InlineKeyboardMarkup{InlineKeyboard: &keyb}
			replyText = function.GetStringf(
				chat.Id,
				"handlers/warn",
				map[string]string{
					"1": strconv.Itoa(user.Id),
					"2": html.EscapeString(user.FirstName),
					"3": strconv.Itoa(warns),
					"4": strconv.Itoa(limit),
					"5": strconv.Itoa(user.Id),
				})
		}

		reply := b.NewSendableMessage(chat.Id, replyText)
		reply.ParseMode = parsemode.Html
		reply.ReplyToMessageId = msg.MessageId
		reply.ReplyMarkup = markup
		_, err := reply.Send()
		if err != nil {
			reply.ReplyToMessageId = 0
			_, err := reply.Send()
			err_handler.HandleErr(err)
		}

		notif := sql.GetNotification(user.Id)
		if notif != nil && notif.Notification == "true" {
			reply.ReplyMarkup = nil
			reply.ReplyToMessageId = 0
			reply.ChatId = user.Id
			_, err = reply.Send()
			err_handler.HandleErr(err)
		}
	}

	if db.Deletion == "true" {
		_, err := msg.Delete()
		err_handler.HandleErr(err)
	}

	err := function.SendLog(b, u, "username", "")
	return err
}

func pictureScan(b ext.Bot, u *gotgbot.Update) error {
	db := sql.GetPicture(u.EffectiveChat.Id)

	if db == nil {
		return nil
	}

	if db.Option != "true" {
		return gotgbot.ContinueGroups{}
	}

	user := u.EffectiveUser
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	if chat.Type != "supergroup" {
		return nil
	}

	if chat_status.IsUserAdmin(chat, user.Id) {
		return nil
	}

	if !chat_status.CanRestrict(b, chat) {
		return nil
	}

	photo, err := user.GetProfilePhotos(0, 0)
	if err != nil {
		err_handler.HandleErr(err)
		return gotgbot.ContinueGroups{}
	}

	if photo != nil && photo.TotalCount > 0 {
		return gotgbot.ContinueGroups{}
	}

	replyText := function.GetStringf(
		msg.Chat.Id,
		"handlers/listener/listener.go:173",
		map[string]string{
			"1": strconv.Itoa(user.Id),
			"2": html.EscapeString(user.FirstName),
			"3": db.Action,
			"4": strconv.Itoa(user.Id),
		})
	markup := &ext.InlineKeyboardMarkup{}
	reply := b.NewSendableMessage(chat.Id, replyText)
	reply.ParseMode = parsemode.Html
	reply.ReplyToMessageId = msg.MessageId

	if db.Action != "warn" {
		switch db.Action {
		case "mute":
			restrictSend := b.NewSendableRestrictChatMember(chat.Id, user.Id)
			restrictSend.UntilDate = extraction.ExtractTime(b, msg, sql.GetSetting(chat.Id).Time)
			_, err := restrictSend.Send()
			err_handler.HandleErr(err)
			keyb := function.BuildKeyboardf(
				"data/keyboard/usermonitor.json",
				1,
				map[string]string{
					"1": "Cara Pasang Foto Profil",
					"2": function.GetString(chat.Id, "handlers/listener/listener.go:179"),
					"3": fmt.Sprintf("pmute_%v_%v", user.Id, chat.Id),
					"4": "https://www.wikihow.com/Use-Telegram",
				})
			markup = &ext.InlineKeyboardMarkup{InlineKeyboard: &keyb}
		case "kick":
			_, err := b.UnbanChatMember(chat.Id, user.Id)
			err_handler.HandleErr(err)
			markup = nil
		case "ban":
			restrictSend := b.NewSendableKickChatMember(chat.Id, user.Id)
			restrictSend.UntilDate = -1
			_, err := restrictSend.Send()
			err_handler.HandleErr(err)
			keyb := function.BuildKeyboardf(
				"data/keyboard/usermonitor.json",
				1,
				map[string]string{
					"1": "Cara Pasang Foto Profil",
					"2": function.GetString(chat.Id, "handlers/listener/listener.go:184"),
					"3": fmt.Sprintf("pban_%v_%v", user.Id, chat.Id),
					"4": "https://www.wikihow.com/Use-Telegram",
				})
			markup = &ext.InlineKeyboardMarkup{InlineKeyboard: &keyb}
		}

		reply.ReplyMarkup = markup
		_, err := reply.Send()
		if err != nil {
			if err.Error() == "Bad Request: reply message not found" {
				reply.ReplyToMessageId = 0
				_, err := reply.Send()
				err_handler.HandleErr(err)
			}
		}

		notif := sql.GetNotification(user.Id)
		if notif != nil {
			if notif.Notification == "true" {
				txt := function.GetStringf(
					user.Id,
					"picturep",
					map[string]string{
						"1": strconv.Itoa(user.Id),
						"2": html.EscapeString(user.FirstName),
						"3": db.Action,
						"4": strconv.Itoa(user.Id),
						"5": chat.Title})
				reply.Text = txt
				reply.ReplyToMessageId = 0
				reply.ChatId = user.Id
				_, err := reply.Send()
				err_handler.HandleErr(err)
			}
		}
	} else {
		limit := sql.GetWarnSetting(strconv.Itoa(chat.Id))
		warns, _ := sql.WarnUser(strconv.Itoa(user.Id), strconv.Itoa(chat.Id), "No Profile Picture")

		var reply = ""
		if warns >= limit {
			sql.ResetWarns(strconv.Itoa(user.Id), strconv.Itoa(chat.Id))
			reply = function.GetStringf(
				msg.Chat.Id,
				"handlers/warn2",
				map[string]string{
					"1": strconv.Itoa(user.Id),
					"2": html.EscapeString(user.FirstName),
					"3": strconv.Itoa(warns),
					"4": strconv.Itoa(limit),
					"5": strconv.Itoa(user.Id)})
			_, err := chat.UnbanMember(user.Id)
			err_handler.HandleErr(err)
		} else {
			keyb := function.BuildKeyboardf(
				"data/keyboard/usermonitor.json",
				1,
				map[string]string{
					"1": "Cara Pasang Foto Profil",
					"2": function.GetString(chat.Id, "rmwarn"),
					"3": fmt.Sprintf("rmWarn(%v)", user.Id),
					"4": "https://www.wikihow.com/Use-Telegram",
				})
			markup = &ext.InlineKeyboardMarkup{InlineKeyboard: &keyb}
			reply = function.GetStringf(
				msg.Chat.Id,
				"handlers/warn3",
				map[string]string{
					"1": strconv.Itoa(user.Id),
					"2": html.EscapeString(user.FirstName),
					"3": strconv.Itoa(warns),
					"4": strconv.Itoa(limit),
					"5": strconv.Itoa(user.Id)})
		}

		msgs := b.NewSendableMessage(chat.Id, reply)
		msgs.ParseMode = parsemode.Html
		msgs.ReplyToMessageId = msg.MessageId
		msgs.ReplyMarkup = markup
		_, err := msgs.Send()

		if err != nil {
			msgs.ReplyToMessageId = 0
			_, err := msgs.Send()
			err_handler.HandleErr(err)
		}

		if sql.GetNotification(user.Id).Notification == "true" {
			msgs.ReplyMarkup = nil
			msgs.ReplyToMessageId = 0
			msgs.ChatId = user.Id
			_, err = msgs.Send()
			err_handler.HandleErr(err)
		}
	}

	if db.Deletion == "true" {
		_, err := msg.Delete()
		err_handler.HandleErr(err)
	}

	err = function.SendLog(b, u, "picture", "")
	return err
}

func welcomeScan(b ext.Bot, u *gotgbot.Update) error {
	db := sql.GetVerify(u.EffectiveChat.Id)
	if db == nil {
		return nil
	}

	if db.Option != "true" {
		return gotgbot.ContinueGroups{}
	}

	user := u.EffectiveUser
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	if chat.Type != "supergroup" {
		return gotgbot.ContinueGroups{}
	}

	if chat_status.IsUserAdmin(chat, user.Id) {
		return nil
	}

	if !chat_status.CanRestrict(b, chat) {
		return nil
	}

	replyText := function.GetStringf(
		msg.Chat.Id,
		"handlers/listener/listener.go:298",
		map[string]string{
			"1": strconv.Itoa(user.Id),
			"2": html.EscapeString(user.FirstName),
			"3": chat.Title,
			"4": strconv.Itoa(user.Id)})

	kb := make([][]ext.InlineKeyboardButton, 1)
	kb[0] = make([]ext.InlineKeyboardButton, 1)
	kb[0][0] = ext.InlineKeyboardButton{
		Text:         function.GetString(chat.Id, "handlers/listener/listener.go:303"),
		CallbackData: fmt.Sprintf("wlcm_%v", user.Id),
	}

	restrictSend := b.NewSendableRestrictChatMember(chat.Id, user.Id)
	restrictSend.UntilDate = extraction.ExtractTime(b, msg, sql.GetSetting(chat.Id).Time)
	_, err := restrictSend.Send()

	err_handler.HandleErr(err)
	reply := b.NewSendableMessage(chat.Id, replyText)
	reply.ReplyMarkup = &ext.InlineKeyboardMarkup{InlineKeyboard: &kb}
	reply.ParseMode = parsemode.Html
	reply.ReplyToMessageId = msg.MessageId
	_, err = reply.Send()

	if err != nil {
		if err.Error() == "Bad Request: reply message not found" {
			reply := b.NewSendableMessage(chat.Id, replyText)
			reply.ReplyMarkup = &ext.InlineKeyboardMarkup{InlineKeyboard: &kb}
			reply.ParseMode = parsemode.Html
			_, err = reply.Send()
			return err
		}
	}

	if db.Deletion == "true" {
		_, err = msg.Delete()
		err_handler.HandleErr(err)
	}

	err = function.SendLog(b, u, "welcome", "")
	return err
}

func spamScan(b ext.Bot, u *gotgbot.Update) error {
	user := u.EffectiveUser
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	if chat.Type != "supergroup" {
		return gotgbot.ContinueGroups{}
	}

	if !chat_status.CanRestrict(b, chat) {
		return nil
	}

	if chat_status.IsUserAdmin(chat, msg.From.Id) {
		return gotgbot.ContinueGroups{}
	}

	if sql.GetEnforceGban(chat.Id).Option != "true" {
		return gotgbot.ContinueGroups{}
	}

	ban := sql.GetUserSpam(user.Id)
	if ban != nil {
		err := function.SpamFunc(b, u)
		err_handler.HandleErr(err)
		err = function.SendLog(b, u, "spam", ban.Reason)
		return gotgbot.EndGroups{}
	}

	_ = function.CasListener(b, u)
	return gotgbot.ContinueGroups{}
}

func linkScan(b ext.Bot, u *gotgbot.Update) error {
	db := sql.GetAntispam(u.Message.Chat.Id)
	if db == nil {
		return gotgbot.ContinueGroups{}
	}

	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat.Type != "supergroup" {
		return nil
	}

	if !chat_status.CanDelete(b, chat) {
		return gotgbot.ContinueGroups{}
	}

	if chat_status.IsUserAdmin(chat, user.Id) {
		return nil
	}

	if db.Deletion != "true" {
		return gotgbot.ContinueGroups{}
	}

	if db.Link == "true" {
		accepted := make(map[string]struct{})
		accepted["url"] = struct{}{}
		accepted["text_link"] = struct{}{}

		entities := msg.ParseEntityTypes(accepted)

		var ent *ext.ParsedMessageEntity = nil
		if len(entities) > 0 {
			ent = &entities[0]
		} else {
			ent = nil
		}

		if entities != nil && ent != nil {
			_, err := msg.Delete()
			err_handler.HandleErr(err)
			val := map[string]string{
				"1": fmt.Sprint(user.Id),
				"2": html.EscapeString(user.FirstName),
				"3": fmt.Sprint(user.Id),
				"4": "Link"}
			replyText := function.GetStringf(chat.Id, "handlers/listener/listener.go:dellink", val)
			reply := b.NewSendableMessage(chat.Id, replyText)
			reply.ParseMode = parsemode.Html
			_, err = reply.Send()
			err = function.SendLog(b, u, "link", "")
			return err
		}
	}
	if db.Forward == "true" {
		if msg.ForwardFrom != nil || msg.ForwardFromChat != nil {
			_, err := msg.Delete()
			err_handler.HandleErr(err)
			val := map[string]string{
				"1": fmt.Sprint(user.Id),
				"2": html.EscapeString(user.FirstName),
				"3": fmt.Sprint(user.Id),
				"4": " Forwarded Message",
			}
			replyText := function.GetStringf(chat.Id, "handlers/listener/listener.go:dellink", val)
			reply := b.NewSendableMessage(chat.Id, replyText)
			reply.ParseMode = parsemode.Html
			_, err = reply.Send()
			err = function.SendLog(b, u, "link", "")
			return err
		}
	}
	if db.Arabs == "true" {
		if function.CheckArabs(msg.Text) {
			_, err := msg.Delete()
			err_handler.HandleErr(err)
			val := map[string]string{
				"1": fmt.Sprint(user.Id),
				"2": html.EscapeString(user.FirstName),
				"3": fmt.Sprint(user.Id),
				"4": " Arabic Text",
			}
			replyText := function.GetStringf(chat.Id, "handlers/listener/listener.go:dellink", val)
			reply := b.NewSendableMessage(chat.Id, replyText)
			reply.ParseMode = parsemode.Html
			_, err = reply.Send()
			err = function.SendLog(b, u, "link", "")
			return err
		}

		if function.CheckChinese(msg.Text) {
			_, err := msg.Delete()
			err_handler.HandleErr(err)
			val := map[string]string{
				"1": fmt.Sprint(user.Id),
				"2": html.EscapeString(user.FirstName),
				"3": fmt.Sprint(user.Id),
				"4": " Chinese Text",
			}
			replyText := function.GetStringf(chat.Id, "handlers/listener/listener.go:dellink", val)
			reply := b.NewSendableMessage(chat.Id, replyText)
			reply.ParseMode = parsemode.Html
			_, err = reply.Send()
			err = function.SendLog(b, u, "link", "")
			return err
		}

		if function.CheckChinese(user.FirstName) {
			_, err := msg.Delete()
			err_handler.HandleErr(err)
			_, err = b.UnbanChatMember(chat.Id, user.Id)
			err_handler.HandleErr(err)
			val := map[string]string{
				"1": fmt.Sprint(user.Id),
				"2": html.EscapeString(user.FirstName),
				"3": fmt.Sprint(user.Id),
				"4": " Chinese Name",
			}
			replyText := function.GetStringf(chat.Id, "handlers/listener/listener.go:dellink", val)
			reply := b.NewSendableMessage(chat.Id, replyText)
			reply.ParseMode = parsemode.Html
			_, err = reply.Send()
			err = function.SendLog(b, u, "link", "")
			return err
		}

		if function.CheckArabs(user.FirstName) {
			_, err := msg.Delete()
			err_handler.HandleErr(err)
			_, err = b.UnbanChatMember(chat.Id, user.Id)
			err_handler.HandleErr(err)
			val := map[string]string{
				"1": fmt.Sprint(user.Id),
				"2": html.EscapeString(user.FirstName),
				"3": fmt.Sprint(user.Id),
				"4": " Arabic Name",
			}
			replyText := function.GetStringf(chat.Id, "handlers/listener/listener.go:dellink", val)
			reply := b.NewSendableMessage(chat.Id, replyText)
			reply.ParseMode = parsemode.Html
			_, err = reply.Send()
			err = function.SendLog(b, u, "link", "")
			return err
		}

	}
	return gotgbot.ContinueGroups{}
}

// Thanks to https://github.com/mojurasu/kantek for the kriminalamt plugin
func kriminalamtHandler(b ext.Bot, u *gotgbot.Update) error {
	user := u.EffectiveUser
	chat := u.EffectiveChat
	//msg := u.EffectiveMessage

	time.Sleep(1 * time.Second)
	getChat, err := b.GetChatMember(chat.Id, user.Id)
	if err != nil {
		err_handler.HandleErr(err)
		fmt.Print(getChat)
	}

	if getChat != nil && getChat.Status == "left" {
		rson := "AutoKriminalamt #" + strings.Replace(fmt.Sprint(chat.Id), "-100", "", 1) + " No. 1"
		sql.UpdateUserSpam(
			user.Id,
			rson,
			fmt.Sprint(b.Id),
			int(time.Now().Unix()),
		)
		err := function.SpamFunc(b, u)
		err_handler.HandleErr(err)
		err = function.SendLog(b, u, "spam", rson)
		return gotgbot.EndGroups{}
	}

	return gotgbot.ContinueGroups{}
}

func update(_ ext.Bot, u *gotgbot.Update) error {
	user := u.EffectiveUser
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	if msg != nil {
		sql.UpdateChat(strconv.Itoa(chat.Id), chat.Title, chat.Type, chat.InviteLink)
		sql.UpdateUser(user.Id, user.Username, user.FirstName, user.LastName)

		if sql.GetLang(chat.Id) == nil {
			caching.REDIS.Set(fmt.Sprintf("lang_%v", chat.Id), "id", 7200)
			sql.UpdateLang(chat.Id, "id")
		}
		if msg.ForwardFrom != nil {
			usr := msg.ForwardFrom
			sql.UpdateUser(usr.Id, usr.Username, usr.FirstName, usr.LastName)
		}

		if chat.Type == "supergroup" {
			if sql.GetVerify(chat.Id) == nil {
				sql.UpdateVerify(chat.Id, "true", "-", "true")
			}
			if sql.GetUsername(chat.Id) == nil {
				sql.UpdateUsername(chat.Id, "true", "mute", "-", "true")
			}
			if sql.GetPicture(chat.Id) == nil {
				sql.UpdatePicture(chat.Id, "true", "mute", "-", "true")
			}
			if sql.GetSetting(chat.Id) == nil {
				sql.UpdateSetting(chat.Id, "5m", "true")
			}
			if sql.GetEnforceGban(chat.Id) == nil {
				sql.UpdateEnforceGban(chat.Id, "true")
			}
			if sql.GetAntispam(chat.Id) == nil {
				sql.UpdateAntispam(chat.Id, "false", "true", "false", "false")
			}
		}

		if sql.GetNotification(user.Id) == nil {
			sql.UpdateNotification(user.Id, "true")
		}
		return gotgbot.ContinueGroups{}
	}
	return nil
}

func usernameQuery(b ext.Bot, u *gotgbot.Update) error {
	msg := u.CallbackQuery
	user := msg.From
	chat := msg.Message.Chat

	data, _ := regexp.MatchString("^umute_\\d+_.", msg.Data)
	data2, _ := regexp.MatchString("^uba_\\d+_.", msg.Data)
	if data {
		splt := strings.Split(msg.Data, "umute_")[1]
		if strings.Split(splt, "_")[0] == strconv.Itoa(user.Id) {
			if user.Username == "" {
				cb := b.NewSendableAnswerCallbackQuery(msg.Id)
				cb.Text = function.GetString(chat.Id, "handlers/listener/listener.go:454")
				cb.ShowAlert = true
				cb.CacheTime = 10
				_, err := cb.Send()
				return err
			}

			if chat.Type == "supergroup" {
				_, err := b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id,
					"handlers/listener/listener.go:441"), true)
				if err != nil {
					_, err = b.AnswerCallbackQueryText(msg.Id, err.Error(), true)
					return err
				}
				_, err = msg.Message.Delete()
				if err != nil {
					_, err = b.AnswerCallbackQueryText(msg.Id, err.Error(), true)
					return err
				}
				_, err = b.UnRestrictChatMember(chat.Id, user.Id)
				return err
			}

			_, err := b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id,
				"handlers/listener/listener.go:441"), true)
			err_handler.HandleErr(err)
			_, err = msg.Message.Delete()
			err_handler.HandleErr(err)
			cid, _ := strconv.Atoi(strings.Split(splt, "_")[1])
			_, err = b.UnRestrictChatMember(cid, user.Id)
			return err
		}
		_, err := b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id,
			"handlers/listener/listener.go:458"), true)
		return err
	} else if data2 == true {
		if chat.Type == "supergroup" {
			if chat_status.IsUserAdmin(chat, user.Id) == true {
				i, _ := strconv.Atoi(strings.Split(msg.Data, "uba_")[1])
				_, err := b.UnbanChatMember(chat.Id, i)
				err_handler.HandleErr(err)
				_, err = b.AnswerCallbackQueryText(msg.Id, "Unbanned!", true)
				err_handler.HandleErr(err)
				_, err = msg.Message.Delete()
				err_handler.HandleErr(err)
				return err
			}
		} else if chat.Type == "private" {
			_, err := b.AnswerCallbackQuery(msg.Id)
			return err
		}
	}
	return gotgbot.ContinueGroups{}
}

func pictureQuery(b ext.Bot, u *gotgbot.Update) error {
	msg := u.CallbackQuery
	user := msg.From
	chat := msg.Message.Chat

	data, _ := regexp.MatchString("^pmute_\\d+_.", msg.Data)
	data2, _ := regexp.MatchString("^pba_\\d+_.", msg.Data)
	if data == true {
		splt := strings.Split(msg.Data, "pmute_")[1]
		if strings.Split(splt, "_")[0] == strconv.Itoa(user.Id) {
			photo, _ := user.GetProfilePhotos(0, 0)
			if photo != nil && photo.TotalCount == 0 {
				_, err := b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id, "handlers/listener/listener.go:511"), true)
				return err
			}

			if chat.Type == "supergroup" {
				_, err := b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id, "handlers/listener/listener.go:441"), true)
				err_handler.HandleErr(err)
				_, err = msg.Message.Delete()
				err_handler.HandleErr(err)
				_, err = b.UnRestrictChatMember(chat.Id, user.Id)
				err_handler.HandleErr(err)
				return err
			}

			_, err := b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id, "handlers/listener/listener.go:441"), true)
			err_handler.HandleErr(err)
			_, err = msg.Message.Delete()
			err_handler.HandleErr(err)
			cid, _ := strconv.Atoi(strings.Split(splt, "_")[1])
			_, err = b.UnRestrictChatMember(cid, user.Id)
			return err
		}
		_, err := b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id, "handlers/listener/listener.go:515"), true)
		err_handler.HandleErr(err)
		return err
	} else if data2 == true {
		if chat.Type == "supergroup" {
			if chat_status.IsUserAdmin(chat, user.Id) == true {
				i, _ := strconv.Atoi(strings.Split(msg.Data, "uba_")[1])
				_, err := b.UnbanChatMember(chat.Id, i)
				_, err = b.AnswerCallbackQueryText(msg.Id, "Unbanned!", true)
				err_handler.HandleErr(err)
				_, err = msg.Message.Delete()
				return err
			}
		} else if chat.Type == "private" {
			_, err := b.AnswerCallbackQuery(msg.Id)
			return err
		}

	}
	return gotgbot.ContinueGroups{}
}

func verifyQuery(b ext.Bot, u *gotgbot.Update) error {
	msg := u.CallbackQuery
	user := msg.From
	chat := msg.Message.Chat

	if chat.Type == "supergroup" {
		if strings.Split(msg.Data, "wlcm_")[1] == strconv.Itoa(user.Id) {
			_, err := b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id, "handlers/listener/listener.go:552"), true)
			err_handler.HandleErr(err)
			_, err = msg.Message.Delete()
			err_handler.HandleErr(err)
			_, err = b.UnRestrictChatMember(chat.Id, user.Id)
			return err
		}
		_, err := b.AnswerCallbackQueryText(msg.Id,
			function.GetString(chat.Id, "handlers/listener/listener.go:566"), true)
		return err

	}
	return gotgbot.ContinueGroups{}
}

func warnQuery(bot ext.Bot, u *gotgbot.Update) error {
	query := u.CallbackQuery
	user := u.EffectiveUser
	chat := u.EffectiveChat
	pattern, _ := regexp.Compile(`rmWarn\((.+?)\)`)

	if chat_status.IsUserAdmin(chat, user.Id) == false {
		_, err := bot.AnswerCallbackQueryText(query.Id, "You need to be an admin to do this.", true)
		return err
	}

	if pattern.MatchString(query.Data) {
		userID := pattern.FindStringSubmatch(query.Data)[1]
		chat := u.EffectiveChat
		res := sql.RemoveWarn(userID, strconv.Itoa(chat.Id))
		if res {
			_, err := bot.AnswerCallbackQueryText(query.Id, "Warn Removed.", true)
			err_handler.HandleErr(err)
			_, err = u.EffectiveMessage.Delete()
			return err
		}
	}
	return gotgbot.ContinueGroups{}
}

func LoadUserListener(u *gotgbot.Updater) {
	defer logrus.Info("Usermonitor listener Loaded...")
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, update))
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.NewChatMembers(), kriminalamtHandler))
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, spamScan))
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, linkScan))
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, usernameScan))
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, pictureScan))
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.NewChatMembers(), welcomeScan))
	u.Dispatcher.AddHandler(handlers.NewCallback("^(umute|uba)_", usernameQuery))
	u.Dispatcher.AddHandler(handlers.NewCallback("^(pmute|pban)_", pictureQuery))
	u.Dispatcher.AddHandler(handlers.NewCallback("wlcm", verifyQuery))
	u.Dispatcher.AddHandler(handlers.NewCallback("rmWarn", warnQuery))
}
