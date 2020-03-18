package keyboard

import (
	"encoding/json"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"GoAntispamBot/bot/helpers/errHandler"
	"GoAntispamBot/bot/model"
)

// This function will handle buttons creation from a json file
// and defined size.
func BuildKeyboard(path string, size int) (res [][]ext.InlineKeyboardButton) {
	jsonFile, err := os.Open(path)
	errHandler.Error(err)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result []model.Button
	_ = json.Unmarshal(byteValue, &result)

	var btnList []ext.InlineKeyboardButton
	for _, data := range result {
		if isValidUrl(data.Data) {
			btnList = append(btnList, ext.InlineKeyboardButton{Text: data.Text, Url: data.Data})
		} else {
			btnList = append(btnList, ext.InlineKeyboardButton{Text: data.Text, CallbackData: data.Data})
		}
	}

	for size < len(btnList) {
		btnList, res = btnList[size:], append(res, btnList[0:size:size])
	}

	return append(res, btnList)
}

// This function will handle buttons creation from a json file
// and defined size but with extra args to be placed inside the button.
func BuildKeyboardf(path string, size int, dataMap map[string]string) (res [][]ext.InlineKeyboardButton) {
	jsonFile, err := os.Open(path)
	errHandler.Error(err)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result []model.Button
	_ = json.Unmarshal(byteValue, &result)

	var btnList []ext.InlineKeyboardButton
	for _, data := range result {
		var replData []string
		for k, v := range dataMap {
			replData = append(replData, "{"+k+"}", v)
		}
		repl := strings.NewReplacer(replData...)

		if isValidUrl(repl.Replace(data.Data)) {
			btnList = append(btnList, ext.InlineKeyboardButton{
				Text: repl.Replace(data.Text),
				Url:  repl.Replace(data.Data),
			})
		} else {
			btnList = append(btnList, ext.InlineKeyboardButton{
				Text:         repl.Replace(data.Text),
				CallbackData: repl.Replace(data.Data),
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

