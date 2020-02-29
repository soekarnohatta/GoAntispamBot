package listener

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
)

func handleReport(b ext.Bot, u *gotgbot.Update) error {
	query := u.CallbackQuery
	pattern, _ := regexp.Compile(`report\((.+?)\)\((.+?)\)\((.+?)\)`)

	if pattern.MatchString(query.Data) {
		action := pattern.FindStringSubmatch(query.Data)[1]
		chatID, _ := strconv.Atoi(pattern.FindStringSubmatch(query.Data)[2])
		userID, _ := strconv.Atoi(pattern.FindStringSubmatch(query.Data)[3])

		switch action {
		case "kick":
			_, _ = b.KickChatMember(chatID, userID)
			_, _ = b.AnswerCallbackQueryText(query.Id, "Kicked.", true)
			_, _ = query.Message.Delete()
		case "ban":
			ban := b.NewSendableKickChatMember(chatID, userID)
			ban.UntilDate = -1
			_, _ = ban.Send()
			_, _ = b.AnswerCallbackQueryText(query.Id, "Banned.", true)
			_, _ = query.Message.Delete()
		case "del":
			_, _ = b.DeleteMessage(chatID, userID)
			_, _ = b.AnswerCallbackQueryText(query.Id, "Deleted.", true)
			_, _ = query.Message.Delete()
		}
	}
	return gotgbot.ContinueGroups{}
}

func LoadReportListener(u *gotgbot.Updater) {
	defer logrus.Info("Help Listener Loaded...")
	u.Dispatcher.AddHandler(handlers.NewCallback("report", handleReport))
}
