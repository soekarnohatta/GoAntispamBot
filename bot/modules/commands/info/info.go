package info

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/extraction"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
	"github.com/shirou/gopsutil/host"
	"github.com/sirupsen/logrus"
	"strconv"
)

func getUser(b ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat

	var replyText = "👤*User Info*\n"

	userId := extraction.ExtractUser(msg, args)
	if userId != 0 {
		userInfo := sql.GetUser(userId)
		if userInfo != nil {
			val := map[string]string{"1": strconv.Itoa(userId), "2": userInfo.FirstName, "3": userInfo.LastName, "4": userInfo.UserName}
			replyText += function.GetStringf(chat.Id, "modules/info/info.go:29", val)
		}

		spamStatus := sql.GetUserSpam(userId)
		if spamStatus != nil {
			val := map[string]string{"1": spamStatus.Reason}
			replyText += function.GetStringf(chat.Id, "modules/info/info.go:35", val)
		}

		if replyText != "" {
			message := b.NewSendableMessage(chat.Id, replyText)
			message.ParseMode = "markdown"
			message.ReplyToMessageId = msg.MessageId
			_, err := message.Send()
			if err != nil {
				if err.Error() == "Bad Request: reply message not found" {
					message.ReplyToMessageId = 0
					_, err := message.Send()
					return err
				}
			}
			return err
		}
	} else {
		_, err := msg.ReplyText(function.GetString(chat.Id, "modules/info/info.go:51"))
		return err
	}
	return nil
}

func getBot(b ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat

	infoButtons := [][]ext.InlineKeyboardButton{make([]ext.InlineKeyboardButton, 2), make([]ext.InlineKeyboardButton, 2)}
	infoButtons[0][0] = ext.InlineKeyboardButton{
		Text:         "📡 About",
		CallbackData: fmt.Sprintf("info(%v)", "about"),
	}
	infoButtons[0][1] = ext.InlineKeyboardButton{
		Text: "📀 Source",
		Url:  "https://github.com/soekarnohatta/antispambot",
	}
	infoButtons[1][0] = ext.InlineKeyboardButton{
		Text:         "💹 Donate",
		CallbackData: fmt.Sprintf("info(%v)", "donate"),
	}

	info, _ := host.Info()
	replyTxt := fmt.Sprintf("🤖*Bot Info*\n"+
		"👤Bot Name : %v\n"+
		"🤖Bot Username : @%v\n"+
		"🖥Host OS : %v\n"+
		"⚙Host Name : %v\n"+
		"⏱Host Uptime : %v\n"+
		"💽Kernel Version : %v\n"+
		"💾Platform : %v\n", b.FirstName, b.UserName, info.OS,
		info.Hostname, convertseconds(info.Uptime), info.KernelVersion, info.Platform)

	replyMsg := b.NewSendableMessage(chat.Id, replyTxt)
	replyMsg.ParseMode = "Markdown"
	replyMsg.ReplyMarkup = &ext.InlineKeyboardMarkup{&infoButtons}
	replyMsg.ReplyToMessageId = msg.MessageId
	_, err := replyMsg.Send()
	return err
}

func convertseconds(seconds uint64) string {
	if seconds != 0 {
		days := seconds / 86400
		hours := seconds / 3600
		minutes := (seconds / 60) - (hours * 60)
		return fmt.Sprintf("`%v Day(s), %v Hour(s), %v Minute(s)`", days, hours, minutes)
	}
	return ""
}

func LoadInfo(u *gotgbot.Updater) {
	defer logrus.Info("Info Module Loaded...")
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("whois", []rune{'/', '.'}, getUser))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("info", []rune{'/', '.'}, getBot))
}
