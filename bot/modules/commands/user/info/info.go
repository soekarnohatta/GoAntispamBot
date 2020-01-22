package info

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/ext/helpers"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/shirou/gopsutil/host"
	"github.com/sirupsen/logrus"
	"math"
	"strconv"
	"time"

	"github.com/jumatberkah/antispambot/bot/modules/helpers/extraction"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
)

var btnList = function.BuildKeyboard("data/keyboard/info.json", 2)

func getUser(b ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat

	userId := extraction.ExtractUser(msg, args)
	if userId != 0 {
		replyText := "👤*User Info*\n"
		userInfo := sql.GetUser(userId)
		if userInfo != nil {
			val := map[string]string{
				"1": strconv.Itoa(userId),
				"2": userInfo.FirstName,
				"3": userInfo.LastName,
				"4": userInfo.UserName,
			}
			replyText += function.GetStringf(chat.Id, "modules/info/info.go:29", val)
		}

		spamStatus := sql.GetUserSpam(userId)
		if spamStatus != nil {
			timeBanned, _ := strconv.ParseInt(fmt.Sprint(spamStatus.TimeAdded), 10, 64)
			val := map[string]string{"1": spamStatus.Reason, "2": spamStatus.Banner,
				"3": fmt.Sprint(time.Unix(timeBanned, 0))}
			replyText += function.GetStringf(chat.Id, "modules/info/info.go:35", val)
		}

		if replyText != "👤*User Info*\n" {
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

	info, _ := host.Info()
	replyTxt := fmt.Sprintf("*Bot Info*\n"+
		"👤*Bot Name :* %v\n"+
		"🤖*Bot Username :* @%v\n"+
		"🖥*Host OS :* %v\n"+
		"⚙*Host Name :* %v\n"+
		"⏱*Host Uptime :* `%v`\n"+
		"💽*Kernel Version :* %v\n"+
		"💾*Platform :* %v\n",
		helpers.EscapeMarkdown(b.FirstName),
		helpers.EscapeMarkdown(b.UserName),
		helpers.EscapeMarkdown(info.OS),
		helpers.EscapeMarkdown(info.Hostname),
		_convertSeconds(info.Uptime),
		info.KernelVersion,
		info.Platform,
	)

	replyMsg := b.NewSendableMessage(chat.Id, replyTxt)
	replyMsg.ParseMode = "Markdown"
	replyMsg.ReplyMarkup = &ext.InlineKeyboardMarkup{InlineKeyboard: &btnList}
	replyMsg.ReplyToMessageId = msg.MessageId
	_, err := replyMsg.Send()
	return err
}

func _convertSeconds(input uint64) (result string) {
	if input != 0 {
		years := math.Floor(float64(input) / 60 / 60 / 24 / 7 / 30 / 12)
		seconds := input % (60 * 60 * 24 * 7 * 30 * 12)
		months := math.Floor(float64(seconds) / 60 / 60 / 24 / 7 / 30)
		seconds = input % (60 * 60 * 24 * 7 * 30)
		weeks := math.Floor(float64(seconds) / 60 / 60 / 24 / 7)
		seconds = input % (60 * 60 * 24 * 7)
		days := math.Floor(float64(seconds) / 60 / 60 / 24)
		seconds = input % (60 * 60 * 24)
		hours := math.Floor(float64(seconds) / 60 / 60)
		seconds = input % (60 * 60)
		minutes := math.Floor(float64(seconds) / 60)
		seconds = input % 60

		if years > 0 {
			result = _plural(int(years), "year") + _plural(int(months), "month") + _plural(int(weeks), "week") + _plural(int(days), "day") + _plural(int(hours), "hour") + _plural(int(minutes), "minute") + _plural(int(seconds), "second")
		} else if months > 0 {
			result = _plural(int(months), "month") + _plural(int(weeks), "week") + _plural(int(days), "day") + _plural(int(hours), "hour") + _plural(int(minutes), "minute") + _plural(int(seconds), "second")
		} else if weeks > 0 {
			result = _plural(int(weeks), "week") + _plural(int(days), "day") + _plural(int(hours), "hour") + _plural(int(minutes), "minute") + _plural(int(seconds), "second")
		} else if days > 0 {
			result = _plural(int(days), "day") + _plural(int(hours), "hour") + _plural(int(minutes), "minute") + _plural(int(seconds), "second")
		} else if hours > 0 {
			result = _plural(int(hours), "hour") + _plural(int(minutes), "minute") + _plural(int(seconds), "second")
		} else if minutes > 0 {
			result = _plural(int(minutes), "minute") + _plural(int(seconds), "second")
		} else {
			result = _plural(int(seconds), "second")
		}
		return
	}
	return
}

func _plural(count int, singular string) (result string) {
	if (count == 1) || (count == 0) {
		result = strconv.Itoa(count) + " " + singular + " "
	} else {
		result = strconv.Itoa(count) + " " + singular + "s "
	}
	return
}

func LoadInfo(u *gotgbot.Updater) {
	defer logrus.Info("Info Module Loaded...")
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("whois", []rune{'/', '.'}, getUser))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("info", []rune{'/', '.'}, getBot))
}
