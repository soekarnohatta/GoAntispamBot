package events

import (
	"GoAntispamBot/bot/helpers/message"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"html"
	"strconv"
	"strings"

	"GoAntispamBot/bot/helpers/trans"
	"GoAntispamBot/bot/model"
	"GoAntispamBot/bot/providers/telegramProvider"
	"GoAntispamBot/bot/services/welcomeService"
)

//var VALID_WELCOME_FORMATTERS = []string{"first", "last", "fullname", "username", "id", "count", "chatname", "mention"}

// EnumFuncMap map of welcome type to function to execute
var EnumFuncMap = map[int]func(ext.Bot, int, string) (*ext.Message, error){
	welcomeService.TEXT:        ext.Bot.SendMessage,
	welcomeService.BUTTON_TEXT: ext.Bot.SendMessage,
	welcomeService.STICKER:     ext.Bot.SendStickerStr,
	welcomeService.DOCUMENT:    ext.Bot.SendDocumentStr,
	welcomeService.PHOTO:       ext.Bot.SendPhotoStr,
	welcomeService.AUDIO:       ext.Bot.SendAudioStr,
	welcomeService.VOICE:       ext.Bot.SendVoiceStr,
	welcomeService.VIDEO:       ext.Bot.SendVideoStr,
}

// NewMember is a collection of handlers that are handling new members
type NewMember struct {
	TelegramProvider telegramProvider.TelegramProvider
}

func Send(bot ext.Bot, u *gotgbot.Update, message string, keyboard *ext.InlineKeyboardMarkup, backupMessage string, reply bool) *ext.Message {
	msg := bot.NewSendableMessage(u.EffectiveChat.Id, message)
	msg.ParseMode = parsemode.Html
	if reply {
		msg.ReplyToMessageId = u.EffectiveMessage.MessageId
	}
	msg.ReplyMarkup = keyboard
	m, err := msg.Send()
	if err != nil {
		m, _ = u.EffectiveMessage.ReplyText(backupMessage + trans.GetString(msg.ChatId, "error/invalidwelc"))
	}
	return m
}

// NewMemberHandler handles all incoming updates of new members.
func (r NewMember) NewMemberHandler(bot ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	newMembers := u.EffectiveMessage.NewChatMembers
	welcPrefs := welcomeService.GetWelcomePrefs(chat.Id)
	var firstName = ""
	var fullName = ""
	var username = ""
	var res = ""
	var keyb = make([][]ext.InlineKeyboardButton, 0)

	if welcPrefs.ShouldWelcome {
		for _, mem := range newMembers {
			if mem.Id == bot.Id {
				continue
			}

			if welcPrefs.WelcomeType != welcomeService.TEXT && welcPrefs.WelcomeType != welcomeService.BUTTON_TEXT {
				_, err := EnumFuncMap[welcPrefs.WelcomeType](bot, chat.Id, welcPrefs.CustomWelcome)
				if err != nil {
					return err
				}
			}
			firstName = mem.FirstName
			if len(mem.FirstName) <= 0 {
				firstName = "Unknown"
			}

			if welcPrefs.CustomWelcome != "" {
				if mem.LastName != "" {
					fullName = firstName + " " + mem.LastName
				} else {
					fullName = firstName
				}
				count, _ := chat.GetMembersCount()
				mention := message.MentionHtml(mem.Id, firstName)

				if mem.Username != "" {
					username = "@" + html.EscapeString(mem.Username)
				} else {
					username = mention
				}

				r := strings.NewReplacer(
					"{first}", html.EscapeString(firstName),
					"{last}", html.EscapeString(mem.LastName),
					"{fullname}", html.EscapeString(fullName),
					"{username}", username,
					"{mention}", mention,
					"{count}", strconv.Itoa(count),
					"{chatname}", html.EscapeString(chat.Title),
					"{id}", strconv.Itoa(mem.Id),
					"{rules}", "",
				)
				res = r.Replace(welcPrefs.CustomWelcome)
				buttons := welcomeService.GetWelcomeButtons(chat.Id)
				if strings.Contains(welcPrefs.CustomWelcome, "{rules}") {
					rulesButton := model.WelcomeButton{
						ChatId:   u.EffectiveChat.Id,
						Name:     "Rules",
						Url:      fmt.Sprintf("t.me/%v?start=%v", bot.UserName, u.EffectiveChat.Id),
						SameLine: false,
					}
					buttons = append(buttons, rulesButton)
				}
				keyb = message.BuildWelcomeKeyboard(buttons)
			} else {
				r := strings.NewReplacer("{first}", firstName)
				res = r.Replace(welcomeService.DefaultWelcome)
			}

			if welcPrefs.ShouldMute {
				if !welcomeService.IsUserHuman(mem.Id, chat.Id) {
					if !welcomeService.HasUserClickedButton(mem.Id, chat.Id) {
						_, _ = bot.RestrictChatMember(chat.Id, mem.Id)
					}
				}
				kb := make([]ext.InlineKeyboardButton, 1)
				kb[0] = ext.InlineKeyboardButton{Text: "Click here to prove you're human", CallbackData: "unmute"}
				keyb = append(keyb, kb)
			}

			keyboard := &ext.InlineKeyboardMarkup{InlineKeyboard: &keyb}
			r := strings.NewReplacer("{first}", firstName)
			sent := Send(bot, u, res, keyboard, r.Replace(welcomeService.DefaultWelcome), !welcPrefs.DelJoined)

			if welcPrefs.CleanWelcome != 0 {
				_, _ = bot.DeleteMessage(chat.Id, welcPrefs.CleanWelcome)
				welcomeService.SetCleanWelcome(chat.Id, sent.MessageId)
			}

			if welcPrefs.DelJoined {
				_, _ = u.EffectiveMessage.Delete()
			}
		}
	}
	return nil
}

// UnmuteCallback is a handler to handle unmute callback from verification button
// this should help to prevent spammer.
func (r NewMember) UnmuteCallback(bot ext.Bot, u *gotgbot.Update) error {
	query := u.CallbackQuery
	effectiveUser := u.EffectiveUser
	chat := u.EffectiveChat

	if !welcomeService.IsUserHuman(effectiveUser.Id, chat.Id) {
		if !welcomeService.HasUserClickedButton(effectiveUser.Id, chat.Id) {
			_, err := bot.UnRestrictChatMember(chat.Id, effectiveUser.Id)
			if err != nil {
				return err
			}
			welcomeService.UserClickedButton(effectiveUser.Id, chat.Id)
			_, _ = bot.AnswerCallbackQueryText(query.Id, "You've proved that you are a human! "+
				"You can now talk in this group.", false)
			return nil
		}
	}

	_, _ = bot.AnswerCallbackQueryText(query.Id, "This action is invalid for you.", false)
	return gotgbot.EndGroups{}
}
