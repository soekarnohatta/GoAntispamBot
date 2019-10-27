package extraction

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/google/uuid"
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func GetUserId(username string) int {
	if len(username) <= 5 {
		return 0
	}
	if username[0] == '@' {
		username = username[1:]
	}
	users := sql.GetUserIdByName(username)
	if users == nil {
		return 0
	}

	return users.UserId
}

func IdFromReply(m *ext.Message) (int, string) {
	prevMessage := m.ReplyToMessage
	if prevMessage == nil {
		return 0, ""
	}
	userId := prevMessage.From.Id
	res := strings.SplitN(m.Text, " ", 2)
	if len(res) < 2 {
		return userId, ""
	}
	return userId, res[1]
}

func ExtractUserAndText(m *ext.Message, args []string) (int, string) {
	prevMessage := m.ReplyToMessage
	splitText := strings.SplitN(m.Text, " ", 2)

	if len(splitText) < 2 {
		return IdFromReply(m)
	}

	textToParse := splitText[1]

	text := ""

	var userId int
	accepted := make(map[string]struct{})
	accepted["text_mention"] = struct{}{}

	entities := m.ParseEntityTypes(accepted)

	var ent *ext.ParsedMessageEntity
	var isId = false
	if len(entities) > 0 {
		ent = &entities[0]
	} else {
		ent = nil
	}

	if entities != nil && ent != nil && ent.Offset == (len(m.Text)-len(textToParse)) {
		ent = &entities[0]
		userId = ent.User.Id
		text = strconv.Itoa(int(m.Text[ent.Offset+ent.Length]))
	} else if len(args) >= 1 && args[0][0] == '@' {
		user := args[0]
		userId = GetUserId(user)
		if userId == 0 {
			_, err := m.ReplyText("Saya belum memiliki data dia di database saya.")
			err_handler.HandleErr(err)
			return 0, ""
		} else {
			res := strings.SplitN(m.Text, " ", 3)
			if len(res) >= 3 {
				text = res[2]
			}
		}
	} else if len(args) >= 1 {
		isId = true
		for _, arg := range args[0] {
			if unicode.IsDigit(arg) {
				continue
			} else {
				isId = false
				break
			}
		}
		if isId {
			userId, _ = strconv.Atoi(args[0])
			res := strings.SplitN(m.Text, " ", 3)
			if len(res) >= 3 {
				text = res[2]
			}
		}
	}
	if !isId && prevMessage != nil {
		_, parseErr := uuid.Parse(args[0])
		userId, text = IdFromReply(m)
		if parseErr == nil {
			return userId, text
		}
	} else if !isId {
		_, parseErr := uuid.Parse(args[0])
		if parseErr == nil {
			return userId, text
		}
	}

	_, err := m.Bot.GetChat(userId)
	if err != nil {
		_, err := m.ReplyText("Saya perlu melihat dia terlebih dahulu.")
		err_handler.HandleErr(err)
		return userId, text
	}
	return userId, text
}

func ExtractUser(message *ext.Message, args []string) int {
	userId, _ := ExtractUserAndText(message, args)
	return userId
}

func ExtractTime(b ext.Bot, m *ext.Message, timeVal string) int64 {
	lastLetter := timeVal[len(timeVal)-1:]
	lastLetter = strings.ToLower(lastLetter)
	var banTime int64
	if strings.ContainsAny(lastLetter, "m & d & h") {
		t := timeVal[:len(timeVal)-1]
		timeNum, err := strconv.Atoi(t)
		if err != nil {
			_, err := b.SendMessage(m.Chat.Id, "Invalid time amount specified.")
			err_handler.HandleErr(err)
			return -1
		}

		if lastLetter == "m" {
			banTime = time.Now().Unix() + int64(timeNum*60)
		} else if lastLetter == "h" {
			banTime = time.Now().Unix() + int64(timeNum*60*60)
		} else if lastLetter == "d" {
			banTime = time.Now().Unix() + int64(timeNum*24*60*60)
		} else {
			return 0
		}
		return banTime
	} else {
		_, err := b.SendMessage(m.Chat.Id,
			fmt.Sprintf("Invalid time type specified. Expected m, h, or d got: %s", lastLetter))
		err_handler.HandleErr(err)
		return -1
	}
}

func GetEmoji(chatId int) [][]string {
	chat := sql.GetUsername(chatId)
	pic := sql.GetPicture(chatId)
	ver := sql.GetVerify(chatId)
	tim := sql.GetSetting(chatId)
	spm := sql.GetEnforceGban(chatId)
	lastLetter := tim.Time[len(tim.Time)-1:]
	lastLetter = strings.ToLower(lastLetter)
	lst := make([][]string, 0)

	if chat.Option == "true" {
		chat.Option = "🔵"
	} else {
		chat.Option = "⚪"
	}
	if pic.Option == "true" {
		pic.Option = "🔵"
	} else {
		pic.Option = "⚪"
	}
	if ver.Option == "true" {
		ver.Option = "🔵"
	} else {
		ver.Option = "⚪"
	}
	if spm.Option == "true" {
		spm.Option = "🔵"
	} else {
		spm.Option = "⚪"
	}

	if chat.Deletion == "true" {
		chat.Deletion = "+ 🗑"
	} else {
		chat.Deletion = "-"
	}
	if pic.Deletion == "true" {
		pic.Deletion = "+ 🗑"
	} else {
		pic.Deletion = "-"
	}
	if tim.Deletion == "true" {
		tim.Deletion = "+ 🗑"
	} else {
		tim.Deletion = "-"
	}
	if ver.Deletion == "true" {
		ver.Deletion = "+ 🗑"
	} else {
		ver.Deletion = "-"
	}

	if chat.Action == "mute" {
		chat.Action = "+ 🔇"
	} else if chat.Action == "ban" {
		chat.Action = "+ ⛔"
	} else if chat.Action == "kick" {
		chat.Action = "+ 🚷"
	} else if chat.Action == "warn" {
		chat.Action = "+ Warn"
	} else {
		chat.Action = "+ None"
	}

	if pic.Action == "mute" {
		pic.Action = "+ 🔇"
	} else if pic.Action == "ban" {
		pic.Action = "+ ⛔"
	} else if pic.Action == "kick" {
		pic.Action = "+ 🚷"
	} else if pic.Action == "warn" {
		pic.Action = "+ Warn"
	} else {
		pic.Action = "+ None"
	}

	opt := make([]string, 4)
	opt[0] = chat.Option
	opt[1] = pic.Option
	opt[2] = ver.Option
	opt[3] = spm.Option

	act := make([]string, 2)
	act[0] = chat.Action
	act[1] = pic.Action

	del := make([]string, 4)
	del[0] = chat.Deletion
	del[1] = pic.Deletion
	del[2] = tim.Deletion
	del[3] = ver.Deletion

	ti := make([]string, 1)
	ti[0] = tim.Time

	gu := make([]string, 1)
	gu[0] = lastLetter

	lst = append(lst, opt, act, del, ti, gu)
	return lst
}
