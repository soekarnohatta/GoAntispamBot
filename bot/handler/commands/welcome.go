package commands

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"strconv"
	"strings"

	"GoAntispamBot/bot/handler/events"
	"GoAntispamBot/bot/helpers/chatStatus"
	"GoAntispamBot/bot/helpers/user"
	"GoAntispamBot/bot/model"
	"GoAntispamBot/bot/services/welcomeService"
)

func welcome(bot ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat

	if u.EffectiveChat.Type == "private" {
		_, err := u.EffectiveMessage.ReplyText("This command is meant to be used in a group!")
		return err
	}

	if !chatStatus.IsAdmin(chat, u.EffectiveUser.Id) {
		_, _ = u.EffectiveMessage.ReplyText("You need to be an admin to do this.")
		return gotgbot.ContinueGroups{}
	}

	if len(args) == 0 || strings.ToLower(args[0]) == "noformat" {
		noformat := len(args) > 0 && strings.ToLower(args[0]) == "noformat"
		welcPrefs := welcomeService.GetWelcomePrefs(strconv.Itoa(chat.Id))
		_, _ = u.EffectiveMessage.ReplyHTMLf("I am currently welcoming users: <code>%v</code>"+
			"\nI am currently deleting old welcomes: <code>%v</code>"+
			"\nI am currently deleting service messages: <code>%v</code>"+
			"\nOn joining, I am currently muting users: <code>%v</code>"+
			"\nThe welcome message not filling the {} is:",
			welcPrefs.ShouldWelcome,
			welcPrefs.CleanWelcome != 0,
			welcPrefs.DelJoined,
			welcPrefs.ShouldMute)

		if welcPrefs.WelcomeType == welcomeService.BUTTON_TEXT {
			buttons := welcomeService.GetWelcomeButtons(chat.Id)
			if strings.Contains(welcPrefs.CustomWelcome, "{rules}") {
				rulesButton := model.WelcomeButton{
					Id:       0,
					ChatId:   u.EffectiveChat.Id,
					Name:     "Rules",
					Url:      fmt.Sprintf("t.me/%v?start=%v", bot.UserName, u.EffectiveChat.Id),
					SameLine: false,
				}
				buttons = append(buttons, rulesButton)
				strings.ReplaceAll(welcPrefs.CustomWelcome, "{rules}", "")
			}
			if noformat {
				welcPrefs.CustomWelcome += user.RevertButtons(buttons)
				_, err := u.EffectiveMessage.ReplyHTML(welcPrefs.CustomWelcome)
				return err
			} else {
				keyb := user.BuildWelcomeKeyboard(buttons)
				keyboard := ext.InlineKeyboardMarkup{InlineKeyboard: &keyb}
				events.send(bot, u, welcPrefs.CustomWelcome, &keyboard, welcomeService.DefaultWelcome, !welcPrefs.DelJoined)
			}
		} else {
			_, err := events.EnumFuncMap[welcPrefs.WelcomeType](bot, chat.Id, welcPrefs.CustomWelcome) // needs change
			return err
		}
	} else if len(args) >= 1 {
		switch strings.ToLower(args[0]) {
		case "on", "yes":
			welcomeService.SetWelcPref(strconv.Itoa(chat.Id), true)
			_, err := u.EffectiveMessage.ReplyText("I'll welcome users from now on.")
			return err
		case "off", "no":
			welcomeService.SetWelcPref(strconv.Itoa(chat.Id), false)
			_, err := u.EffectiveMessage.ReplyText("I'll not welcome users from now on.")
			return err
		default:
			_, err := u.EffectiveMessage.ReplyText("I understand 'on/yes' or 'off/no' only!")
			return err
		}
	}
	return nil
}

func setWelcome(_ ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	if u.EffectiveChat.Type == "private" {
		_, err := u.EffectiveMessage.ReplyText("This command is meant to be used in a group!")
		return err
	}

	if !chatStatus.IsAdmin(chat, u.EffectiveUser.Id) {
		_, _ = u.EffectiveMessage.ReplyText("You need to be an admin to do this.")
		return gotgbot.ContinueGroups{}
	}

	text, dataType, content, buttons := user.GetWelcomeType(msg)
	if dataType == -1 {
		_, err := msg.ReplyText("You didn't specify what to reply with!")
		return err
	}

	btns := make([]model.WelcomeButton, len(buttons))
	for i, btn := range buttons {
		btns[i] = model.WelcomeButton{
			ChatId:   chat.Id,
			Name:     btn.Name,
			Url:      btn.Content,
			SameLine: btn.SameLine,
		}
	}

	if text != "" {
		welcomeService.SetCustomWelcome(chat.Id, text, btns, dataType)
	} else {
		welcomeService.SetCustomWelcome(chat.Id, content, btns, dataType)
	}

	_, err := msg.ReplyText("Successfully set custom welcome message!")
	return err
}

