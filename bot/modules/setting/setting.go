package setting

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/caching"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/chat_status"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

func setusername(_ ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat_status.RequireSupergroup(chat, msg) == true {
		if chat_status.RequireUserAdmin(chat, msg, user.Id) {
			if len(args) != 0 {
				if strings.ToLower(args[0]) == "true" {
					err := sql.UpdateUsername(chat.Id, "true", "mute", "-", "true")
					err_handler.HandleErr(err)
					_, err = msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					err_handler.HandleErr(err)
					return err
				} else if strings.ToLower(args[0]) == "false" {
					err := sql.UpdateUsername(chat.Id, "false", "mute", "-", "true")
					err_handler.HandleErr(err)
					_, err = msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					return err
				} else {
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:2"))
					err_handler.HandleErr(err)
					return err
				}
			} else {
				_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:2"))
				err_handler.HandleErr(err)
				return err
			}
		}
	}
	return nil
}

func setverify(_ ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat_status.RequireSupergroup(chat, msg) == true {
		if chat_status.RequireUserAdmin(chat, msg, user.Id) {
			if len(args) != 0 {
				if strings.ToLower(args[0]) == "true" {
					err := sql.UpdateVerify(chat.Id, "true", "-", "true")
					err_handler.HandleErr(err)
					_, err = msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					err_handler.HandleErr(err)
					return err
				} else if strings.ToLower(args[0]) == "false" {
					err := sql.UpdateVerify(chat.Id, "false", "-", "true")
					err_handler.HandleErr(err)
					_, err = msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					return err
				} else {
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:2"))
					err_handler.HandleErr(err)
					return err
				}
			} else {
				_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:2"))
				err_handler.HandleErr(err)
				return err
			}
		}
	}
	return nil
}

func setenforce(_ ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat_status.RequireSupergroup(chat, msg) == true {
		if chat_status.RequireUserAdmin(chat, msg, user.Id) {
			if len(args) != 0 {
				if strings.ToLower(args[0]) == "true" {
					db := make(chan error)
					go func() { db <- sql.UpdateEnforceGban(chat.Id, "true") }()
					err_handler.HandleErr(<-db)
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					err_handler.HandleErr(err)
					return err
				} else if strings.ToLower(args[0]) == "false" {
					db := make(chan error)
					go func() { db <- sql.UpdateEnforceGban(chat.Id, "false") }()
					err_handler.HandleErr(<-db)
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					err_handler.HandleErr(err)
					return err
				} else {
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:2"))
					err_handler.HandleErr(err)
					return err
				}
			} else {
				_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:2"))
				err_handler.HandleErr(err)
				return err
			}
		}
	}
	return nil
}

func setpicture(_ ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat_status.RequireSupergroup(chat, msg) == true {
		if chat_status.RequireUserAdmin(chat, msg, user.Id) {
			if len(args) != 0 {
				if strings.ToLower(args[0]) == "true" {
					db := make(chan error)
					go func() { db <- sql.UpdatePicture(chat.Id, "true", "mute", "-", "true") }()
					err_handler.HandleErr(<-db)
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					err_handler.HandleErr(err)
					return err
				} else if strings.ToLower(args[0]) == "false" {
					db := make(chan error)
					go func() { db <- sql.UpdatePicture(chat.Id, "false", "mute", "-", "true") }()
					err_handler.HandleErr(<-db)
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					err_handler.HandleErr(err)
					return err
				} else {
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:2"))
					err_handler.HandleErr(err)
					return err
				}
			} else {
				_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:2"))
				err_handler.HandleErr(err)
				return err
			}
		}
	}
	return nil
}

func settime(_ ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat_status.RequireSupergroup(chat, msg) == true {
		if chat_status.RequireUserAdmin(chat, msg, user.Id) {
			if len(args) != 0 {
				match, err := regexp.MatchString("^\\d+[mhd]", strings.ToLower(args[0]))
				if match == true {
					db := make(chan error)
					go func() { db <- sql.UpdateSetting(chat.Id, args[0], "true") }()
					err_handler.HandleErr(<-db)
					_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
					err_handler.HandleErr(err)
					return err
				}
				_, err = msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:3"))
				err_handler.HandleErr(err)
				return err
			}
			_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:3"))
			err_handler.HandleErr(err)
			return err
		}
	}
	return nil
}

func setnotif(_ ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat_status.RequirePrivate(chat, msg) == true {
		if len(args) != 0 {
			if strings.ToLower(args[0]) == "true" {
				db := make(chan error)
				go func() { db <- sql.UpdateNotification(user.Id, "true") }()
				err_handler.HandleErr(<-db)
				_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
				err_handler.HandleErr(err)
				return err
			} else if strings.ToLower(args[0]) == "false" {
				db := make(chan error)
				go func() { db <- sql.UpdateNotification(user.Id, "false") }()
				err_handler.HandleErr(<-db)
				_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:1"))
				err_handler.HandleErr(err)
				return err
			} else {
				_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:2"))
				err_handler.HandleErr(err)
				return err
			}
		} else {
			_, err := msg.ReplyHTML(function.GetString(chat.Id, "modules/setting/setting.go:2"))
			err_handler.HandleErr(err)
			return err
		}
	}
	return nil
}

func admincache(b ext.Bot, u *gotgbot.Update) error {
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
		err_handler.HandleTgErr(b, u, err)
		return err
	}
	_, err = msg.ReplyHTML("Admin cache has been cleaned")
	err_handler.HandleErr(err)
	return err
}

func LoadSetting(u *gotgbot.Updater) {
	defer logrus.Info("Setting Module Loaded...")
	go u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("username", []rune{'/', '.'}, setusername))
	go u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("verify", []rune{'/', '.'}, setverify))
	go u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("profilepicture", []rune{'/', '.'}, setpicture))
	go u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("time", []rune{'/', '.'}, settime))
	go u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("enforce", []rune{'/', '.'}, setenforce))
	go u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("notif", []rune{'/', '.'}, setnotif))
	go u.Dispatcher.AddHandler(handlers.NewPrefixCommand("admincache", []rune{'/', '.'}, admincache))
}
