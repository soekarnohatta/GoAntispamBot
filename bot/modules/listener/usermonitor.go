package listener

import (
	"encoding/json"
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
	"github.com/sirupsen/logrus"
	"html"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type casban struct {
	Status bool `json:"ok"`
}

func username(b ext.Bot, u *gotgbot.Update) error {
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

	if chat_status.IsUserAdmin(chat, user.Id) == true {
		return nil
	}

	if chat.Type != "supergroup" {
		return nil
	}

	if user.Username != "" {
		return gotgbot.ContinueGroups{}
	}

	replyText := function.GetStringf(chat.Id, "modules/listener/listener.go:45",
		map[string]string{"1": strconv.Itoa(user.Id), "2": html.EscapeString(user.FirstName),
			"3": db.Action, "4": strconv.Itoa(user.Id)})
	markup := &ext.InlineKeyboardMarkup{}

	if db.Action != "warn" {
		reply := b.NewSendableMessage(chat.Id, replyText)
		reply.ParseMode = parsemode.Html
		reply.ReplyToMessageId = msg.MessageId

		switch db.Action {
		case "mute":
			restrictSend := b.NewSendableRestrictChatMember(chat.Id, user.Id)
			restrictSend.UntilDate = extraction.ExtractTime(b, msg, sql.GetSetting(chat.Id).Time)
			_, err := restrictSend.Send()

			if err != nil {
				if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
					err_handler.HandleTgErr(b, u, err)
					return err
				}
			}

			kb := make([][]ext.InlineKeyboardButton, 1)
			kb[0] = make([]ext.InlineKeyboardButton, 1)
			kb[0][0] = ext.InlineKeyboardButton{Text: function.GetString(chat.Id, "modules/listener/listener.go:51"), CallbackData: fmt.Sprintf("umute_%v_%v", user.Id, chat.Id)}
			markup = &ext.InlineKeyboardMarkup{InlineKeyboard: &kb}
		case "kick":
			_, err := b.UnbanChatMember(chat.Id, user.Id)

			if err != nil {
				if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
					err_handler.HandleTgErr(b, u, err)
					return err
				}
			}
			markup = nil
		case "ban":
			restrictSend := b.NewSendableKickChatMember(chat.Id, user.Id)
			restrictSend.UntilDate = -1
			_, err := restrictSend.Send()

			if err != nil {
				if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
					err_handler.HandleTgErr(b, u, err)
					return err
				}
			}

			kbk := make([][]ext.InlineKeyboardButton, 1)
			kbk[0] = make([]ext.InlineKeyboardButton, 1)
			kbk[0][0] = ext.InlineKeyboardButton{Text: function.GetString(chat.Id, "modules/listener/listener.go:56"), CallbackData: fmt.Sprintf("uba_%v_%v", user.Id, chat.Id)}
			markup = &ext.InlineKeyboardMarkup{InlineKeyboard: &kbk}
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

		if notif != nil && notif.Notification == "true" {
			txt := function.GetStringf(user.Id, "unamep",
				map[string]string{"1": strconv.Itoa(user.Id), "2": user.FirstName, "3": db.Action,
					"4": strconv.Itoa(user.Id), "5": chat.Title})
			reply.Text = txt
			reply.ReplyToMessageId = 0
			reply.ChatId = user.Id
			_, err := reply.Send()
			err_handler.HandleErr(err)
		}
	} else {
		limit := sql.GetWarnSetting(strconv.Itoa(chat.Id))
		warns, _ := sql.WarnUser(strconv.Itoa(user.Id), strconv.Itoa(chat.Id), "Username")

		if warns >= limit {
			go sql.ResetWarns(strconv.Itoa(user.Id), strconv.Itoa(chat.Id))
			val := map[string]string{"1": strconv.Itoa(user.Id), "2": user.FirstName, "3": strconv.Itoa(warns), "4": strconv.Itoa(limit), "5": strconv.Itoa(user.Id)}
			replyText = function.GetStringf(msg.Chat.Id, "modules/warn2", val)
			_, err := chat.UnbanMember(user.Id)

			if err != nil {
				if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
					err_handler.HandleTgErr(b, u, err)
					return err
				}
			}
		} else {
			kb := make([][]ext.InlineKeyboardButton, 1)
			kb[0] = make([]ext.InlineKeyboardButton, 1)
			kb[0][0] = ext.InlineKeyboardButton{Text: function.GetString(chat.Id, "rmwarn"),
				CallbackData: fmt.Sprintf("rmWarn(%v)", user.Id)}
			markup = &ext.InlineKeyboardMarkup{InlineKeyboard: &kb}
			val := map[string]string{"1": strconv.Itoa(user.Id), "2": user.FirstName, "3": strconv.Itoa(warns), "4": strconv.Itoa(limit), "5": strconv.Itoa(user.Id)}
			replyText = function.GetStringf(msg.Chat.Id, "modules/warn", val)
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
		if err != nil {
			if err.Error() == "Bad Request: message can't be deleted" {
				err_handler.HandleTgErr(b, u, err)
			}
		}
	}
	err := logger.SendLog(b, u, "username", "")
	return err
}

