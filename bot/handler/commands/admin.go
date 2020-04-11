package commands

import (
	"encoding/json"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"strconv"
	"time"

	"GoAntispamBot/bot/helpers/chatStatus"
	"GoAntispamBot/bot/helpers/errHandler"
	"GoAntispamBot/bot/helpers/extraction"
	"GoAntispamBot/bot/helpers/trans"
	"GoAntispamBot/bot/model"
	"GoAntispamBot/bot/providers/telegramProvider"
	"GoAntispamBot/bot/services/securityService"
)

type CommandAdmin struct {
	TelegramProvider telegramProvider.TelegramProvider
}

func (r CommandAdmin) BanUser(_ ext.Bot, u *gotgbot.Update, args []string) error {
	r.TelegramProvider.Init(u)
	msg := u.EffectiveMessage

	if !chatStatus.RequireSudo(msg.From.Id, r.TelegramProvider) {
		return nil
	}

	userID, reason := extraction.ExtractUserText(r.TelegramProvider, args)
	if userID == 0 {
		go r.TelegramProvider.ReplyText(trans.GetString(msg.Chat.Id, "error/userinvalid"))
	} else if chatStatus.IsSudo(userID) {
		go r.TelegramProvider.ReplyText(trans.GetString(msg.Chat.Id, "error/userissudo"))
	} else if reason == "" {
		reason = "None Specified"
	}

	res, err := securityService.FindGlobalBan(userID)
	if res != nil && err == nil {
		res.Reason = reason
		securityService.UpdateGlobalBan(*res)
		go r.TelegramProvider.ReplyText(trans.GetString(msg.Chat.Id, "actions/updgban"))
	} else {
		banStruct := model.GlobalBan{
			UserID:     userID,
			Reason:     reason,
			BannedBy:   msg.From.Id,
			BannedFrom: msg.Chat.Id,
			TimeAdded:  int(time.Now().Unix()),
		}
		securityService.UpdateGlobalBan(banStruct)
		go r.TelegramProvider.ReplyText(trans.GetString(msg.Chat.Id, "actions/newgban"))
	}
	return nil
}

func (r CommandAdmin) UnBanUser(_ ext.Bot, u *gotgbot.Update, args []string) error {
	r.TelegramProvider.Init(u)
	msg := u.EffectiveMessage

	if !chatStatus.RequireSudo(msg.From.Id, r.TelegramProvider) {
		return nil
	}

	target, _ := strconv.Atoi(args[0])
	if chatStatus.IsSudo(target) {
		go r.TelegramProvider.ReplyText(trans.GetString(msg.Chat.Id, "error/userissudo"))
	}

	res, err := securityService.FindGlobalBan(target)
	if res != nil && err == nil {
		securityService.RemoveGlobalBan(target)
		go r.TelegramProvider.ReplyText(trans.GetString(msg.Chat.Id, "actions/ungban"))
	} else {
		go r.TelegramProvider.ReplyText(trans.GetString(msg.Chat.Id, "error/userinvalid"))
	}
	return nil
}

func (r CommandAdmin) Debug(b ext.Bot, u *gotgbot.Update) error {
	r.TelegramProvider.Init(u)
	msg := u.EffectiveMessage

	// Permission Check
	if !chatStatus.RequireSudo(msg.From.Id, r.TelegramProvider) {
		return nil
	}

	if msg.ReplyToMessage != nil {
		//jsonData, err := json.Marshal(msg.ReplyToMessage)
		//err_handler.HandleErr(err)
		output, err := json.MarshalIndent(msg.ReplyToMessage, "", "  ")
		errHandler.Error(err)
		go r.TelegramProvider.SendText(
			string(output),
			msg.Chat.Id,
			0,
			nil,
		)
		return nil
	}

	//jsonData, err := json.Marshal(msg)
	//err_handler.HandleErr(err)
	output, err := json.MarshalIndent(msg, "", "  ")
	errHandler.Error(err)
	go r.TelegramProvider.SendText(
		string(output),
		msg.Chat.Id,
		0,
		nil,
	)
	return nil
}
