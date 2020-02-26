package function

import (
	"encoding/json"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"html"
	"io/ioutil"
	"net/http"
	"strconv"
	"unicode"

	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/helpers/logger"
	"github.com/jumatberkah/antispambot/bot/sql"
)

type casban struct {
	Status bool `json:"ok"`
}

func CasListener(b ext.Bot, u *gotgbot.Update) (bool, error) {
	user := u.EffectiveUser
	spam := &casban{}
	response, err := http.Get(fmt.Sprintf("https://combot.org/api/cas/check?user_id=%v", user.Id))
	err_handler.HandleErr(err)
	defer response.Body.Close()
	if response != nil {
		body, _ := ioutil.ReadAll(response.Body)
		_ = json.Unmarshal(body, &spam)
		if spam.Status == true {
			err = SpamFunc(b, u)
			err_handler.HandleErr(err)
			err = logger.SendLog(b, u, "spam", "CAS Banned (Powered By CAS)")
			return true, err
		}
		return false, gotgbot.ContinueGroups{}
	}
	return false, gotgbot.ContinueGroups{}
}

func CheckChinese(val string) bool {
	for _, rangeTxt := range val {
		if unicode.Is(unicode.Han, rangeTxt) {
			return true
		}
		return false
	}
	return false
}

func CheckArabs(val string) bool {
	for _, rangeTxt := range val {
		if unicode.Is(unicode.Arabic, rangeTxt) {
			return true
		}
		return false
	}
	return false
}

func SpamFunc(b ext.Bot, u *gotgbot.Update) error {
	user := u.EffectiveUser
	chat := u.EffectiveChat
	msg := u.EffectiveMessage
	db := sql.GetSetting(chat.Id)
	if db == nil {
		return nil
	}

	val := map[string]string{
		"1": strconv.Itoa(user.Id),
		"2": html.EscapeString(user.FirstName),
		"3": strconv.Itoa(user.Id),
	}
	txtBan := GetStringf(chat.Id, "handlers/listener/listener.go:580", val)

	restrictSend := b.NewSendableKickChatMember(chat.Id, user.Id)
	restrictSend.UntilDate = -1
	_, err := restrictSend.Send()
	if err != nil {
		if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
			txtBan = GetStringf(chat.Id, "handlers/listener/listener.go:warnspam", val)
			_, err := msg.ReplyHTML(txtBan)
			if err != nil {
				if err.Error() == "Bad Request: reply message not found" {
					_, err = b.SendMessageHTML(chat.Id, txtBan)
					return err
				}
			}
			return err
		}
	}

	_, err = msg.ReplyHTML(txtBan)
	if err != nil {
		if err.Error() == "Bad Request: reply message not found" {
			_, err = b.SendMessageHTML(chat.Id, txtBan)
			return err
		}
	}

	if db.Deletion == "true" {
		_, err = b.DeleteMessage(chat.Id, msg.MessageId)
		if err != nil {
			if err.Error() == "Bad Request: message can't be deleted" {
				err_handler.HandleTgErr(b, u, err)
				return err
			}
		}
	}
	return nil
}