func picture(b ext.Bot, u *gotgbot.Update) error {
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

	if chat_status.IsUserAdmin(chat, user.Id) == true {
		return nil
	}

	if chat.Type != "supergroup" {
		return nil
	}

	photo, _ := user.GetProfilePhotos(0, 0)

	if photo != nil && photo.TotalCount != 0 {
		return gotgbot.ContinueGroups{}
	}

	val := map[string]string{"1": strconv.Itoa(user.Id), "2": user.FirstName, "3": db.Action, "4": strconv.Itoa(user.Id)}
	replyText := function.GetStringf(msg.Chat.Id, "modules/listener/listener.go:173", val)
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

			if err != nil {
				if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
					err_handler.HandleTgErr(b, u, err)
					return err
				}
			}

			kb := make([][]ext.InlineKeyboardButton, 1)
			kb[0] = make([]ext.InlineKeyboardButton, 1)
			kb[0][0] = ext.InlineKeyboardButton{Text: function.GetString(chat.Id, "modules/listener/listener.go:179"),
				CallbackData: fmt.Sprintf("pmute_%v_%v", user.Id, chat.Id)}

			markup = &ext.InlineKeyboardMarkup{InlineKeyboard: &kb}
		case "kick":
			_, err := b.UnbanChatMember(chat.Id, user.Id)

			if err != nil {
				if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
					err_handler.HandleTgErr(b, u, err)
					return err
				}
			}

			markup = nil
		case "ban":
			restrictSend := b.NewSendableKickChatMember(chat.Id, user.Id)
			restrictSend.UntilDate = -1
			_, err := restrictSend.Send()

			if err != nil {
				if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
					err_handler.HandleTgErr(b, u, err)
					return err
				}
			}

			kbk := make([][]ext.InlineKeyboardButton, 1)
			kbk[0] = make([]ext.InlineKeyboardButton, 1)
			kbk[0][0] = ext.InlineKeyboardButton{Text: function.GetString(chat.Id, "modules/listener/listener.go:184"), CallbackData: fmt.Sprintf("pban_%v_%v", user.Id, chat.Id)}
			markup = &ext.InlineKeyboardMarkup{InlineKeyboard: &kbk}
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
		if notif != nil && notif.Notification == "true" {
			txtVal := map[string]string{"1": strconv.Itoa(user.Id), "2": user.FirstName, "3": db.Action, "4": strconv.Itoa(user.Id), "5": chat.Title}
			txt := function.GetStringf(user.Id, "picturep", txtVal)
			reply.Text = txt
			reply.ReplyToMessageId = 0
			reply.ChatId = user.Id
			_, err := reply.Send()
			err_handler.HandleErr(err)
		}
	} else {
		limit := sql.GetWarnSetting(strconv.Itoa(chat.Id))
		warns, _ := sql.WarnUser(strconv.Itoa(user.Id), strconv.Itoa(chat.Id), "No Profile Picture")

		var reply = ""
		if warns >= limit {
			go sql.ResetWarns(strconv.Itoa(user.Id), strconv.Itoa(chat.Id))
			val = map[string]string{"1": strconv.Itoa(user.Id), "2": user.FirstName, "3": strconv.Itoa(warns), "4": strconv.Itoa(limit), "5": strconv.Itoa(user.Id)}
			reply = function.GetStringf(msg.Chat.Id, "modules/warn2", val)
			_, err := chat.UnbanMember(user.Id)
			if err != nil {
				if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
					err_handler.HandleTgErr(b, u, err)
					return err
				}
			}
		} else {
			kb := make([][]ext.InlineKeyboardButton, 1)
			kb[0] = make([]ext.InlineKeyboardButton, 1)
			kb[0][0] = ext.InlineKeyboardButton{Text: function.GetString(chat.Id, "rmwarn"),
				CallbackData: fmt.Sprintf("rmWarn(%v)", user.Id)}
			markup = &ext.InlineKeyboardMarkup{InlineKeyboard: &kb}
			val = map[string]string{"1": strconv.Itoa(user.Id), "2": user.FirstName, "3": strconv.Itoa(warns), "4": strconv.Itoa(limit), "5": strconv.Itoa(user.Id)}
			reply = function.GetStringf(msg.Chat.Id, "modules/warn3", val)
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
		if err != nil {
			if err.Error() == "Bad Request: message can't be deleted" {
				err_handler.HandleTgErr(b, u, err)
			}
		}
	}
	err := logger.SendLog(b, u, "picture", "")
	return err
}

