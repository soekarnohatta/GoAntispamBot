package telegramProvider

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/parsemode"

	"GoAntispamBot/bot/helpers/errHandler"
	"GoAntispamBot/bot/providers/errorProvider"
)

type TelegramProvider struct {
	Update          *gotgbot.Update
	Bot             ext.Bot
	Message         *ext.Message
	EditedMessage   *ext.Message
	MessageOrEdited *ext.Message
	SentMessageID   int

	// Self Added Type
	TimeInit float64
	TimeProc float64
}

func (r TelegramProvider) Init(u *gotgbot.Update) {
	r.Update = u
	r.Bot = u.Message.Bot

	if u.Message != nil {
		r.Message = u.Message
		r.MessageOrEdited = u.Message
	}

	if u.CallbackQuery != nil {
		r.Message = u.CallbackQuery.Message
	}

	if u.EditedMessage != nil {
		r.EditedMessage = u.EditedMessage
		r.MessageOrEdited = u.EditedMessage
	}
}

func (r TelegramProvider) ReplyText(txtToSend string) {
	chatID := r.Message.Chat.Id

	sendMsg := r.Bot.NewSendableMessage(chatID, txtToSend)
	sendMsg.ReplyToMessageId = r.Message.MessageId
	sendMsg.DisableWebPreview = true
	sendMsg.ParseMode = parsemode.Html
	sendMsg.DisableNotification = false
	send, err := sendMsg.Send()
	if err != nil {
		if err.Error() == errorProvider.RepMsgNotFound {
			sendMsg.ReplyToMessageId = 0
			_, _ = sendMsg.Send()
		} else {
			return
		}
	}

	if send != nil {
		r.SentMessageID = send.MessageId
	}
}

func (r TelegramProvider) SendText(txtToSend string, chatID int, repMsgID int, btn *ext.InlineKeyboardMarkup) {
	if chatID == 0 {
		chatID = r.Message.Chat.Id
	}
	sendMsg := r.Bot.NewSendableMessage(chatID, txtToSend)
	sendMsg.ReplyToMessageId = repMsgID
	sendMsg.DisableWebPreview = true
	sendMsg.ReplyMarkup = btn
	sendMsg.ParseMode = parsemode.Html
	sendMsg.DisableNotification = false
	send, err := sendMsg.Send()
	if err != nil {
		if err.Error() == errorProvider.RepMsgNotFound {
			sendMsg.ReplyToMessageId = 0
			_, _ = sendMsg.Send()
		} else {
			return
		}
	}

	if send != nil {
		r.SentMessageID = send.MessageId
	}
}

func (r TelegramProvider) EditMessage(txtEdit string, btn *ext.InlineKeyboardMarkup) {
	editMsg := r.Bot.NewSendableEditMessageText(r.Message.Chat.Id, r.SentMessageID, txtEdit)
	editMsg.DisableWebPreview = true
	editMsg.ReplyMarkup = btn
	editMsg.ParseMode = parsemode.Html
	edited, err := editMsg.Send()
	errHandler.Error(err)

	if edited != nil {
		r.EditedMessage = edited
	}
}

func (r TelegramProvider) KickMember(userID int, chatID int, dur int64) {
	kickMem := r.Bot.NewSendableKickChatMember(chatID, userID)
	kickMem.UntilDate = dur
	_, err := kickMem.Send()
	errHandler.SendError(err)
}
