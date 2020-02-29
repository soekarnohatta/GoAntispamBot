package reporting

import (
	"encoding/json"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/ext/helpers"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"

	"github.com/jumatberkah/antispambot/bot/helpers/caching"
	"github.com/jumatberkah/antispambot/bot/helpers/chat_status"
	"github.com/jumatberkah/antispambot/bot/helpers/function"
	"github.com/jumatberkah/antispambot/bot/helpers/telegramProvider"
)

type adminCache struct {
	Admin []string `json:"admin"`
}

var _requestProvider = new(telegramProvider.RequestProvider)

func report(b ext.Bot, u *gotgbot.Update) error {
	_requestProvider.Init(u)
	msg := u.EffectiveMessage
	chat := u.EffectiveChat

	if !chat_status.RequireSupergroup(chat, msg) {
		return nil
	}

	reason := "No reason has been specified"
	if len(msg.Text) > 7 {
		splitReason := strings.Split(msg.Text, "report")[1]
		if splitReason != "" {
			reason = splitReason
		}
	}

	if msg.ReplyToMessage != nil {
		replyMsg := b.NewSendableMessage(
			chat.Id,
			function.GetString(chat.Id, "handlers/reporting/report.go:start"),
		)
		replyMsg.ParseMode = parsemode.Markdown
		replyMsg.ReplyToMessageId = msg.ReplyToMessage.MessageId
		sent, err := replyMsg.Send()
		go reportUser(b, u, msg, reason, sent)
		return err
	}
	_requestProvider.ReplyMarkdown(
		function.GetString(chat.Id, "handlers/reporting/report.go:reportFailed"))
	return nil
}

func reportUser(b ext.Bot, u *gotgbot.Update, msg *ext.Message, reason string, sent *ext.Message) {
	_requestProvider.Init(u)
	if sent == nil || msg == nil {
		return
	}
	admins, err := caching.CACHE.Get(fmt.Sprintf("admin_%v", msg.Chat.Id))

	if err != nil {
		chat_status.AdminCache(msg.Chat)
		admins, _ = caching.CACHE.Get(fmt.Sprintf("admin_%v", msg.Chat.Id))
	}

	var x adminCache
	_ = json.Unmarshal(admins, &x)

	rep := msg.ReplyToMessage
	reportTxt := fmt.Sprintf("#REPORT\n"+
		"Reported User : [%v](tg://user?id=%v) \\[`%v`] \n"+
		"Chat : %v \\[`%v`]\n"+
		"Message Link : [Here](https://t.me/%v/%v)\n"+
		"Reporter : [%v](tg://user?id=%v) \\[`%v`] \n"+
		"Reason : `%v` \n"+
		"Time Reported : `%v` \n",
		helpers.EscapeMarkdown(rep.From.FirstName),
		rep.From.Id,
		rep.From.Id,
		helpers.EscapeMarkdown(msg.Chat.Title),
		rep.Chat.Id,
		helpers.EscapeMarkdown(msg.Chat.Username),
		rep.MessageId,
		helpers.EscapeMarkdown(msg.From.FirstName),
		msg.From.Id,
		msg.From.Id,
		helpers.EscapeMarkdown(reason),
		time.Now())

	reportButtons := function.BuildKeyboardf(
		"data/keyboard/reporting.json",
		2,
		map[string]string{
			"1": strconv.Itoa(rep.Chat.Id),
			"2": strconv.Itoa(rep.From.Id),
			"3": strconv.Itoa(rep.MessageId),
			"4": rep.Chat.Username,
		})

	counter := 0
	for _, adm := range x.Admin {
		userID, _ := strconv.Atoi(adm)
		sendMsg := b.NewSendableMessage(userID, reportTxt)
		sendMsg.ParseMode = parsemode.Markdown
		sendMsg.ReplyMarkup = &ext.InlineKeyboardMarkup{InlineKeyboard: &reportButtons}
		sendMsg.DisableWebPreview = true
		sendMsg.DisableNotification = true
		_, err := sendMsg.Send()
		if err == nil {
			counter++
		}

		_requestProvider.EditMessageMarkdown(
			sent.MessageId,
			function.GetStringf(
				sent.Chat.Id,
				"handlers/reporting/report.go:report",
				map[string]string{"1": fmt.Sprint(userID)}),
		)
	}

	_requestProvider.EditMessageMarkdown(
		sent.MessageId,
		fmt.Sprintf(
			"`Succesfully Reported to %d/%d admin(s)`",
			counter,
			len(x.Admin)),
	)
}

func LoadReport(u *gotgbot.Updater) {
	defer logrus.Info("Report Module Loaded...")
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("report", []rune{'/', '.'}, report))
}