func verify(b ext.Bot, u *gotgbot.Update) error {
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

	if chat_status.IsUserAdmin(chat, user.Id) == true {
		return nil
	}

	replyText := function.GetStringf(msg.Chat.Id, "modules/listener/listener.go:298",
		map[string]string{"1": strconv.Itoa(user.Id), "2": user.FirstName, "3": chat.Title, "4": strconv.Itoa(user.Id)})

	kb := make([][]ext.InlineKeyboardButton, 1)
	kb[0] = make([]ext.InlineKeyboardButton, 1)
	kb[0][0] = ext.InlineKeyboardButton{Text: function.GetString(chat.Id, "modules/listener/listener.go:303"),
		CallbackData: fmt.Sprintf("wlcm_%v", user.Id)}

	restrictSend := b.NewSendableRestrictChatMember(chat.Id, user.Id)
	restrictSend.UntilDate = extraction.ExtractTime(b, msg, sql.GetSetting(chat.Id).Time)
	_, err := restrictSend.Send()

	if err != nil {
		if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
			err_handler.HandleTgErr(b, u, err)
			return err
		}
	}

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
		if err != nil {
			if err.Error() == "Bad Request: message can't be deleted" {
				err_handler.HandleTgErr(b, u, err)
			}
		}
	}

	err = logger.SendLog(b, u, "welcome", "")
	err_handler.HandleErr(err)
	return err
}

func spam(b ext.Bot, u *gotgbot.Update) error {
	user := u.EffectiveUser
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	if chat_status.IsUserAdmin(chat, msg.From.Id) == true {
		return gotgbot.ContinueGroups{}
	}

	if chat.Type != "supergroup" {
		return gotgbot.ContinueGroups{}
	}

	if sql.GetEnforceGban(chat.Id).Option != "true" {
		return gotgbot.ContinueGroups{}
	}

	ban := sql.GetUserSpam(user.Id)

	if ban != nil {
		err := spamFunc(b, u)
		err_handler.HandleErr(err)
		err = logger.SendLog(b, u, "spam", ban.Reason)
		return gotgbot.EndGroups{}
	} else {
		ret, _ := casListener(b, u)
		if ret {
			return gotgbot.EndGroups{}
		}
	}

	return gotgbot.ContinueGroups{}
}

