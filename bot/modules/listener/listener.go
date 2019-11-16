package listener

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/handlers/Filters"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/caching"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/chat_status"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/extraction"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/logger"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
	"regexp"
	"strconv"
	"strings"
)

/* type spammer struct {
 *	Status bool `json:"ok"`
}*/

func username(b ext.Bot, u *gotgbot.Update) error {
	var err error
	user := u.EffectiveUser
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	db := sql.GetUsername(chat.Id)
	if db.Option != "true" {
		return nil
	}
	if chat_status.IsUserAdmin(chat, msg.From.Id) == true {
		return gotgbot.EndGroups{}
	}
	if chat.Type != "supergroup" {
		return nil
	}
	if user.Username != "" {
		return nil
	}

	bantime := extraction.ExtractTime(b, msg, sql.GetSetting(chat.Id).Time)
	replytext := function.GetStringf(msg.Chat.Id, "modules/listener/listener.go:45",
		map[string]string{"1": strconv.Itoa(user.Id), "2": user.FirstName, "3": db.Action,
			"4": strconv.Itoa(user.Id)})

	kb := make([][]ext.InlineKeyboardButton, 1)
	kb[0] = make([]ext.InlineKeyboardButton, 1)
	kb[0][0] = ext.InlineKeyboardButton{Text: function.GetString(chat.Id, "modules/listener/listener.go:51"),
		CallbackData: fmt.Sprintf("umute_%v_%v", user.Id, chat.Id)}

	kbk := make([][]ext.InlineKeyboardButton, 1)
	kbk[0] = make([]ext.InlineKeyboardButton, 1)
	kbk[0][0] = ext.InlineKeyboardButton{Text: function.GetString(chat.Id, "modules/listener/listener.go:56"),
		CallbackData: fmt.Sprintf("uba_%v_%v", user.Id, chat.Id)}

	markup := &ext.InlineKeyboardMarkup{InlineKeyboard: &kb}

	if db.Action != "warn" {
		reply := b.NewSendableMessage(chat.Id, replytext)
		reply.ParseMode = parsemode.Html
		reply.ReplyToMessageId = msg.MessageId

		if db.Action == "mute" {
			restrictSend := b.NewSendableRestrictChatMember(chat.Id, user.Id)
			restrictSend.UntilDate = bantime
			_, err = restrictSend.Send()
			if err != nil {
				if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
					_, err = b.SendMessage(chat.Id, err.Error())
					return err
				}
			}
			markup = &ext.InlineKeyboardMarkup{InlineKeyboard: &kb}
		} else if db.Action == "kick" {
			_, err = b.UnbanChatMember(chat.Id, user.Id)
			if err != nil {
				if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
					_, err = b.SendMessage(chat.Id, err.Error())
					return err
				}
			}
			markup = nil
		} else if db.Action == "ban" {
			restrictSend := b.NewSendableKickChatMember(chat.Id, user.Id)
			restrictSend.UntilDate = -1
			_, err = restrictSend.Send()
			if err != nil {
				if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
					_, err = b.SendMessage(chat.Id, err.Error())
					return err
				}
			}
			markup = &ext.InlineKeyboardMarkup{InlineKeyboard: &kbk}
		}
		reply.ReplyMarkup = markup
		_, err = reply.Send()
		if err != nil {
			if err.Error() == "Bad Request: reply message not found" {
				reply.ReplyToMessageId = 0
				_, err = reply.Send()
				return err
			}
			err_handler.HandleErr(err)
		}

		if sql.GetNotification(user.Id).Notification == "true" {
			txt := function.GetStringf(user.Id, "unamep",
				map[string]string{"1": strconv.Itoa(user.Id), "2": user.FirstName, "3": db.Action,
					"4": strconv.Itoa(user.Id), "5": chat.Title})
			reply.Text = txt
			reply.ReplyToMessageId = 0
			reply.ChatId = user.Id
			_, err = reply.Send()
			err_handler.HandleErr(err)
		}
	} else {
		limit := sql.GetWarnSetting(strconv.Itoa(chat.Id))
		warns, _ := sql.WarnUser(strconv.Itoa(user.Id), strconv.Itoa(chat.Id), "Username")

		var keyboard ext.InlineKeyboardMarkup
		if warns >= limit {
			go sql.ResetWarns(strconv.Itoa(user.Id), strconv.Itoa(chat.Id))
			replytext = function.GetStringf(msg.Chat.Id, "modules/warn2",
				map[string]string{"1": strconv.Itoa(user.Id), "2": user.FirstName, "3": strconv.Itoa(warns),
					"4": strconv.Itoa(limit), "5": strconv.Itoa(user.Id)})

			_, err = chat.UnbanMember(user.Id)
		} else {
			kb := make([][]ext.InlineKeyboardButton, 1)
			kb[0] = make([]ext.InlineKeyboardButton, 1)
			kb[0][0] = ext.InlineKeyboardButton{Text: function.GetString(chat.Id, "rmwarn"),
				CallbackData: fmt.Sprintf("rmWarn(%v)", user.Id)}
			keyboard = ext.InlineKeyboardMarkup{InlineKeyboard: &kb}
			replytext = function.GetStringf(msg.Chat.Id, "modules/warn",
				map[string]string{"1": strconv.Itoa(user.Id), "2": user.FirstName, "3": strconv.Itoa(warns),
					"4": strconv.Itoa(limit), "5": strconv.Itoa(user.Id)})
		}

		msgs := b.NewSendableMessage(chat.Id, replytext)
		msgs.ParseMode = parsemode.Html
		msgs.ReplyToMessageId = msg.MessageId
		msgs.ReplyMarkup = &keyboard
		_, err = msgs.Send()
		if err != nil {
			msgs.ReplyToMessageId = 0
			_, err = msgs.Send()
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
		_, err = msg.Delete()
		if err != nil {
			if err.Error() == "Bad Request: message can't be deleted" {
				_, err = msg.ReplyText(err.Error())
				return err
			}
		}
	}
	err = logger.SendLog(b, u, "username", "")
	err_handler.HandleErr(err)
	return gotgbot.ContinueGroups{}
}

