package commands

import (
	"GoAntispamBot/bot/handler/events"
	"GoAntispamBot/bot/helpers/chatStatus"
	"GoAntispamBot/bot/helpers/extraction"
	"GoAntispamBot/bot/helpers/message"
	"GoAntispamBot/bot/helpers/trans"
	"GoAntispamBot/bot/model"
	"GoAntispamBot/bot/providers/telegramProvider"
	"GoAntispamBot/bot/services/welcomeService"
	"fmt"
	"strings"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
)

type CommandWelcome struct {
	TelegramProvider telegramProvider.TelegramProvider
}

func (r CommandWelcome) Welcome(bot ext.Bot, u *gotgbot.Update, args []string) error {
	r.TelegramProvider.Init(u)
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if !chatStatus.RequireSuperGroup(r.TelegramProvider) || !chatStatus.RequireAdmin(user.Id, r.TelegramProvider) {
		return nil
	}

	if len(args) == 0 || strings.ToLower(args[0]) == "noformat" {
		noformat := len(args) > 0 && strings.ToLower(args[0]) == "noformat"
		welcPrefs := welcomeService.GetWelcomePrefs(chat.Id)
		go r.TelegramProvider.ReplyText(
			trans.GetStringf(
				chat.Id, "action/welcomepref",
				map[string]string{
					"1": fmt.Sprint(welcPrefs.ShouldWelcome),
					"2": fmt.Sprint(welcPrefs.CleanWelcome != 0),
					"3": fmt.Sprint(welcPrefs.DelJoined),
					"4": fmt.Sprint(welcPrefs.ShouldMute)}),
		)

		if welcPrefs.WelcomeType == welcomeService.BUTTON_TEXT {
			buttons := welcomeService.GetWelcomeButtons(chat.Id)
			if strings.Contains(welcPrefs.CustomWelcome, "{rules}") {
				rulesButton := model.WelcomeButton{
					ChatId:   u.EffectiveChat.Id,
					Name:     "Rules",
					Url:      fmt.Sprintf("t.me/%v?start=%v", bot.UserName, u.EffectiveChat.Id),
					SameLine: false,
				}
				buttons = append(buttons, rulesButton)
				strings.ReplaceAll(welcPrefs.CustomWelcome, "{rules}", "")
			}
			if noformat {
				welcPrefs.CustomWelcome += message.RevertButtons(buttons)
				go r.TelegramProvider.ReplyText(welcPrefs.CustomWelcome)
				return nil
			}

			keyb := message.BuildWelcomeKeyboard(buttons)
			keyboard := ext.InlineKeyboardMarkup{InlineKeyboard: &keyb}
			events.Send(bot, u, welcPrefs.CustomWelcome, &keyboard, welcomeService.DefaultWelcome, !welcPrefs.DelJoined)

		} else {
			_, err := events.EnumFuncMap[welcPrefs.WelcomeType](bot, chat.Id, welcPrefs.CustomWelcome) // needs change
			return err
		}
	} else if len(args) >= 1 {
		extractBool, _ := extraction.ExtractBool(r.TelegramProvider, strings.ToLower(args[0]))
		welcomeService.SetWelcPref(chat.Id, extractBool)
		go r.TelegramProvider.ReplyText(trans.GetString(chat.Id, "action/updset"))
		return nil

	}
	return nil
}