func removeLink(b ext.Bot, u *gotgbot.Update) error {
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
			replyText := fmt.Sprintf("Deleted message from %v\nReason: Link", user.FirstName)
			reply := b.NewSendableMessage(chat.Id, replyText)
			_, _ = reply.Send()
		}
	} else if db.Forward == "true" {
		if msg.ForwardFrom != nil || msg.ForwardFromChat != nil {
			_, err := msg.Delete()
			err_handler.HandleErr(err)
			replyText := fmt.Sprintf("Deleted message from %v\nReason: Forwarded Message", user.FirstName)
			reply := b.NewSendableMessage(chat.Id, replyText)
			_, _ = reply.Send()
		}
	}
	return gotgbot.ContinueGroups{}
}

func update(_ ext.Bot, u *gotgbot.Update) error {
	user := u.EffectiveUser
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	if msg != nil {
		if sql.GetChat(chat.Id) == nil {
			go sql.UpdateChat(strconv.Itoa(chat.Id), chat.Title, chat.Type, chat.InviteLink)
		}
		if sql.GetUser(user.Id) == nil {
			go sql.UpdateUser(user.Id, user.Username, user.FirstName, user.LastName)
		}
		if sql.GetLang(chat.Id) == nil {
			go caching.REDIS.Set(fmt.Sprintf("lang_%v", chat.Id), "id", 7200)
			go sql.UpdateLang(chat.Id, "id")
		}
		if msg.ForwardFrom != nil {
			usr := msg.ForwardFrom
			if sql.GetUser(usr.Id) == nil {
				go sql.UpdateUser(usr.Id, usr.Username, usr.FirstName, usr.LastName)
			}
		}

		if chat.Type == "supergroup" {
			if sql.GetVerify(chat.Id) == nil {
				go sql.UpdateVerify(chat.Id, "true", "-", "true")
			}
			if sql.GetUsername(chat.Id) == nil {
				go sql.UpdateUsername(chat.Id, "true", "mute", "-", "true")
			}
			if sql.GetPicture(chat.Id) == nil {
				go sql.UpdatePicture(chat.Id, "true", "mute", "-", "true")
			}
			if sql.GetSetting(chat.Id) == nil {
				go sql.UpdateSetting(chat.Id, "5m", "true")
			}
			if sql.GetEnforceGban(chat.Id) == nil {
				go sql.UpdateEnforceGban(chat.Id, "true")
			}
		}

		if sql.GetNotification(user.Id) == nil {
			go sql.UpdateNotification(user.Id, "true")

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
				cb.Text = function.GetString(chat.Id, "modules/listener/listener.go:454")
				cb.ShowAlert = true
				cb.CacheTime = 10
				_, err := cb.Send()
				return err
			}

			if chat.Type == "supergroup" {
				_, err := b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id,
					"modules/listener/listener.go:441"), true)
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
				_, err := b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id,
					"modules/listener/listener.go:441"), true)
				err_handler.HandleErr(err)
				_, err = msg.Message.Delete()
				err_handler.HandleErr(err)
				cid, _ := strconv.Atoi(strings.Split(splt, "_")[1])
				_, err = b.UnRestrictChatMember(cid, user.Id)
				return err
			}
		}
		_, err := b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id,
			"modules/listener/listener.go:458"), true)
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
				_, err := b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id, "modules/listener/listener.go:511"), true)
				return err
			}

			if chat.Type == "supergroup" {
				_, err := b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id, "modules/listener/listener.go:441"), true)
				err_handler.HandleErr(err)
				_, err = msg.Message.Delete()
				err_handler.HandleErr(err)
				_, err = b.UnRestrictChatMember(chat.Id, user.Id)
				err_handler.HandleErr(err)
				return err
			} else {
				_, err := b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id, "modules/listener/listener.go:441"), true)
				err_handler.HandleErr(err)
				_, err = msg.Message.Delete()
				err_handler.HandleErr(err)
				cid, _ := strconv.Atoi(strings.Split(splt, "_")[1])
				_, err = b.UnRestrictChatMember(cid, user.Id)
				err_handler.HandleErr(err)
				return err
			}
		}
		_, err := b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id, "modules/listener/listener.go:515"), true)
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
				err_handler.HandleErr(err)
				return err
			}
		} else if chat.Type == "private" {
			_, err := b.AnswerCallbackQuery(msg.Id)
			err_handler.HandleErr(err)
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
			_, err := b.AnswerCallbackQueryText(msg.Id, function.GetString(chat.Id, "modules/listener/listener.go:552"), true)
			err_handler.HandleErr(err)
			_, err = msg.Message.Delete()
			err_handler.HandleErr(err)
			_, err = b.UnRestrictChatMember(chat.Id, user.Id)
			return err

		}
		_, err := b.AnswerCallbackQueryText(msg.Id,
			function.GetString(chat.Id, "modules/listener/listener.go:566"), true)
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
		userId := pattern.FindStringSubmatch(query.Data)[1]
		chat := u.EffectiveChat
		res := sql.RemoveWarn(userId, strconv.Itoa(chat.Id))
		if res {
			msg := bot.NewSendableEditMessageText(chat.Id, u.EffectiveMessage.MessageId,
				fmt.Sprintf("Warn removed by admin %v.", html.EscapeString(user.FirstName)))
			msg.ParseMode = parsemode.Html
			_, err := msg.Send()
			return err
		}
		_, err := u.EffectiveMessage.EditText("User already has no warns.")
		return err
	}
	return gotgbot.ContinueGroups{}
}

