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
	"github.com/jumatberkah/antispambot/bot/sql"
)

type CasBan struct {
	Status bool    `json:"ok"`
	Result Results `json:"result"`
}

type Results struct {
	Offenses  int      `json:"offenses"`
	Messages  []string `json:"messages"`
	TimeAdded int      `json:"time_added"`
}

func CasListener(b ext.Bot, u *gotgbot.Update) error {
	if CheckCas(u) {
		err := SpamFunc(b, u)
		err_handler.HandleErr(err)
		err = SendLog(b, u, "spam", "CAS Banned (Powered By CAS)")
		return gotgbot.EndGroups{}
	}
	return gotgbot.ContinueGroups{}
}

func CheckCas(u *gotgbot.Update) bool {
	user := u.EffectiveUser
	spam := new(CasBan)
	response, err := http.Get(fmt.Sprintf("https://api.cas.chat/check?user_id=%v", user.Id))
	err_handler.HandleErr(err)
	if response != nil {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		err = json.Unmarshal(body, &spam)
		err_handler.HandleErr(err)
		//z, _ := json.Marshal(&CasBan{})
		//fmt.Print(string(z))
		if spam.Status {
			return true
		}
		return false
	}
	return false
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

	a, err := msg.ReplyHTML(txtBan)
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

	_, _ = a.Delete()
	return nil
}
