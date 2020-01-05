package extraction

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/google/uuid"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/err_handler"
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

	var userId = 0
	accepted := make(map[string]struct{})
	accepted["text_mention"] = struct{}{}

	entities := m.ParseEntityTypes(accepted)

	var ent *ext.ParsedMessageEntity = nil
	var isId = false

	if len(entities) > 0 {
		ent = &entities[0]
	} else {
		ent = nil
	}

	if entities != nil && ent != nil && ent.Offset == (len(m.Text)-len(textToParse)) {
		ent = &entities[0]
		userId = ent.User.Id
		res := strings.SplitN(m.Text, ent.Text, 2)
		text = res[1]
	} else if len(args) >= 1 && args[0][0] == '@' {
		user := args[0]
		userId = GetUserId(user)
		if userId == 0 {
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
		err_handler.HandleErr(parseErr)
		userId, text = IdFromReply(m)
		return userId, text
	} else if !isId {
		_, parseErr := uuid.Parse(args[0])
		if parseErr == nil {
			return userId, text
		}
	}

	_, err := m.Bot.GetChat(userId)
	if err != nil {
		return 0, text
	}
	return userId, text
}

func ExtractUser(message *ext.Message, args []string) int {
	ret, _ := ExtractUserAndText(message, args)
	return ret
}

func ExtractTime(b ext.Bot, m *ext.Message, timeVal string) int64 {
	lastLetter := timeVal[len(timeVal)-1:]
	lastLetter = strings.ToLower(lastLetter)
	var ret int64 = 0

	if strings.ContainsAny(lastLetter, "m & d & h") {
		t := timeVal[:len(timeVal)-1]
		timeNum, err := strconv.Atoi(t)

		if err != nil {
			_, err := b.SendMessage(m.Chat.Id, "Invalid time amount specified.")
			err_handler.HandleErr(err)
			return -1
		}

		if lastLetter == "m" {
			ret = time.Now().Unix() + int64(timeNum*60)
		} else if lastLetter == "h" {
			ret = time.Now().Unix() + int64(timeNum*60*60)
		} else if lastLetter == "d" {
			ret = time.Now().Unix() + int64(timeNum*24*60*60)
		} else {
			return -1
		}

		return ret
	} else {
		_, err := b.SendMessage(m.Chat.Id,
			fmt.Sprintf("Invalid time type specified. Expected m, h, or d got: %s", lastLetter))
		err_handler.HandleErr(err)
		return -1
	}
}
