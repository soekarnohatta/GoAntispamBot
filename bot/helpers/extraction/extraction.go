// Package extraction is a package that handles something to extract.
// This package should handle all extractions activity.
package extraction

import (
	"github.com/PaulSonOfLars/gotgbot/ext"
	"strconv"
	"strings"
	"unicode"

	"GoAntispamBot/bot/helpers/trans"
	"GoAntispamBot/bot/providers/telegramProvider"
	"GoAntispamBot/bot/services/logService"
)

// idFromReply will return the user id and the message of the user
// from the replied message.
func idFromReply(m *ext.Message) (int, string) {
	prevMessage := m.ReplyToMessage
	if prevMessage == nil {
		return 0, ""
	}

	userID := prevMessage.From.Id
	res := strings.SplitN(prevMessage.Text, " ", 2)
	if len(res) < 2 {
		return userID, ""
	}
	return userID, res[1]
}

// ExtractUserText will return the user id and the message of the user.
func ExtractUserText(telegramProvider telegramProvider.TelegramProvider, args []string) (int, string) {
	m := telegramProvider.Message
	prevMessage := m.ReplyToMessage
	splitText := strings.SplitN(m.Text, " ", 2)

	if len(splitText) < 2 {
		return idFromReply(m)
	}

	textToParse := splitText[1]
	text := ""

	var ent *ext.ParsedMessageEntity = nil
	accepted := make(map[string]struct{})
	accepted["text_mention"] = struct{}{}
	entities := m.ParseEntityTypes(accepted)
	if len(entities) > 0 {
		ent = &entities[0]
	}

	if entities != nil && ent != nil && ent.Offset == (len(m.Text)-len(textToParse)) {
		ent = &entities[0]
		userID := ent.User.Id
		text = m.Text[ent.Offset+ent.Length:]
		return userID, text
	} else if len(args) >= 1 && args[0][0] == '@' {
		user, _ := strconv.Atoi(args[0])
		userIDs, err := logService.FindUser(user)
		if userIDs == nil || err != nil {
			go telegramProvider.SendText(
				trans.GetString(m.Chat.Id, "error/usernotindb"),
				m.Chat.Id,
				0,
				nil,
			)
			return 0, ""
		} else {
			res := strings.SplitN(m.Text, " ", 2)
			if len(res) >= 3 {
				text = res[2]
				return userIDs.UserID, text
			}
		}
	} else if len(args) >= 1 {
		isID := true
		for _, r := range args[0] {
			if unicode.IsDigit(r) {
				continue
			} else {
				isID = false
				break
			}
		}

		if isID {
			userID, _ := strconv.Atoi(args[0])
			res := strings.SplitN(m.Text, " ", 2)
			if len(res) >= 3 {
				text := res[2]
				return userID, text
			}
		}

	} else if prevMessage != nil {
		_, _ = prevMessage.Delete()
		userID, text := idFromReply(prevMessage)
		return userID, text
	}

	return 0, ""
}
