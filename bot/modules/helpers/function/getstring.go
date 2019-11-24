package function

import (
	"fmt"
	"github.com/PaulSonOfLars/goloc"
	"github.com/go-redis/redis"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/caching"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
	"io/ioutil"
)

func LoadAllLang() {
	files, err := ioutil.ReadDir("trans")
	err_handler.FatalError(err)
	for _, f := range files {
		if f.IsDir() {
			goloc.LoadAll(f.Name())
		}
	}
}

func GetString(chatID int, val string) string {
	var err error
	lang, err := caching.REDIS.Get(fmt.Sprintf("lang_%v", chatID)).Result()
	if err == redis.Nil || lang == "" {
		lang = sql.GetLang(chatID).Lang
	} else if err != nil {
		lang = "en"
	}
	return goloc.Trnl(lang, val)
}

func GetStringf(chatID int, val string, args map[string]string) string {
	var err error
	lang, err := caching.REDIS.Get(fmt.Sprintf("lang_%v", chatID)).Result()
	if err == redis.Nil || lang == "" {
		lang = sql.GetLang(chatID).Lang
	} else if err != nil {
		lang = "en"
	}
	return goloc.Trnlf(lang, val, args)
}
