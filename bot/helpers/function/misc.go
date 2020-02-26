package function

import (
	"encoding/json"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
)

func Contains(key []string, str string) bool {
	if key != nil && str != "" {
		for _, val := range key {
			if val == str {
				return true
			}
		}
		return false
	}
	return false
}

func BuildKeyboard(path string, size int) (res [][]ext.InlineKeyboardButton) {
	jsonFile, err := os.Open(path)
	err_handler.HandleErr(err)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result []map[string]interface{}
	_ = json.Unmarshal(byteValue, &result)

	var btnList []ext.InlineKeyboardButton
	for _, data := range result {
		btnText := fmt.Sprint(data["text"])
		btnData := fmt.Sprint(data["data"])
		if isValidUrl(btnData) {
			btnList = append(btnList, ext.InlineKeyboardButton{Text: btnText, Url: btnData})
		} else {
			btnList = append(btnList, ext.InlineKeyboardButton{Text: btnText, CallbackData: btnData})
		}
	}

	for size < len(btnList) {
		btnList, res = btnList[size:], append(res, btnList[0:size:size])
	}

	return append(res, btnList)
}

func BuildKeyboardf(path string, size int, dataMap map[string]string) (res [][]ext.InlineKeyboardButton) {
	jsonFile, err := os.Open(path)
	err_handler.HandleErr(err)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result []map[string]interface{}
	_ = json.Unmarshal(byteValue, &result)

	var btnList []ext.InlineKeyboardButton
	for _, data := range result {
		btnText := fmt.Sprint(data["text"])
		btnData := fmt.Sprint(data["data"])

		var replData []string
		for k, v := range dataMap {
			replData = append(replData, "{"+k+"}", v)
		}
		repl := strings.NewReplacer(replData...)

		if isValidUrl(repl.Replace(btnData)) {
			btnList = append(btnList, ext.InlineKeyboardButton{
				Text: repl.Replace(btnText),
				Url:  repl.Replace(btnData),
			})
		} else {
			btnList = append(btnList, ext.InlineKeyboardButton{
				Text:         repl.Replace(btnText),
				CallbackData: repl.Replace(btnData),
			})
		}
	}

	for size < len(btnList) {
		btnList, res = btnList[size:], append(res, btnList[0:size:size])
	}

	return append(res, btnList)
}

func isValidUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
