package listener

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/PaulSonOfLars/goloc"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/sirupsen/logrus"

	"github.com/jumatberkah/antispambot/bot/helpers/caching"
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/helpers/function"
	"github.com/jumatberkah/antispambot/bot/helpers/telegramProvider"
	"github.com/jumatberkah/antispambot/bot/sql"
)

var Req = new(telegramProvider.RequestProvider)

func handleLang(b ext.Bot, u *gotgbot.Update) error {
	Req.Init(u)
	query := u.CallbackQuery
	msg := query.Message
	pattern, _ := regexp.Compile(`lang\((.+?)\)\((.+?)\)`)

	if pattern.MatchString(query.Data) {
		lang := pattern.FindStringSubmatch(query.Data)[1]
		chatID, _ := strconv.Atoi(pattern.FindStringSubmatch(query.Data)[2])

		if !goloc.IsLangSupported(lang) {
			ans := b.NewSendableAnswerCallbackQuery(query.Id)
			ans.Text = function.GetString(msg.Chat.Id, "handlers/language/language.go:58")
			ans.ShowAlert = true
			_, err := ans.Send()
			err_handler.HandleErr(err)
			return err
		}

		_, _ = caching.REDIS.Set(fmt.Sprintf("lang_%v", chatID), lang, 7200).Result()
		sql.UpdateLang(chatID, lang)

		Req.EditMessageHTML(
			msg.MessageId,
			function.GetStringf(
				chatID,
				"handlers/language/language.go:51",
				map[string]string{"1": lang}),
		)
	}
	return gotgbot.ContinueGroups{}
}

func LoadLangListener(u *gotgbot.Updater) {
	defer logrus.Info("Lang Listener Loaded...")
	u.Dispatcher.AddHandler(handlers.NewCallback("lang", handleLang))
}
