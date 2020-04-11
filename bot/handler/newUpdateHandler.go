package handler

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"

	"GoAntispamBot/bot/helpers/trans"
	"GoAntispamBot/bot/model"
	"GoAntispamBot/bot/providers/antiSpamProvider"
	"GoAntispamBot/bot/providers/telegramProvider"
	"GoAntispamBot/bot/services/logService"
	"GoAntispamBot/bot/services/settingsService"
)

type Handler struct {
	TelegramProvider telegramProvider.TelegramProvider
}

func (r Handler) UpdateChat(b ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage
	if msg != nil {
		doUpdateLog(u)
		doUpdateSetting(u)
	}
	return gotgbot.ContinueGroups{}
}

func (r Handler) UsernameHandler(_ ext.Bot, u *gotgbot.Update) error {
	r.TelegramProvider.Init(u)
	msg := u.EffectiveMessage
	if msg != nil {
		if settingsService.FindGroupSetting(msg.Chat.Id).Username {
			if msg.From.Username == "" {
				_, _ = msg.Delete()
				go r.TelegramProvider.SendText(
					trans.GetStringf(msg.Chat.Id,
						"actions/handler",
						map[string]string{"1": ""}),
					0,
					0,
					nil,
				)
			}
		}
	}
	return gotgbot.ContinueGroups{}
}

func (r Handler) PictureHandler(_ ext.Bot, u *gotgbot.Update) error {
	r.TelegramProvider.Init(u)
	msg := u.EffectiveMessage
	if msg != nil {
		if settingsService.FindGroupSetting(msg.Chat.Id).ProfilePicture {
			pic, _ := msg.From.GetProfilePhotos(0, 1)
			if pic != nil && !(pic.TotalCount >= 1) {
				_, _ = msg.Delete()
				go r.TelegramProvider.SendText(
					trans.GetStringf(msg.Chat.Id,
						"actions/handler",
						map[string]string{"1": ""}),
					0,
					0,
					nil,
				)
			}
		}
	}
	return gotgbot.ContinueGroups{}
}

func (r Handler) LockHandler(_ ext.Bot, u *gotgbot.Update) error {
	r.TelegramProvider.Init(u)
	msg := u.EffectiveMessage
	if msg != nil {
		//
	}
	return gotgbot.ContinueGroups{}
}

func (r Handler) GbanHandler(_ ext.Bot, u *gotgbot.Update) error {
	r.TelegramProvider.Init(u)
	msg := u.EffectiveMessage
	if msg != nil {
		antiSpamProvider.FilterSpamUser(r.TelegramProvider)
	}
	return gotgbot.ContinueGroups{}
}

//--------------------------------------------------------------------------------------------------------------------//

func doUpdateLog(u *gotgbot.Update) {
	user := u.EffectiveChat
	chat := u.EffectiveChat

	if chat.Type == "private" {
		userStruct := model.UserLog{
			UserID:    user.Id,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			UserName:  user.Username,
		}

		logService.UpdateUser(userStruct)
	}

	chatStruct := model.ChatLog{
		ChatID:    chat.Id,
		ChatLink:  chat.InviteLink,
		ChatTitle: chat.Title,
		ChatType:  chat.Type,
	}
	logService.UpdateChat(chatStruct)
}

func doUpdateSetting(u *gotgbot.Update) {
	chat := u.EffectiveChat

	if chat.Type == "supergroup" {
		if settingsService.FindGroupSetting(chat.Id) == nil {
			settingStruct := model.GroupSetting{
				ChatID:         chat.Id,
				Gban:           true,
				Username:       true,
				ProfilePicture: true,
				Time:           5,
			}
			settingsService.UpdateGroupSetting(settingStruct)
		}
	}
}
