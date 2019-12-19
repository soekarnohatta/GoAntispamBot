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
	for _, data := range files {
		if data.IsDir() {
			go goloc.LoadAll(data.Name())
		}
	}
}

func GetString(chatId int, val string) string {
	if chatId != 0 && val != "" {
		lang := getLang(chatId)
		ret := goloc.Trnl(lang, val)
		if ret != "" {
			return ret
		}
	}
	return "None"
}

func GetStringf(chatId int, val string, args map[string]string) string {
	if args != nil && chatId != 0 && val != "" {
		lang := getLang(chatId)
		ret := goloc.Trnlf(lang, val, args)
		if ret != "" {
			return ret
		}
	}
	return "None"
}

func getLang(chatId int) string {
	lang, err := caching.REDIS.Get(fmt.Sprintf("lang_%v", chatId)).Result()
	if err != nil {
		if err == redis.Nil || lang == "" {
			lg := sql.GetLang(chatId)
			if lg != nil {
				lang = lg.Lang
			} else {
				lang = "en"
			}
		} else {
			lang = "en"
		}
	}
	return lang
}