func resetWelcome(_ ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat

	if u.EffectiveChat.Type == "private" {
		_, err := u.EffectiveMessage.ReplyText("This command is meant to be used in a group!")
		return err
	}

	if !chatStatus.IsAdmin(chat, u.EffectiveUser.Id) {
		_, _ = u.EffectiveMessage.ReplyText("You need to be an admin to do this.")
		return gotgbot.ContinueGroups{}
	}

	go welcomeService.SetCustomWelcome(strconv.Itoa(chat.Id), welcomeService.DefaultWelcome, nil, welcomeService.TEXT)

	_, err := u.EffectiveMessage.ReplyText("Succesfully reset custom welcome message to default!")
	return err
}

func cleanWelcome(_ ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat

	if u.EffectiveChat.Type == "private" {
		_, err := u.EffectiveMessage.ReplyText("This command is meant to be used in a group!")
		return err
	}

	if !chatStatus.IsAdmin(chat, u.EffectiveUser.Id) {
		_, _ = u.EffectiveMessage.ReplyText("You need to be an admin to do this.")
		return gotgbot.ContinueGroups{}
	}

	if len(args) == 0 {
		cleanPref := welcomeService.GetCleanWelcome(strconv.Itoa(chat.Id))
		if cleanPref != 0 {
			_, err := u.EffectiveMessage.ReplyText("I should be deleting welcome messages up to two days old.")
			return err
		} else {
			_, err := u.EffectiveMessage.ReplyText("I'm currently not deleting old welcome messages!")
			return err
		}
	}

	switch strings.ToLower(args[0]) {
	case "off", "no":
		welcomeService.SetCleanWelcome(strconv.Itoa(chat.Id), 0)
		_, err := u.EffectiveMessage.ReplyText("I'll try to delete old welcome messages!")
		return err
	case "on", "yes":
		welcomeService.SetCleanWelcome(strconv.Itoa(chat.Id), 1)
		_, err := u.EffectiveMessage.ReplyText("I'll try to delete old welcome messages!")
		return err
	default:
		_, err := u.EffectiveMessage.ReplyText("I understand 'on/yes' or 'off/no' only!")
		return err
	}
}

func delJoined(_ ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat

	if u.EffectiveChat.Type == "private" {
		_, err := u.EffectiveMessage.ReplyText("This command is meant to be used in a group!")
		return err
	}

	if !chatStatus.IsAdmin(chat, u.EffectiveUser.Id) {
		_, _ = u.EffectiveMessage.ReplyText("You need to be an admin to do this.")
		return gotgbot.ContinueGroups{}
	}

	if len(args) == 0 {
		delPref := welcomeService.GetDelPref(strconv.Itoa(chat.Id))
		if delPref {
			_, err := u.EffectiveMessage.ReplyMarkdown("I should be deleting `user` joined the chat messages now.")
			return err
		} else {
			_, err := u.EffectiveMessage.ReplyText("I'm currently not deleting joined messages.")
			return err
		}
	}

	switch strings.ToLower(args[0]) {
	case "off", "no":
		welcomeService.SetDelPref(chat.Id, false)
		_, err := u.EffectiveMessage.ReplyText("I won't delete joined messages.")
		return err
	case "on", "yes":
		welcomeService.SetDelPref(chat.Id, true)
		_, err := u.EffectiveMessage.ReplyText("I'll try to delete joined messages!")
		return err
	default:
		_, err := u.EffectiveMessage.ReplyText("I understand 'on/yes' or 'off/no' only!")
		return err
	}
}

func welcomeMute(_ ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat

	if u.EffectiveChat.Type == "private" {
		_, err := u.EffectiveMessage.ReplyText("This command is meant to be used in a group!")
		return err
	}

	if !chatStatus.IsAdmin(chat, u.EffectiveUser.Id) {
		_, _ = u.EffectiveMessage.ReplyText("You need to be an admin to do this.")
		return gotgbot.ContinueGroups{}
	}

	if len(args) == 0 {
		welcPref := welcomeService.GetWelcomePrefs(strconv.Itoa(chat.Id))
		if welcPref.ShouldMute {
			_, err := u.EffectiveMessage.ReplyMarkdown("I'm currently muting users when they join.")
			return err
		} else {
			_, err := u.EffectiveMessage.ReplyText("I'm currently not muting users when they join.")
			return err
		}
	}

	switch strings.ToLower(args[0]) {
	case "off", "no":
		welcomeService.SetMutePref(chat.Id, false)
		_, err := u.EffectiveMessage.ReplyText("I won't mute new users when they join.")
		return err
	case "on", "yes":
		welcomeService.SetMutePref(chat.Id, true)
		_, err := u.EffectiveMessage.ReplyText("I'll try to mute new users when they join!")
		return err
	default:
		_, err := u.EffectiveMessage.ReplyText("I understand 'on/yes' or 'off/no' only!")
		return err
	}
}
