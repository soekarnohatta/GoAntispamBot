package private

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
	"github.com/sirupsen/logrus"
)

func start(b ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat

	startButtons := [][]ext.InlineKeyboardButton{make([]ext.InlineKeyboardButton, 2)}
	startButtons[0][0] = ext.InlineKeyboardButton{
		Text:         "📝 Help",
		CallbackData: fmt.Sprintf("start(%v)", "about"),
	}
	startButtons[0][1] = ext.InlineKeyboardButton{
		Text: "🔗 Add Me To Your Groups",
		Url:  "https://t.me/PolyesterBot?startgroup=new",
	}

	txtStart := function.GetString(chat.Id, "modules/private/pm.go:start")
	replyMsg := b.NewSendableMessage(chat.Id, txtStart)
	replyMsg.ParseMode = "Markdown"
	replyMsg.ReplyMarkup = &ext.InlineKeyboardMarkup{&startButtons}
	replyMsg.ReplyToMessageId = msg.MessageId
	_, err := replyMsg.Send()
	return err
}

func LoadPm(u *gotgbot.Updater) {
	defer logrus.Info("PM Module Loaded...")
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("start", []rune{'/', '.'}, start))
}
