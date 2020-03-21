package handlers

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"

	"GoAntispamBot/bot/model"
	"GoAntispamBot/bot/providers"
	"GoAntispamBot/bot/services"
)

type UpdateHandler model.Message

func (r UpdateHandler) UpdateChat(b ext.Bot, u *gotgbot.Update) error {
	doUpdateLog(u)
	doUpdateSetting(u)
	return gotgbot.ContinueGroups{}
}

func (r UpdateHandler) GbanHandler(_ ext.Bot, u *gotgbot.Update) error {
	r.TelegramProvider.Init(u)
	msg := u.EffectiveMessage
	if msg != nil {
		providers.FilterSpamUser(r.TelegramProvider)
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

			services.UpdateUser(userStruct)
		}

		chatStruct := model.ChatLog{
			ChatID:    chat.Id,
			ChatLink:  chat.InviteLink,
			ChatTitle: chat.Title,
			ChatType:  chat.Type,
		}
		services.UpdateChat(chatStruct)
	}
}

func doUpdateSetting(u *gotgbot.Update) {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat

	if msg != nil {
		if chat.Type == "supergroup" {
			if services.FindGroupSetting(chat.Id) == nil {
				settingStruct := model.GroupSetting{
					ChatID:         chat.Id,
					Gban:           true,
					Username:       true,
					ProfilePicture: true,
				}
				services.UpdateGroupSetting(settingStruct)
			}
		}
	}
}