func picture(b ext.Bot, u *gotgbot.Update) error {
	var err error
	user := u.EffectiveUser
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	db := sql.GetPicture(chat.Id)
	if db.Option != "true" {
		return nil
	}
	if chat_status.IsUserAdmin(chat, msg.From.Id) == true {
		return gotgbot.EndGroups{}
	}
	if chat.Type != "supergroup" {
		return nil
	}
	photo, _ := user.GetProfilePhotos(0, 0)
	if photo != nil && photo.TotalCount != 0 {
		return nil
	}

	bantime := extraction.ExtractTime(b, msg, sql.GetSetting(chat.Id).Time)
	replytext := function.GetStringf(msg.Chat.Id, "modules/listener/listener.go:173",
		map[string]string{"1": strconv.Itoa(user.Id), "2": user.FirstName, "3": db.Action,
			"4": strconv.Itoa(user.Id)})

	kb := make([][]ext.InlineKeyboardButton, 1)
	kb[0] = make([]ext.InlineKeyboardButton, 1)
	kb[0][0] = ext.InlineKeyboardButton{Text: function.GetString(chat.Id, "modules/listener/listener.go:179"),
		CallbackData: fmt.Sprintf("pmute_%v_%v", user.Id, chat.Id)}

	kbk := make([][]ext.InlineKeyboardButton, 1)
	kbk[0] = make([]ext.InlineKeyboardButton, 1)
	kbk[0][0] = ext.InlineKeyboardButton{Text: function.GetString(chat.Id, "modules/listener/listener.go:184"),
		CallbackData: fmt.Sprintf("pban_%v_%v", user.Id, chat.Id)}

	markup := &ext.InlineKeyboardMarkup{}

	reply := b.NewSendableMessage(chat.Id, replytext)
	reply.ParseMode = parsemode.Html
	reply.ReplyToMessageId = msg.MessageId

	if db.Action != "warn" {
		if db.Action == "mute" {
			restrictSend := b.NewSendableRestrictChatMember(chat.Id, user.Id)
			restrictSend.UntilDate = bantime
			_, err = restrictSend.Send()
			if err != nil {
				if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
					_, err = b.SendMessage(chat.Id, err.Error())
					return err
				}
			}
			markup = &ext.InlineKeyboardMarkup{InlineKeyboard: &kb}
		} else if db.Action == "kick" {
			_, err = b.UnbanChatMember(chat.Id, user.Id)
			if err != nil {
				if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
					_, err = b.SendMessage(chat.Id, err.Error())
					return err
				}
			}
			markup = nil
		} else if db.Action == "ban" {
			restrictSend := b.NewSendableKickChatMember(chat.Id, user.Id)
			restrictSend.UntilDate = -1
			_, err = restrictSend.Send()
			if err != nil {
				if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
					_, err = b.SendMessage(chat.Id, err.Error())
					return err
				}
			}
			markup = &ext.InlineKeyboardMarkup{InlineKeyboard: &kbk}
		}
		reply.ReplyMarkup = markup
		_, err = reply.Send()
		if err != nil {
			if err.Error() == "Bad Request: reply message not found" {
				reply.ReplyToMessageId = 0
				_, err = reply.Send()
				return err
			}
			err_handler.HandleErr(err)
		}
		if sql.GetNotification(user.Id).Notification == "true" {
			txt := function.GetStringf(user.Id, "picturep",
				map[string]string{"1": strconv.Itoa(user.Id), "2": user.FirstName, "3": db.Action,
					"4": strconv.Itoa(user.Id), "5": chat.Title})
			reply.Text = txt
			reply.ReplyToMessageId = 0
			reply.ChatId = user.Id
			_, err = reply.Send()
			err_handler.HandleErr(err)
		}
	} else {
		limit := sql.GetWarnSetting(strconv.Itoa(chat.Id))
		warns, _ := sql.WarnUser(strconv.Itoa(user.Id), strconv.Itoa(chat.Id), "No Profile Picture")

		var reply string
		var keyboard ext.InlineKeyboardMarkup
		if warns >= limit {
			go sql.ResetWarns(strconv.Itoa(user.Id), strconv.Itoa(chat.Id))
			reply = function.GetStringf(msg.Chat.Id, "modules/warn2",
				map[string]string{"1": strconv.Itoa(user.Id), "2": user.FirstName, "3": strconv.Itoa(warns),
					"4": strconv.Itoa(limit), "5": strconv.Itoa(user.Id)})
			_, err = chat.UnbanMember(user.Id)
		} else {
			kb := make([][]ext.InlineKeyboardButton, 1)
			kb[0] = make([]ext.InlineKeyboardButton, 1)
			kb[0][0] = ext.InlineKeyboardButton{Text: function.GetString(chat.Id, "rmwarn"),
				CallbackData: fmt.Sprintf("rmWarn(%v)", user.Id)}
			keyboard = ext.InlineKeyboardMarkup{InlineKeyboard: &kb}
			reply = function.GetStringf(msg.Chat.Id, "modules/warn3",
				map[string]string{"1": strconv.Itoa(user.Id), "2": user.FirstName, "3": strconv.Itoa(warns),
					"4": strconv.Itoa(limit), "5": strconv.Itoa(user.Id)})
		}

		msgs := b.NewSendableMessage(chat.Id, reply)
		msgs.ParseMode = parsemode.Html
		msgs.ReplyToMessageId = msg.MessageId
		msgs.ReplyMarkup = &keyboard
		_, err = msgs.Send()
		if err != nil {
			msgs.ReplyToMessageId = 0
			_, err = msgs.Send()
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
		_, err = msg.Delete()
		if err != nil {
			if err.Error() == "Bad Request: message can't be deleted" {
				_, err = msg.ReplyText(err.Error())
				return err
			}
		}
	}
	err = logger.SendLog(b, u, "picture", "")
	err_handler.HandleErr(err)
	return gotgbot.ContinueGroups{}
}

