package reporting

import (
	"encoding/json"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/caching"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/chat_status"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type adminList struct {
	Admin []string `json:"admin"`
}

func report(b ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat

	reason := "No reason has been specified"
	splitReason := strings.Split(msg.Text, "report")[1]
	if splitReason != "" {
		reason = splitReason
	}

	if msg.ReplyToMessage != nil {
		go reportUser(b, msg, reason)
		replyMsg := b.NewSendableMessage(chat.Id, function.GetString(chat.Id, "a"))
		replyMsg.ParseMode = parsemode.Markdown
		replyMsg.ReplyToMessageId = msg.ReplyToMessage.MessageId
		_, err := replyMsg.Send()
		return err
	}
	return nil
}

func reportUser(b ext.Bot, msg *ext.Message, reason string) {
	admins, err := caching.CACHE.Get(fmt.Sprintf("admin_%v", msg.Chat.Id))

	if err != nil {
		chat_status.AdminCache(msg.Chat)
		admins, _ = caching.CACHE.Get(fmt.Sprintf("admin_%v", msg.Chat.Id))
	}

	var x adminList
	_ = json.Unmarshal(admins, &x)

	rep := msg.ReplyToMessage
	reportTxt := fmt.Sprintf("#REPORT\n"+
		"Reported User : [%v](tg://user?id=%v) \\[`%v`] \n"+
		"Message Link : [Here](https://t.me/%v/%v)\n"+
		"Reporter : [%v](tg://user?id=%v) \\[`%v`] \n"+
		"Reason : `%v` \n"+
		"Time Reported : `%v` \n", rep.From.FirstName, rep.From.Id, rep.From.Id, msg.Chat.Username, rep.MessageId, msg.From.FirstName,
		msg.From.Id, msg.From.Id, reason, time.Now())

	reportButtons := [][]ext.InlineKeyboardButton{make([]ext.InlineKeyboardButton, 1), make([]ext.InlineKeyboardButton, 2),
		make([]ext.InlineKeyboardButton, 1)}
	reportButtons[0][0] = ext.InlineKeyboardButton{
		Text: "📝 Message Link",
		Url:  fmt.Sprintf("https://t.me/%v/%v", msg.Chat.Username, rep.MessageId),
	}
	reportButtons[1][0] = ext.InlineKeyboardButton{
		Text:         "🚷 Kick",
		CallbackData: fmt.Sprintf("report(kick)_%v", rep.From.Id),
	}
	reportButtons[1][1] = ext.InlineKeyboardButton{
		Text:         "🚫 Ban",
		CallbackData: fmt.Sprintf("report(ban)_%v", rep.From.Id),
	}
	reportButtons[2][0] = ext.InlineKeyboardButton{
		Text:         "❌ Delete Message",
		CallbackData: fmt.Sprintf("report(del)_%v", rep.MessageId),
	}

	for _, adm := range x.Admin {
		uId, _ := strconv.Atoi(adm)
		sendMsg := b.NewSendableMessage(uId, reportTxt)
		sendMsg.ParseMode = parsemode.Markdown
		sendMsg.ReplyMarkup = &ext.InlineKeyboardMarkup{&reportButtons}
		sendMsg.DisableWebPreview = true
		sendMsg.DisableNotification = true
		_, err := sendMsg.Send()
		err_handler.HandleErr(err)
	}
}

func LoadReport(u *gotgbot.Updater) {
	defer logrus.Info("Report Module Loaded...")
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("report", []rune{'/', '.'}, report))
}
