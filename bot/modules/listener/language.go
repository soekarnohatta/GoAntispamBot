package listener

import (
	"fmt"
	"github.com/PaulSonOfLars/goloc"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"

	"github.com/jumatberkah/antispambot/bot/modules/helpers/caching"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
)

func handleLang(b ext.Bot, u *gotgbot.Update) error {
	query := u.CallbackQuery
	msg := query.Message
	pattern, _ := regexp.Compile(`lang\((.+?)\)\((.+?)\)`)

	if pattern.MatchString(query.Data) {
		lang := pattern.FindStringSubmatch(query.Data)[1]
		chatId, _ := strconv.Atoi(pattern.FindStringSubmatch(query.Data)[2])

		if !goloc.IsLangSupported(lang) {
			ans := b.NewSendableAnswerCallbackQuery(query.Id)
			ans.Text = function.GetString(msg.Chat.Id, "modules/language/language.go:58")
			ans.ShowAlert = true
			_, err := ans.Send()
			err_handler.HandleErr(err)
			return err
		}

		_, _ = caching.REDIS.Set(fmt.Sprintf("lang_%v", chatId), lang, 7200).Result()
		sql.UpdateLang(chatId, lang)

		_, err := msg.EditHTML(function.GetStringf(
			chatId,
			"modules/language/language.go:51",
			map[string]string{"1": lang}),
		)
		err_handler.HandleErr(err)
	}
	return gotgbot.ContinueGroups{}
}

func LoadLangListener(u *gotgbot.Updater) {
	defer logrus.Info("Lang Listener Loaded...")
	u.Dispatcher.AddHandler(handlers.NewCallback("lang", handleLang))
}
