package modules

import (
	"fmt"
	"github.com/PaulSonOfLars/goloc"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/go-redis/redis"
	"github.com/jumatberkah/antispambot/bot/helpers/chat_status"
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
	"io/ioutil"
)

func loadlang() {
	files, err := ioutil.ReadDir("trans")
	err_handler.FatalError(err)
	for _, f := range files {
		if f.IsDir() {
			goloc.LoadAll(f.Name())
		}
	}
}

// GetString -> Get the desired string
func GetString(chatID int, val string) string {
	var err error
	lang, err := sql.REDIS.Get(fmt.Sprintf("lang_%v", chatID)).Result()
	if err == redis.Nil || lang == "" {
		lang = sql.GetLang(chatID).Lang
	} else if err != nil {
		lang = "en"
	}
	return goloc.Trnl(lang, val)
}

// GetStringf -> Get the desired string
func GetStringf(chatID int, val string, args map[string]string) string {
	var err error
	lang, err := sql.REDIS.Get(fmt.Sprintf("lang_%v", chatID)).Result()
	if err == redis.Nil || lang == "" {
		lang = sql.GetLang(chatID).Lang
	} else if err != nil {
		lang = "id"
	}
	err_handler.HandleErr(err)
	return goloc.Trnlf(lang, val, args)
}

func setlang(b ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat_status.IsUserAdmin(chat, user.Id) {
		if len(args) != 0 {
			if goloc.IsLangSupported(args[0]) {
				_, err := sql.REDIS.Set(fmt.Sprintf("lang_%v", chat.Id), args[0], 0).Result()
				if err != nil {
					err = sql.UpdateLang(chat.Id, args[0])
					err_handler.HandleTgErr(b, u, err)
					_, err = msg.ReplyText(GetStringf(chat.Id, "modules/language.go:51",
						map[string]string{"1": args[0]}))
					return err
				}
				_, err = msg.ReplyText(GetStringf(chat.Id, "modules/language.go:51",
					map[string]string{"1": args[0]}))
				return err

			}
			_, err := msg.ReplyText(GetString(chat.Id, "modules/language.go:58"))
			return err
		}
		_, err := msg.ReplyText("Please insert the language code so that i can change your language")
		return err

	}
	_, err := msg.Delete()
	return err
}

// LoadLang -> Register handlers
func LoadLang(u *gotgbot.Updater) {
	loadlang()
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("setlang", setlang))
}
