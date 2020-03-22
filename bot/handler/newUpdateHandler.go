package handler

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"

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
	doUpdateLog(u)
	doUpdateSetting(u)
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

func doUpdateLog(u *gotgbot.Update) {
	msg := u.EffectiveMessage
	user := u.EffectiveChat
	chat := u.EffectiveChat

	if msg != nil {
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
}

func doUpdateSetting(u *gotgbot.Update) {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat

	if msg != nil {
		if chat.Type == "supergroup" {
			if settingsService.FindGroupSetting(chat.Id) == nil {
				settingStruct := model.GroupSetting{
					ChatID:         chat.Id,
					Gban:           true,
					Username:       true,
					ProfilePicture: true,
				}
				settingsService.UpdateGroupSetting(settingStruct)
			}
		}
	}
}
