package function

import (
	"fmt"
	"github.com/PaulSonOfLars/goloc"
	"github.com/go-redis/redis"
	"io/ioutil"

	"github.com/jumatberkah/antispambot/bot/modules/helpers/caching"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
)

var dummy = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum eu viverra lacus. Aliquam a pellentesque libero. Ut semper ornare nulla eget suscipit. Aenean dictum scelerisque urna a sagittis. Morbi vel ex luctus, tristique eros sit amet, sagittis neque. Nulla quis odio massa. Donec vitae odio quis elit ultrices porttitor. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Praesent eu justo odio. Donec ac arcu lectus."

func LoadAllLang() {
	files, err := ioutil.ReadDir("./data/trans")
	err_handler.FatalError(err)
	for _, data := range files {
		if data.IsDir() {
			goloc.LoadAll(data.Name())
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
	return dummy
}

func GetStringf(chatId int, val string, args map[string]string) string {
	if args != nil && chatId != 0 && val != "" {
		lang := getLang(chatId)
		ret := goloc.Trnlf(lang, val, args)
		if ret != "" {
			return ret
		}
	}
	return dummy
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