func verify(b ext.Bot, u *gotgbot.Update) error {
	var err error
	user := u.EffectiveUser
	chat := u.EffectiveChat
	msg := u.EffectiveMessage
	db := sql.GetVerify(chat.Id)

	if db.Option != "true" {
		return nil
	}
	if chat_status.IsUserAdmin(chat, msg.From.Id) == true {
		return gotgbot.EndGroups{}
	}

	bantime := extraction.ExtractTime(b, msg, sql.GetSetting(chat.Id).Time)
	replytext := function.GetStringf(msg.Chat.Id, "modules/listener/listener/listener.go:298",
		map[string]string{"1": strconv.Itoa(user.Id), "2": user.FirstName, "3": chat.Title, "4": strconv.Itoa(user.Id)})

	kb := make([][]ext.InlineKeyboardButton, 1)
	kb[0] = make([]ext.InlineKeyboardButton, 1)
	kb[0][0] = ext.InlineKeyboardButton{Text: function.GetString(chat.Id, "modules/listener/listener.go:303"),
		CallbackData: fmt.Sprintf("wlcm_%v", user.Id)}

	restrictSend := b.NewSendableRestrictChatMember(chat.Id, user.Id)
	restrictSend.UntilDate = bantime
	_, err = restrictSend.Send()
	if err != nil {
		if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
			_, err = b.SendMessage(chat.Id, err.Error())
			return err
		}
	}

	reply := b.NewSendableMessage(chat.Id, replytext)
	reply.ReplyMarkup = &ext.InlineKeyboardMarkup{InlineKeyboard: &kb}
	reply.ParseMode = parsemode.Html
	reply.ReplyToMessageId = msg.MessageId
	_, err = reply.Send()
	if err != nil {
		if err.Error() == "Bad Request: reply message not found" {
			reply := b.NewSendableMessage(chat.Id, replytext)
			reply.ReplyMarkup = &ext.InlineKeyboardMarkup{InlineKeyboard: &kb}
			reply.ParseMode = parsemode.Html
			_, err = reply.Send()
			return err
		}
	}
	if db.Deletion == "true" {
		_, err = msg.Delete()
		if err != nil {
			if err.Error() == "Bad Request: message can't be deleted" {
				_, err = msg.ReplyText(err.Error())
				return err
			}
		}
	}
	err = logger.SendLog(b, u, "welcome", "")
	return gotgbot.ContinueGroups{}
}

