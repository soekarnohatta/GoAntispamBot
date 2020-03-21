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
	"GoAntispamBot/bot/services"
)

type CommandHandler model.Command

func (r CommandHandler) BanUser(_ ext.Bot, u *gotgbot.Update, args []string) error {
	r.TelegramProvider.Init(u)
	msg := u.EffectiveMessage

	if !chatStatus.RequireSudo(msg.From.Id, r.TelegramProvider) {
		return nil
	}

	userID, reason := extraction.ExtractUserText(r.TelegramProvider, args)
	if userID == 0 {
		go r.TelegramProvider.SendText(
			trans.GetString(msg.Chat.Id, "error/userinvalid"),
			msg.Chat.Id,
			0,
			nil,
		)
	} else if chatStatus.IsSudo(userID) {
		go r.TelegramProvider.SendText(
			trans.GetString(msg.Chat.Id, "error/userissudo"),
			msg.Chat.Id,
			0,
			nil,
		)
	} else if reason == "" {
		reason = "None Specified"
	}

	res, err := services.FindGlobalBan(userID)
	if res != nil && err == nil {
		res.Reason = reason
		services.UpdateGlobalBan(*res)
		go r.TelegramProvider.SendText(
			trans.GetString(msg.Chat.Id, "actions/updgban"),
			msg.Chat.Id,
			0,
			nil,
		)
	} else {
		banStruct := model.GlobalBan{
			UserID:     userID,
			Reason:     reason,
			BannedBy:   msg.From.Id,
			BannedFrom: msg.Chat.Id,
			TimeAdded:  int(time.Now().Unix()),
		}
		services.UpdateGlobalBan(banStruct)
		go r.TelegramProvider.SendText(
			trans.GetString(msg.Chat.Id, "actions/newgban"),
			msg.Chat.Id,
			0,
			nil,
		)
	}
	return nil
}

func (r CommandHandler) UnBanUser(_ ext.Bot, u *gotgbot.Update, args []string) error {
	r.TelegramProvider.Init(u)
	msg := u.EffectiveMessage
	userID := msg.From.Id

	if !chatStatus.RequireSudo(msg.From.Id, r.TelegramProvider) {
		return nil
	}

	target, _ := strconv.Atoi(args[0])
	if chatStatus.IsSudo(target) {
		go r.TelegramProvider.SendText(
			trans.GetString(msg.Chat.Id, "error/userissudo"),
			msg.Chat.Id,
			0,
			nil,
		)
	}

	res, err := services.FindGlobalBan(target)
	if res != nil && err == nil {
		services.RemoveGlobalBan(target)
		go r.TelegramProvider.SendText(
			trans.GetString(msg.Chat.Id, "actions/ungban"),
			msg.Chat.Id,
			0,
			nil,
		)
	} else {
		go r.TelegramProvider.SendText(
			trans.GetString(msg.Chat.Id, "error/userinvalid"),
			msg.Chat.Id,
			0,
			nil,
		)
	}
}

func (r CommandHandler) Debug(b ext.Bot, u *gotgbot.Update) error {
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