func (r CommandWelcome) SetWelcome(_ ext.Bot, u *gotgbot.Update) error {
	r.TelegramProvider.Init(u)
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	if !chatStatus.RequireSuperGroup(r.TelegramProvider) || !chatStatus.RequireAdmin(msg.From.Id, r.TelegramProvider) {
		return nil
	}

	text, dataType, content, buttons := message.GetWelcomeType(msg)
	if dataType == -1 {
		go r.TelegramProvider.ReplyText("error/wlcspecifyerror")
		return nil
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

	go r.TelegramProvider.ReplyText(trans.GetString(chat.Id, "action/updwlcm"))
	return nil
}

func (r CommandWelcome) ResetWelcome(_ ext.Bot, u *gotgbot.Update) error {
	r.TelegramProvider.Init(u)
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	if !chatStatus.RequireSuperGroup(r.TelegramProvider) || !chatStatus.RequireAdmin(msg.From.Id, r.TelegramProvider) {
		return nil
	}

	welcomeService.SetCustomWelcome(chat.Id, welcomeService.DefaultWelcome, nil, welcomeService.TEXT)
	go r.TelegramProvider.ReplyText(trans.GetString(chat.Id, "action/updreset"))
	return nil
}

func (r CommandWelcome) CleanWelcome(_ ext.Bot, u *gotgbot.Update, args []string) error {
	r.TelegramProvider.Init(u)
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	if !chatStatus.RequireSuperGroup(r.TelegramProvider) || !chatStatus.RequireAdmin(msg.From.Id, r.TelegramProvider) {
		return nil
	}

	if len(args) == 0 {
		cleanPref := welcomeService.GetCleanWelcome(chat.Id)
		if cleanPref != 0 {
			go r.TelegramProvider.ReplyText(trans.GetString(chat.Id, "actions/upddelwlcm"))
			return nil
		} else {
			go r.TelegramProvider.ReplyText(trans.GetString(chat.Id, "actions/updnotdelwlcm"))
			return nil
		}
	}

	extractBool, _ := extraction.ExtractBool(r.TelegramProvider, strings.ToLower(args[0]))
	switch extractBool {
	case true:
		welcomeService.SetCleanWelcome(chat.Id, 1)
		go r.TelegramProvider.ReplyText(trans.GetString(chat.Id, "actions/upddelwlcm"))
	case false:
		welcomeService.SetCleanWelcome(chat.Id, 0)
		go r.TelegramProvider.ReplyText(trans.GetString(chat.Id, "actions/updnotdelwlcm"))
	}
	return nil
}

func (r CommandWelcome) DelJoined(_ ext.Bot, u *gotgbot.Update, args []string) error {
	r.TelegramProvider.Init(u)
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	if !chatStatus.RequireSuperGroup(r.TelegramProvider) || !chatStatus.RequireAdmin(msg.From.Id, r.TelegramProvider) {
		return nil
	}

	if len(args) == 0 {
		delPref := welcomeService.GetDelPref(chat.Id)
		if delPref {
			_, _ = u.EffectiveMessage.ReplyMarkdown("I should be deleting `user` joined the chat messages now.")
			return nil
		} else {
			_, _ = u.EffectiveMessage.ReplyText("I'm currently not deleting joined messages.")
			return nil
		}
	}

	extractBool, _ := extraction.ExtractBool(r.TelegramProvider, args[0])
	switch extractBool {
	case false:
		welcomeService.SetDelPref(chat.Id, false)
		go r.TelegramProvider.ReplyText(trans.GetString(chat.Id, "actions/upddeljoin"))
		return nil
	case true:
		welcomeService.SetDelPref(chat.Id, true)
		go r.TelegramProvider.ReplyText(trans.GetString(chat.Id, "actions/upddeljoin"))
		return nil
	}
	return nil
}

func (r CommandWelcome) WelcomeMute(_ ext.Bot, u *gotgbot.Update, args []string) error {
	r.TelegramProvider.Init(u)
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	if !chatStatus.RequireSuperGroup(r.TelegramProvider) || !chatStatus.RequireAdmin(msg.From.Id, r.TelegramProvider) {
		return nil
	}

	if len(args) == 0 {
		welcPref := welcomeService.GetWelcomePrefs(chat.Id)
		if welcPref.ShouldMute {
			_, err := u.EffectiveMessage.ReplyMarkdown("I'm currently muting users when they join.")
			return err
		} else {
			_, err := u.EffectiveMessage.ReplyText("I'm currently not muting users when they join.")
			return err
		}
	}

	extractBool, _ := extraction.ExtractBool(r.TelegramProvider, args[0])
	switch extractBool {
	case false:
		welcomeService.SetMutePref(chat.Id, false)
		go r.TelegramProvider.ReplyText(trans.GetString(chat.Id, "actions/updwelcmute"))
		return nil
	case true:
		welcomeService.SetMutePref(chat.Id, true)
		go r.TelegramProvider.ReplyText(trans.GetString(chat.Id, "actions/updwelcmute"))
		return nil
	}
	return nil
}