func spam(b ext.Bot, u *gotgbot.Update) error {
	user := u.EffectiveUser
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	if chat_status.IsUserAdmin(chat, msg.From.Id) == true {
		return gotgbot.EndGroups{}
	}
	if chat_status.IsBotAdmin(chat, nil) == false {
		return gotgbot.EndGroups{}
	}
	if chat.Type != "supergroup" {
		return nil
	}
	if sql.GetEnforceGban(chat.Id).Option != "true" {
		return nil
	}
	ban := sql.GetUserSpam(user.Id)
	if ban != nil {
		err := spamfunc(b, u)
		err_handler.HandleErr(err)
		err = logger.SendLog(b, u, "spam", ban.Reason)
		return gotgbot.EndGroups{}
	}

	return gotgbot.ContinueGroups{}
}

func update(_ ext.Bot, u *gotgbot.Update) error {
	user := u.EffectiveUser
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	if msg != nil {
		db := make(chan error)
		go func() { db <- sql.UpdateUser(user.Id, user.Username, user.FirstName) }()
		err_handler.HandleErr(<-db)
		db = make(chan error)
		go func() { db <- sql.UpdateChat(strconv.Itoa(chat.Id), chat.Title, chat.Type, chat.InviteLink) }()
		err_handler.HandleErr(<-db)

		if msg.ForwardFrom != nil {
			usr := msg.ForwardFrom
			db := make(chan error)
			go func() { db <- sql.UpdateUser(usr.Id, usr.Username, usr.FirstName) }()
			err_handler.HandleErr(<-db)
		}

		if sql.GetVerify(chat.Id) == nil {
			db := make(chan error)
			go func() { db <- sql.UpdateVerify(chat.Id, "true", "-", "true") }()
			err_handler.HandleErr(<-db)
		}
		if sql.GetUsername(chat.Id) == nil {
			db := make(chan error)
			go func() { db <- sql.UpdateUsername(chat.Id, "true", "mute", "-", "true") }()
			err_handler.HandleErr(<-db)
		}
		if sql.GetPicture(chat.Id) == nil {
			db := make(chan error)
			go func() { db <- sql.UpdatePicture(chat.Id, "true", "mute", "-", "true") }()
			err_handler.HandleErr(<-db)
		}
		if sql.GetSetting(chat.Id) == nil {
			db := make(chan error)
			go func() { db <- sql.UpdateSetting(chat.Id, "5m", "true") }()
			err_handler.HandleErr(<-db)
		}
		if sql.GetEnforceGban(chat.Id) == nil {
			db := make(chan error)
			go func() { db <- sql.UpdateEnforceGban(chat.Id, "true") }()
			err_handler.HandleErr(<-db)
		}
		if sql.GetLang(chat.Id) == nil {
			caching.REDIS.Set(fmt.Sprintf("lang_%v", chat.Id), "id", 0)
			caching.REDIS.BgSave()
			db := make(chan error)
			go func() { db <- sql.UpdateLang(chat.Id, "id") }()
			err_handler.HandleErr(<-db)
		}
		if sql.GetNotification(user.Id) == nil {
			db := make(chan error)
			go func() { db <- sql.UpdateNotification(user.Id, "true") }()
			err_handler.HandleErr(<-db)
		}
	}
	return gotgbot.ContinueGroups{}
}

