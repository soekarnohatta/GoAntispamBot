package extraction

import (
	"github.com/PaulSonOfLars/gotgbot/ext"
	"strconv"
	"strings"

	"GoAntispamBot/bot/helpers/trans"
	"GoAntispamBot/bot/services"
)

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

func extractUserText(m *ext.Message, args []string) (int, string) {
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

	if entities != nil && ent != nil && ent.Offset == (len(m.Text) - len(textToParse)) {
		ent = &entities[0]
		userID := ent.User.Id
		text = m.Text[ent.Offset+ent.Length:]
		return userID, text
	} else if len(args) >= 1 && args[0][0] == '@' {
		user, _ := strconv.Atoi(args[0])
		userIDs, err := services.FindUser(user)
		if userIDs == nil || err != nil {
			_, _ = m.ReplyHTML(trans.GetString(m.Chat.Id, "error/usernotindb"))
			return 0, ""
		} else {
			res := strings.SplitN(m.Text, " ", 2)
			if len(res) >= 3 {
				text = res[2]
				return userIDs.UserID, text
			}
		}
	}

	return 0, ""
}