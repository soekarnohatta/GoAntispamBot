package trans

import (
	"github.com/PaulSonOfLars/goloc"
	"io/ioutil"

	"GoAntispamBot/bot/helpers/errHandler"
	"GoAntispamBot/bot/services/langService"
)

const notFound = "Error: The desired string is not found"

func init() {
	files, err := ioutil.ReadDir("./bot/data/trans")
	errHandler.Fatal(err)
	for _, data := range files {
		if data.IsDir() {
			goloc.LoadAll(data.Name())
		}
	}
}

func GetString(chatId int, key string) string {
	if chatId != 0 && key != "" {
		lang := langService.FindLang(chatId)
		ret := goloc.Trnl(lang, key)
		if ret != "" {
			return ret
		}
	}
	return notFound
}

func GetStringf(chatId int, key string, args map[string]string) string {
	if args != nil && chatId != 0 && key != "" {
		lang := langService.FindLang(chatId)
		ret := goloc.Trnlf(lang, key, args)
		if ret != "" {
			return ret
		}
	}
	return notFound
}
