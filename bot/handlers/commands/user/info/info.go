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

	"github.com/jumatberkah/antispambot/bot/helpers/extraction"
	"github.com/jumatberkah/antispambot/bot/helpers/function"
	"github.com/jumatberkah/antispambot/bot/sql"
)

var btnList = function.BuildKeyboard("data/keyboard/info.json", 2)

func getUser(b ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat

	userID := extraction.ExtractUser(msg, args)
	if userID != 0 {
		replyText := "*User Info*"
		userInfo := sql.GetUser(userID)
		if userInfo != nil {
			val := map[string]string{
				"1": strconv.Itoa(userID),
				"2": helpers.EscapeMarkdown(userInfo.FirstName),
				"3": helpers.EscapeMarkdown(userInfo.LastName),
				"4": helpers.EscapeMarkdown(userInfo.UserName),
			}
			replyText += function.GetStringf(chat.Id, "handlers/info/info.go:29", val)
		}

		spamStatus := sql.GetUserSpam(userID)
		if spamStatus != nil {
			timeBanned, _ := strconv.ParseInt(
				fmt.Sprint(spamStatus.TimeAdded),
				10,
				64,
			)
			val := map[string]string{
				"1": helpers.EscapeMarkdown(spamStatus.Reason),
				"2": helpers.EscapeMarkdown(spamStatus.Banner),
				"3": fmt.Sprint(time.Unix(timeBanned, 0)),
			}

			replyText += function.GetStringf(chat.Id, "handlers/info/info.go:35", val)
		}

		checkCas := function.CheckCas(u)
		if checkCas {
			replyText += "\n*🏴󠁣󠁯󠁣󠁡󠁳󠁿 CAS Banned*: " + fmt.Sprint(checkCas)
		}

		if !(replyText == "*User Info*") {
			_, err := msg.ReplyMarkdown(replyText)
			return err
		}

		_, err := msg.ReplyHTML(function.GetString(chat.Id, "handlers/info/info.go:51"))
		return err
	}

	_, err := msg.ReplyHTML(function.GetString(chat.Id, "handlers/info/info.go:51"))
	return err
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