func usernamequery(b ext.Bot, u *gotgbot.Update) error {
	var err error
	msg := u.CallbackQuery
	user := msg.From
	chat := msg.Message.Chat

	data, _ := regexp.MatchString("^umute_\\d+_.", msg.Data)
	data2, _ := regexp.MatchString("^uba_\\d+_.", msg.Data)
	if data == true {
		splt := strings.Split(msg.Data, "umute_")[1]
		if strings.Split(splt, "_")[0] == strconv.Itoa(user.Id) {
			if user.Username == "" {
				_, err = b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id, "modules/listener/listener.go:454"), true)
				return err
			}

			if chat.Type == "supergroup" {
				_, err = b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id, "modules/listener/listener.go:441"), true)
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
			} else {
				_, err = b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id, "modules/listener/listener.go:441"), true)
				if err != nil {
					err_handler.HandleCbErr(b, u, err)
				}
				_, err = msg.Message.Delete()
				if err != nil {
					_, err = b.AnswerCallbackQueryText(msg.Id, err.Error(), true)
					return err
				}
				cid, _ := strconv.Atoi(strings.Split(splt, "_")[1])
				_, err = b.UnRestrictChatMember(cid, user.Id)
				return err
			}
		}
		_, err = b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id, "modules/listener/listener.go:458"), true)
		return err
	} else if data2 == true {
		if chat.Type == "supergroup" {
			if chat_status.IsUserAdmin(chat, user.Id) == true {
				i, _ := strconv.Atoi(strings.Split(msg.Data, "uba_")[1])
				_, err = b.UnbanChatMember(chat.Id, i)
				_, err = b.AnswerCallbackQueryText(msg.Id, "Unbanned!", true)
				if err != nil {
					_, err = b.AnswerCallbackQueryText(msg.Id, err.Error(), true)
					return err
				}
				_, err = msg.Message.Delete()
				if err != nil {
					_, err = b.AnswerCallbackQueryText(msg.Id, err.Error(), true)
					return err
				}
				return err
			}
		} else if chat.Type == "private" {
			_, err = b.AnswerCallbackQuery(msg.Id)
			return err
		}

	}

	return gotgbot.ContinueGroups{}
}

func picturequery(b ext.Bot, u *gotgbot.Update) error {
	var err error
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
				_, err = b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id, "modules/listener/listener.go:511"), true)
				return err
			}

			if chat.Type == "supergroup" {
				_, err = b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id, "modules/listener/listener.go:441"), true)
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
			} else {
				_, err = b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id, "modules/listener/listener.go:441"), true)
				if err != nil {
					err_handler.HandleCbErr(b, u, err)
				}
				_, err = msg.Message.Delete()
				if err != nil {
					_, err = b.AnswerCallbackQueryText(msg.Id, err.Error(), true)
					return err
				}
				cid, _ := strconv.Atoi(strings.Split(splt, "_")[1])
				_, err = b.UnRestrictChatMember(cid, user.Id)
				return err
			}
		}
		_, err = b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id, "modules/listener/listener.go:515"), true)
		return err
	} else if data2 == true {
		if chat.Type == "supergroup" {
			if chat_status.IsUserAdmin(chat, user.Id) == true {
				i, _ := strconv.Atoi(strings.Split(msg.Data, "uba_")[1])
				_, err = b.UnbanChatMember(chat.Id, i)
				_, err = b.AnswerCallbackQueryText(msg.Id, "Unbanned!", true)
				if err != nil {
					_, err = b.AnswerCallbackQueryText(msg.Id, err.Error(), true)
					return err
				}
				_, err = msg.Message.Delete()
				if err != nil {
					_, err = b.AnswerCallbackQueryText(msg.Id, err.Error(), true)
					return err
				}
				return err
			}
		} else if chat.Type == "private" {
			_, err = b.AnswerCallbackQuery(msg.Id)
			return err
		}

	}

	return gotgbot.ContinueGroups{}
}

