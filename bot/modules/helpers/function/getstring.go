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
	lang, err := caching.REDIS.Get(fmt.Sprintf("lang_%v", chatID)).Result()
	if err != nil {
		if err == redis.Nil || lang == "" {
			lg := sql.GetLang(chatID)
			if lg != nil {
				lang = lg.Lang
			} else {
				lang = "en"
			}
		} else {
			lang = "en"
		}
	}

	ret := goloc.Trnl(lang, val)
	if ret != "" {
		return ret
	}
	return "None"
}

func GetStringf(chatID int, val string, args map[string]string) string {
	if args != nil && chatID != 0 && val != "" {
		lang, err := caching.REDIS.Get(fmt.Sprintf("lang_%v", chatID)).Result()
		if err != nil {
			if err == redis.Nil || lang == "" {
				lg := sql.GetLang(chatID)
				if lg != nil {
					lang = lg.Lang
				} else {
					lang = "en"
				}
			} else {
				lang = "en"
			}
		}
		ret := goloc.Trnlf(lang, val, args)
		if ret != "" {
			return ret
		}
	}
	return "None"
}
