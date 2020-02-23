package setting

import (
	"fmt"
	"github.com/PaulSonOfLars/goloc"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"

	"github.com/jumatberkah/antispambot/bot/modules/helpers/caching"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/chat_status"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
)

func setUsername(_ ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat_status.RequireSupergroup(chat, msg) {
		if chat_status.RequireUserAdmin(chat, msg, user.Id) {
			if len(args) != 0 {
				switch strings.ToLower(args[0]) {
				case "true", "on", "yes":
					sql.UpdateUsername(chat.Id, "true", "mute", "-", "true")
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					return err
				case "false", "off", "no":
					sql.UpdateUsername(chat.Id, "false", "mute", "-", "true")
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					return err
				default:
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:2"))
					return err
				}
			} else {
				_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:2"))
				return err
			}
		}
	}
	return nil
}

func setVerify(_ ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat_status.RequireSupergroup(chat, msg) {
		if chat_status.RequireUserAdmin(chat, msg, user.Id) {
			if len(args) != 0 {
				switch strings.ToLower(args[0]) {
				case "true", "on", "yes":
					sql.UpdateVerify(chat.Id, "true", "-", "true")
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					return err
				case "false", "off", "no":
					sql.UpdateVerify(chat.Id, "false", "-", "true")
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					return err
				default:
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:2"))
					return err
				}
			} else {
				_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:2"))
				return err
			}
		}
	}
	return nil
}

func setEnforce(_ ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat_status.RequireSupergroup(chat, msg) {
		if chat_status.RequireUserAdmin(chat, msg, user.Id) {
			if len(args) != 0 {
				switch strings.ToLower(args[0]) {
				case "true", "on", "yes":
					sql.UpdateEnforceGban(chat.Id, "true")
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					return err
				case "false", "off", "no":
					sql.UpdateEnforceGban(chat.Id, "false")
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					return err
				default:
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:2"))
					return err
				}
			} else {
				_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:2"))
				return err
			}
		}
	}
	return nil
}

func setPicture(_ ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat_status.RequireSupergroup(chat, msg) {
		if chat_status.RequireUserAdmin(chat, msg, user.Id) {
			if len(args) != 0 {
				switch strings.ToLower(args[0]) {
				case "true", "on", "yes":
					sql.UpdatePicture(chat.Id, "true", "mute", "-", "true")
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					return err
				case "false", "off", "no":
					sql.UpdatePicture(chat.Id, "false", "mute", "-", "true")
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					return err
				default:
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:2"))
					return err
				}
			} else {
				_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:2"))
				return err
			}
		}
	}
	return nil
}

func setTime(_ ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat_status.RequireSupergroup(chat, msg) {
		if chat_status.RequireUserAdmin(chat, msg, user.Id) {
			if len(args) != 0 {
				match, err := regexp.MatchString("^\\d+[mhd]", strings.ToLower(args[0]))
				if match {
					sql.UpdateSetting(chat.Id, args[0], "true")
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					return err
				}
				_, err = msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:3"))
				return err
			}
			_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:3"))
			return err
		}
	}
	return nil
}

func setNotif(_ ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat_status.RequirePrivate(chat, msg) {
		if len(args) != 0 {
			switch strings.ToLower(args[0]) {
			case "true", "on", "yes":
				sql.UpdateNotification(user.Id, "true")
				_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
				return err
			case "false", "off", "no":
				sql.UpdateNotification(user.Id, "false")
				_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
				return err
			default:
				_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:2"))
				return err
			}
		} else {
			_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:2"))
			return err
		}
	}
	return nil
}

func setLang(b ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if !chat_status.RequireUserAdmin(chat, msg, user.Id) {
		return nil
	}

	var btnLang = function.BuildKeyboardf(
		"data/keyboard/language.json",
		2,
		map[string]string{"1": fmt.Sprint(chat.Id)},
	)

	if len(args) == 0 {
		newMsg := b.NewSendableMessage(chat.Id, "*Available Language(s):*")
		newMsg.ParseMode = parsemode.Markdown
		newMsg.ReplyMarkup = &ext.InlineKeyboardMarkup{InlineKeyboard: &btnLang}
		_, err := newMsg.Send()
		return err
	}

	if !goloc.IsLangSupported(args[0]) {
		_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/language/language.go:58"))
		err_handler.HandleErr(err)
		newMsg := b.NewSendableMessage(chat.Id, "*Available Language(s):*")
		newMsg.ParseMode = parsemode.Markdown
		newMsg.ReplyMarkup = &ext.InlineKeyboardMarkup{InlineKeyboard: &btnLang}
		_, err = newMsg.Send()
		return err
	}

	_, _ = caching.REDIS.Set(fmt.Sprintf("lang_%v", chat.Id), args[0], 7200).Result()
	sql.UpdateLang(chat.Id, args[0])

	_, err := msg.ReplyHTML(function.GetStringf(
		chat.Id,
		"modules/language/language.go:51",
		map[string]string{"1": args[0]}),
	)
	return err
}

func adminCache(_ ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if !chat_status.RequireSupergroup(chat, msg) {
		return gotgbot.EndGroups{}
	}

	if !chat_status.RequireUserAdmin(chat, msg, user.Id) {
		return gotgbot.EndGroups{}
	}

	err := caching.CACHE.Delete(fmt.Sprintf("admin_%v", chat.Id))

	if err != nil {
		_, err = msg.ReplyText("Admin cache was refreshed!")
		chat_status.AdminCache(chat)
		return err
	}

	_, err = msg.ReplyText("Admin cache has been updated")
	chat_status.AdminCache(chat)
	return err
}

func LoadSetting(u *gotgbot.Updater) {
	defer logrus.Info("Setting Module Loaded...")
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("username", []rune{'/', '.'}, setUsername))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("verify", []rune{'/', '.'}, setVerify))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("profilepicture", []rune{'/', '.'}, setPicture))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("time", []rune{'/', '.'}, setTime))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("enforce", []rune{'/', '.'}, setEnforce))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("notif", []rune{'/', '.'}, setNotif))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("setlang", []rune{'/', '.'}, setLang))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("admincache", []rune{'/', '.'}, adminCache))
}
