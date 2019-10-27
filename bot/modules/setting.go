package modules

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/jumatberkah/antispambot/bot/helpers/chat_status"
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
	"regexp"
	"strings"
)

func setusername(_ ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat.Type == "supergroup" {
		if chat_status.RequireUserAdmin(chat, msg, user.Id, nil) {
			if len(args) != 0 {
				if strings.ToLower(args[0]) == "true" {
					db := make(chan error)
					go func() { db <- sql.UpdateUsername(chat.Id, "true", "mute", "-", "true") }()
					err_handler.HandleErr(<-db)
					_, err := msg.ReplyHTML(GetString(chat.Id, "modules/setting.go:1"))
					return err
				} else if strings.ToLower(args[0]) == "false" {
					db := make(chan error)
					go func() { db <- sql.UpdateUsername(chat.Id, "false", "mute", "-", "true") }()
					err_handler.HandleErr(<-db)
					_, err := msg.ReplyHTML(GetString(chat.Id, "modules/setting.go:1"))
					return err
				} else {
					_, err := msg.ReplyHTML(GetString(chat.Id, "modules/setting.go:2"))
					return err
				}
			} else {
				_, err := msg.ReplyHTML(GetString(chat.Id, "modules/setting.go:2"))
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

	if chat.Type == "supergroup" {
		if chat_status.RequireUserAdmin(chat, msg, user.Id, nil) {
			if len(args) != 0 {
				if strings.ToLower(args[0]) == "true" {
					db := make(chan error)
					go func() { db <- sql.UpdateVerify(chat.Id, "true", "-", "true") }()
					err_handler.HandleErr(<-db)
					_, err := msg.ReplyHTML(GetString(chat.Id, "modules/setting.go:1"))
					return err
				} else if strings.ToLower(args[0]) == "false" {
					db := make(chan error)
					go func() { db <- sql.UpdateVerify(chat.Id, "false", "-", "true") }()
					err_handler.HandleErr(<-db)
					_, err := msg.ReplyHTML(GetString(chat.Id, "modules/setting.go:1"))
					return err
				} else {
					_, err := msg.ReplyHTML(GetString(chat.Id, "modules/setting.go:2"))
					return err
				}
			} else {
				_, err := msg.ReplyHTML(GetString(chat.Id, "modules/setting.go:2"))
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

	if chat.Type == "supergroup" {
		if chat_status.RequireUserAdmin(chat, msg, user.Id, nil) {
			if len(args) != 0 {
				if strings.ToLower(args[0]) == "true" {
					db := make(chan error)
					go func() { db <- sql.UpdateEnforceGban(chat.Id, "true") }()
					err_handler.HandleErr(<-db)
					_, err := msg.ReplyHTML(GetString(chat.Id, "modules/setting.go:1"))
					return err
				} else if strings.ToLower(args[0]) == "false" {
					db := make(chan error)
					go func() { db <- sql.UpdateEnforceGban(chat.Id, "false") }()
					err_handler.HandleErr(<-db)
					_, err := msg.ReplyHTML(GetString(chat.Id, "modules/setting.go:1"))
					return err
				} else {
					_, err := msg.ReplyHTML(GetString(chat.Id, "modules/setting.go:2"))
					return err
				}
			} else {
				_, err := msg.ReplyHTML(GetString(chat.Id, "modules/setting.go:2"))
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

	if chat.Type == "supergroup" {
		if chat_status.RequireUserAdmin(chat, msg, user.Id, nil) {
			if len(args) != 0 {
				if strings.ToLower(args[0]) == "true" {
					db := make(chan error)
					go func() { db <- sql.UpdatePicture(chat.Id, "true", "mute", "-", "true") }()
					err_handler.HandleErr(<-db)
					_, err := msg.ReplyHTML(GetString(chat.Id, "modules/setting.go:1"))
					return err
				} else if strings.ToLower(args[0]) == "false" {
					db := make(chan error)
					go func() { db <- sql.UpdatePicture(chat.Id, "false", "mute", "-", "true") }()
					err_handler.HandleErr(<-db)
					_, err := msg.ReplyHTML(GetString(chat.Id, "modules/setting.go:1"))
					return err
				} else {
					_, err := msg.ReplyHTML(GetString(chat.Id, "modules/setting.go:2"))
					return err
				}
			} else {
				_, err := msg.ReplyHTML(GetString(chat.Id, "modules/setting.go:2"))
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

	if chat.Type == "supergroup" {
		if chat_status.RequireUserAdmin(chat, msg, user.Id, nil) {
			if len(args) != 0 {
				match, err := regexp.MatchString("^\\d+[mhd]", strings.ToLower(args[0]))
				if match == true {
					db := make(chan error)
					go func() { db <- sql.UpdateSetting(chat.Id, args[0], "true") }()
					err_handler.HandleErr(<-db)
					_, err := msg.ReplyHTML(GetString(chat.Id, "modules/setting.go:1"))
					return err
				} else {
					_, err = msg.ReplyHTML(GetString(chat.Id, "modules/setting.go:3"))
					return err
				}
			} else {
				_, err := msg.ReplyHTML(GetString(chat.Id, "modules/setting.go:3"))
				return err
			}
		}
	}
	return gotgbot.ContinueGroups{}
}

func LoadSetting(u *gotgbot.Updater) {
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("username", setusername))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("verify", setverify))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("profilepicture", setpicture))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("time", settime))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("enforce", setenforce))
}
