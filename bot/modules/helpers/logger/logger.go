package logger

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"github.com/jumatberkah/antispambot/bot"
	"html"
	"strconv"
	"time"
)

func SendBanLog(b ext.Bot, uid int, rson string, u *gotgbot.Update) error {
	user := u.EffectiveUser
	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	txtLog := fmt.Sprintf("#GBAN\n"+
		"<b>Sudo:</b> <a href=\"tg://user?id=%v\">%v</a>\n"+
		"<b>User ID:</b> <code>%v</code>\n"+
		"<b>Time:</b> <code>%v</code>\n"+
		"<b>Reason:</b> <code>%v</code>", user.Id, html.EscapeString(user.FirstName), strconv.Itoa(uid), formatted, rson)

	sendLog := b.NewSendableMessage(bot.BotConfig.LogBan, txtLog)
	sendLog.ParseMode = parsemode.Html
	_, err := sendLog.Send()
	return err
}

func SendLog(b ext.Bot, u *gotgbot.Update, param string, args string) error {
	user := u.EffectiveUser
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	waktu := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		waktu.Year(), waktu.Month(), waktu.Day(),
		waktu.Hour(), waktu.Minute(), waktu.Second())

	switch param {
	case "username":
		txtLog := fmt.Sprintf("#NOUSERNAME\n"+
			"<b>User ID:</b> <a href=\"tg://user?id=%v\">%v</a>[%v]\n"+
			"<b>Chat ID:</b> <code>%v</code>\n"+
			"<b>Chat Title:</b> <code>%v</code>\n"+
			"<b>Time:</b> <code>%v</code>\n"+
			"<b>Message:</b>\n%v", user.Id, html.EscapeString(user.FirstName), user.Id, chat.Id, chat.Title, formatted, msg.Text)

		sendLog := b.NewSendableMessage(bot.BotConfig.LogEvent, txtLog)
		sendLog.ParseMode = parsemode.Html
		_, err := sendLog.Send()
		return err
	case "picture":
		txtLog := fmt.Sprintf("#NOPROFILEPICTURE\n"+
			"<b>User ID:</b> <a href=\"tg://user?id=%v\">%v</a>[%v]\n"+
			"<b>Chat ID:</b> <code>%v</code>\n"+
			"<b>Chat Title:</b> <code>%v</code>\n"+
			"<b>Time:</b> <code>%v</code>\n"+
			"<b>Message:</b>\n%v", user.Id, html.EscapeString(user.FirstName), user.Id, chat.Id, chat.Title, formatted, msg.Text)

		sendLog := b.NewSendableMessage(bot.BotConfig.LogEvent, txtLog)
		sendLog.ParseMode = parsemode.Html
		_, err := sendLog.Send()
		return err
	case "welcome":
		txtLog := fmt.Sprintf("#NEWMEMBER\n"+
			"<b>User ID:</b> <a href=\"tg://user?id=%v\">%v</a>[%v]\n"+
			"<b>Chat ID:</b> <code>%v</code>\n"+
			"<b>Chat Title:</b> <code>%v</code>\n"+
			"<b>Time:</b> <code>%v</code>\n"+
			"<b>Event:</b>\n%v", user.Id, html.EscapeString(user.FirstName), user.Id, chat.Id, chat.Title, formatted, "NewChatMembers")

		sendLog := b.NewSendableMessage(bot.BotConfig.LogEvent, txtLog)
		sendLog.ParseMode = parsemode.Html
		_, err := sendLog.Send()
		return err
	case "error":
		txtLog := fmt.Sprintf("#ERROR\n"+
			"<b>Time:</b> <code>%v</code>\n"+
			"<b>Error Message:</b>\n%v", formatted, args)
		sendLog := b.NewSendableMessage(bot.BotConfig.LogEvent, txtLog)
		sendLog.ParseMode = parsemode.Html
		_, err := sendLog.Send()
		return err
	case "spam":
		txtLog := fmt.Sprintf("#SPAMMER\n"+
			"<b>User ID:</b> <a href=\"tg://user?id=%v\">%v</a>[%v]\n"+
			"<b>Chat ID:</b> <code>%v</code>\n"+
			"<b>Chat Title:</b> <code>%v</code>\n"+
			"<b>Time:</b> <code>%v</code>\n"+
			"<b>Reason:</b> <code>%v</code>\n"+
			"<b>Message:</b>\n%v", user.Id, html.EscapeString(user.FirstName), user.Id, chat.Id, chat.Title, formatted, args, msg.Text)

		sendLog := b.NewSendableMessage(bot.BotConfig.LogBan, txtLog)
		sendLog.ParseMode = parsemode.Html
		_, err := sendLog.Send()
		return err
	case "link":
		txtLog := fmt.Sprintf("#LINK/AR/CH\n"+
			"<b>User ID:</b> <a href=\"tg://user?id=%v\">%v</a>[%v]\n"+
			"<b>Chat ID:</b> <code>%v</code>\n"+
			"<b>Chat Title:</b> <code>%v</code>\n"+
			"<b>Time:</b> <code>%v</code>\n"+
			"<b>Message:</b>\n%v", user.Id, html.EscapeString(user.FirstName), user.Id, chat.Id, chat.Title, formatted, msg.Text)

		sendLog := b.NewSendableMessage(bot.BotConfig.LogBan, txtLog)
		sendLog.ParseMode = parsemode.Html
		_, err := sendLog.Send()
		return err
	}
	return nil
}
