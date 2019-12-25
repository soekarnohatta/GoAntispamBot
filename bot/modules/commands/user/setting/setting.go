package setting

import (
	"fmt"
	"github.com/PaulSonOfLars/goloc"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/caching"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/chat_status"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

func setUsername(_ ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat_status.RequireSupergroup(chat, msg) == true {
		if chat_status.RequireUserAdmin(chat, msg, user.Id) {
			if len(args) != 0 {
				if strings.ToLower(args[0]) == "true" {
					go sql.UpdateUsername(chat.Id, "true", "mute", "-", "true")
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					return err
				} else if strings.ToLower(args[0]) == "false" {
					go sql.UpdateUsername(chat.Id, "false", "mute", "-", "true")
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					return err
				} else {
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
	return nil
}

func setVerify(_ ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat_status.RequireSupergroup(chat, msg) == true {
		if chat_status.RequireUserAdmin(chat, msg, user.Id) {
			if len(args) != 0 {
				if strings.ToLower(args[0]) == "true" {
					go sql.UpdateVerify(chat.Id, "true", "-", "true")
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					return err
				} else if strings.ToLower(args[0]) == "false" {
					go sql.UpdateVerify(chat.Id, "false", "-", "true")
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					return err
				} else {
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

	if chat_status.RequireSupergroup(chat, msg) == true {
		if chat_status.RequireUserAdmin(chat, msg, user.Id) {
			if len(args) != 0 {
				if strings.ToLower(args[0]) == "true" {
					go sql.UpdateEnforceGban(chat.Id, "true")
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					return err
				} else if strings.ToLower(args[0]) == "false" {
					go sql.UpdateEnforceGban(chat.Id, "false")
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					return err
				} else {
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

	if chat_status.RequireSupergroup(chat, msg) == true {
		if chat_status.RequireUserAdmin(chat, msg, user.Id) {
			if len(args) != 0 {
				if strings.ToLower(args[0]) == "true" {
					go sql.UpdatePicture(chat.Id, "true", "mute", "-", "true")
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					return err
				} else if strings.ToLower(args[0]) == "false" {
					go sql.UpdatePicture(chat.Id, "false", "mute", "-", "true")
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					return err
				} else {
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

	if chat_status.RequireSupergroup(chat, msg) == true {
		if chat_status.RequireUserAdmin(chat, msg, user.Id) {
			if len(args) != 0 {
				match, err := regexp.MatchString("^\\d+[mhd]", strings.ToLower(args[0]))
				if match == true {
					go sql.UpdateSetting(chat.Id, args[0], "true")
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

	if chat_status.RequirePrivate(chat, msg) == true {
		if len(args) != 0 {
			if strings.ToLower(args[0]) == "true" {
				go sql.UpdateNotification(user.Id, "true")
				_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
				return err
			} else if strings.ToLower(args[0]) == "false" {
				go sql.UpdateNotification(user.Id, "false")
				_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
				return err
			} else {
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

	if chat_status.RequireUserAdmin(chat, msg, user.Id) == false {
		return nil
	}

	if len(args) == 0 {
		_, err := msg.ReplyText("Please insert the language code so that i can change your language")
		return err
	}

	if !goloc.IsLangSupported(args[0]) {
		_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/language/language.go:58"))
		return err
	}

	_, err := caching.REDIS.Set(fmt.Sprintf("lang_%v", chat.Id), args[0], 7200).Result()
	if err != nil {
		go sql.UpdateLang(chat.Id, args[0])
	}

	_, err = msg.ReplyHTML(function.GetStringf(chat.Id, "modules/language/language.go:51",
		map[string]string{"1": args[0]}))
	return err
}

func adminCache(b ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat_status.RequireSupergroup(chat, msg) == false {
		return gotgbot.EndGroups{}
	}

	if chat_status.RequireUserAdmin(chat, msg, user.Id) == false {
		return gotgbot.EndGroups{}
	}

	err := caching.CACHE.Delete(fmt.Sprintf("admin_%v", chat.Id))

	if err != nil {
		_, err = msg.ReplyText("Admin cache was refreshed!")
		go chat_status.AdminCache(chat)
		return err
	}

	_, err = msg.ReplyText("Admin cache has been updated")
	go chat_status.AdminCache(chat)
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