func spamFunc(b ext.Bot, u *gotgbot.Update) error {
	user := u.EffectiveUser
	chat := u.EffectiveChat
	msg := u.EffectiveMessage
	db := sql.GetSetting(chat.Id)
	if db == nil {
		return nil
	}

	val := map[string]string{"1": strconv.Itoa(user.Id), "2": user.FirstName, "3": strconv.Itoa(user.Id)}
	txtBan := function.GetStringf(chat.Id, "modules/listener/listener.go:580", val)

	_, err := msg.ReplyHTML(txtBan)
	if err != nil {
		if err.Error() == "Bad Request: reply message not found" {
			_, err = b.SendMessageHTML(chat.Id, txtBan)
			return err
		}
	}
	restrictSend := b.NewSendableKickChatMember(chat.Id, user.Id)
	restrictSend.UntilDate = -1
	_, err = restrictSend.Send()
	if err != nil {
		if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
			txtBan = function.GetStringf(chat.Id, "modules/listener/listener.go:warnspam", val)
			_, err := msg.ReplyHTML(txtBan)
			if err != nil {
				if err.Error() == "Bad Request: reply message not found" {
					_, err = b.SendMessageHTML(chat.Id, txtBan)
					return err
				}
			}
			return err
		}
	}

	if db.Deletion == "true" {
		_, err = b.DeleteMessage(chat.Id, msg.MessageId)
		if err != nil {
			if err.Error() == "Bad Request: message can't be deleted" {
				err_handler.HandleTgErr(b, u, err)
				return err
			}
		}
	}
	return nil
}

func casListener(b ext.Bot, u *gotgbot.Update) (bool, error) {
	user := u.EffectiveUser
	spam := &casban{}
	response, err := http.Get(fmt.Sprintf("https://combot.org/api/cas/check?user_id=%v", user.Id))
	err_handler.HandleErr(err)
	defer response.Body.Close()
	if response != nil {
		body, _ := ioutil.ReadAll(response.Body)
		_ = json.Unmarshal(body, &spam)
		if spam.Status == true {
			err = spamFunc(b, u)
			err_handler.HandleErr(err)
			err = logger.SendLog(b, u, "spam", "CAS Banned (Powered By CAS)")
			return true, err
		}
		return false, gotgbot.ContinueGroups{}
	}
	return false, gotgbot.ContinueGroups{}
}

func LoadUserListener(u *gotgbot.Updater) {
	defer logrus.Info("Listeners Module Loaded...")
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, update))
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, spam))
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, username))
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, picture))
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.NewChatMembers(), verify))
	u.Dispatcher.AddHandler(handlers.NewCallback("^(umute|uba)_", usernameQuery))
	u.Dispatcher.AddHandler(handlers.NewCallback("^(pmute|pban)_", pictureQuery))
	u.Dispatcher.AddHandler(handlers.NewCallback("wlcm_", verifyQuery))
	u.Dispatcher.AddHandler(handlers.NewCallback("rmWarn", warnQuery))
}
