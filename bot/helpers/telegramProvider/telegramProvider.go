package telegramProvider

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"time"

	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/helpers/function"
)

type RequestProvider struct {
	Update            *gotgbot.Update
	appendText        string
	Bot               ext.Bot
	Message           *ext.Message
	SentMessageId     int
	EditedMessageId   int
	CallbackMessageId int
	timeInit          float64
	timeProc          float64
}

func (r *RequestProvider) Init(u *gotgbot.Update) {
	var msg = u.EffectiveMessage
	r.Update = u
	r.Bot = msg.Bot

	if u.CallbackQuery != nil {
		r.Message = u.CallbackQuery.Message
		r.timeInit = function.GetJeda(time.Unix(int64(u.CallbackQuery.Message.Date), 0))
	} else if msg != nil {
		r.Message = msg
		r.timeInit = function.GetJeda(time.Unix(int64(r.Message.Date), 0))
	}
}

func (r *RequestProvider) ReplyHTML(text string) {
	r.timeProc = function.GetJeda(time.Unix(int64(r.Message.Date), 0))
	if text != "" {
		text += fmt.Sprintf("\n\n⏱ <code>%.3f s|⌛️%.3f s</code>", r.timeInit, r.timeProc)
	}

	_, err := r.Message.ReplyHTML(text)
	if err != nil {
		if err.Error() == "Bad Request: reply message not found" {
			_, err := r.Bot.SendMessageHTML(r.Message.Chat.Id, text)
			err_handler.HandleErr(err)
		}
	}
}

func (r *RequestProvider) ReplyMarkdown(text string) {
	r.timeProc = function.GetJeda(time.Unix(int64(r.Message.Date), 0))
	if text != "" {
		text += fmt.Sprintf("\n\n⏱ `%.3f s|⌛️%.3f s`", r.timeInit, r.timeProc)
	}

	_, err := r.Message.ReplyMarkdown(text)
	if err != nil {
		if err.Error() == "Bad Request: reply message not found" {
			_, err := r.Bot.SendMessageMarkdown(r.Message.Chat.Id, text)
			err_handler.HandleErr(err)
		}
	}
}

func (r *RequestProvider) EditMessageMarkdown(msgId int, text string) {
	r.timeProc = function.GetJeda(time.Unix(int64(r.Message.Date), 0))
	if text != "" {
		text += fmt.Sprintf("\n\n⏱ `%.3f s|⌛️%.3f s`", r.timeInit, r.timeProc)
	}

	_, err := r.Bot.EditMessageMarkdown(r.Message.Chat.Id, msgId, text)
	err_handler.HandleErr(err)
}

func (r *RequestProvider) EditMessageHTML(msgId int, text string) {
	r.timeProc = function.GetJeda(time.Unix(int64(r.Message.Date), 0))
	if text != "" {
		text += fmt.Sprintf("\n\n⏱ <code>%.3f s|⌛️%.3f s</code>", r.timeInit, r.timeProc)
	}

	_, err := r.Bot.EditMessageHTML(r.Message.Chat.Id, msgId, text)
	err_handler.HandleErr(err)
}

func (r *RequestProvider) SendTextAsync(txtToSend string, repMsgId int, customChatId int, parseMode string, button *ext.InlineKeyboardMarkup) {
	r.timeProc = function.GetJeda(time.Unix(int64(r.Message.Date), 0))
	if txtToSend != "" {
		switch parseMode {
		case parsemode.Html:
			txtToSend += fmt.Sprintf("\n\n⏱ <code>%.3f s|⌛️%.3f s</code>", r.timeInit, r.timeProc)
		case parsemode.Markdown:
			txtToSend += fmt.Sprintf("\n\n⏱ `%.3f s|⌛️%.3f s`", r.timeInit, r.timeProc)
		}
	}

	var chatId = r.Message.Chat.Id
	if customChatId != 0 {
		chatId = customChatId
	}

	var send = r.Bot.NewSendableMessage(chatId, txtToSend)
	send.ParseMode = parseMode
	send.ReplyMarkup = button
	send.ReplyToMessageId = repMsgId
	var _, err = send.Send()
	if err != nil {
		if err.Error() == "Bad Request: reply message not found" {
			send.ReplyToMessageId = 0
			var _, err = send.Send()
			err_handler.HandleErr(err)
		}
	}
}