func verifyquery(b ext.Bot, u *gotgbot.Update) error {
	var err error
	msg := u.CallbackQuery
	user := msg.From
	chat := msg.Message.Chat

	if chat.Type == "supergroup" {
		data, _ := regexp.MatchString("^wlcm_\\d+$", msg.Data)
		if data == true {
			if strings.Split(msg.Data, "wlcm_")[1] == strconv.Itoa(user.Id) {
				_, err = b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id, "modules/listener/listener.go:552"), true)
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
			_, err = b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id, "modules/listener/listener.go:566"), true)
			return err
		}
	}

	return gotgbot.ContinueGroups{}
}

func warnquery(bot ext.Bot, u *gotgbot.Update) error {
	query := u.CallbackQuery
	user := u.EffectiveUser
	chat := u.EffectiveChat
	pattern, _ := regexp.Compile(`rmWarn\((.+?)\)`)

	// Check permissions
	if chat_status.IsUserAdmin(chat, user.Id) == false {
		_, _ = bot.AnswerCallbackQueryText(query.Id, "You need to be an admin to do this.", true)
		return gotgbot.EndGroups{}
	}

	if pattern.MatchString(query.Data) {
		userId := pattern.FindStringSubmatch(query.Data)[1]
		chat := u.EffectiveChat
		res := sql.RemoveWarn(userId, strconv.Itoa(chat.Id))
		if res {
			msg := bot.NewSendableEditMessageText(chat.Id, u.EffectiveMessage.MessageId,
				fmt.Sprintf("Warn removed by admin %v.", user.FirstName))
			msg.ParseMode = parsemode.Html
			_, err := msg.Send()
			return err
		}
		_, err := u.EffectiveMessage.EditText("User already has no warns.")
		return err

	}
	return nil
}

func spamfunc(b ext.Bot, u *gotgbot.Update) error {
	user := u.EffectiveUser
	chat := u.EffectiveChat
	msg := u.EffectiveMessage
	db := sql.GetSetting(chat.Id)
	txtBan := function.GetStringf(chat.Id, "modules/listener/listener.go:580",
		map[string]string{"1": strconv.Itoa(user.Id), "2": user.FirstName, "3": strconv.Itoa(user.Id)})

	_, err := msg.ReplyHTMLf(txtBan)
	if err != nil {
		if err.Error() == "Bad Request: reply message not found" {
			_, err = b.SendMessageHTML(chat.Id, txtBan)
			return err
		}
	}
	_, err = b.KickChatMember(chat.Id, user.Id)
	if err != nil {
		if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
			_, err = b.SendMessage(chat.Id, err.Error())
			return err
		}
	}

	if db.Deletion == "true" {
		_, err = b.DeleteMessage(chat.Id, msg.MessageId)
		if err != nil {
			if err.Error() == "Bad Request: message can't be deleted" {
				_, err = b.SendMessage(chat.Id, err.Error())
				return err
			}
		}
	}
	return nil
}

// LoadListeners -> Register handlers
func LoadListeners(u *gotgbot.Updater) {
	go u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, update))
	go u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, spam))
	go u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, username))
	go u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, picture))
	go u.Dispatcher.AddHandler(handlers.NewMessage(Filters.NewChatMembers(), verify))
	go u.Dispatcher.AddHandler(handlers.NewCallback("^(umute|uba)_", usernamequery))
	go u.Dispatcher.AddHandler(handlers.NewCallback("^(pmute|pban)_", picturequery))
	go u.Dispatcher.AddHandler(handlers.NewCallback("wlcm_", verifyquery))
	go u.Dispatcher.AddHandler(handlers.NewCallback("rmWarn", warnquery))
}
