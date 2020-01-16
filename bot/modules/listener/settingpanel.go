package listener

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/caching"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/chat_status"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
)

func panel(b ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat_status.RequireSupergroup(chat, msg) {
		if chat_status.IsUserAdmin(chat, user.Id) {
			replyText, replyButtons := mainMenu(chat.Id)
			reply := b.NewSendableMessage(chat.Id, replyText)
			reply.ReplyMarkup = &ext.InlineKeyboardMarkup{&replyButtons}
			reply.ParseMode = parsemode.Html
			reply.ReplyToMessageId = msg.MessageId
			_, err := reply.Send()
			err_handler.HandleErr(err)
			return err
		}
		return nil
	}
	return nil
}

func backQuery(b ext.Bot, u *gotgbot.Update) error {
	msg := u.CallbackQuery
	user := msg.From
	chat := msg.Message.Chat

	if msg != nil {
		if chat.Type == "supergroup" {
			if chat_status.IsUserAdmin(chat, user.Id) == true {
				replyText, replyButtons := mainMenu(chat.Id)
				_, err := b.EditMessageTextMarkup(chat.Id, msg.Message.MessageId, replyText, parsemode.Html,
					&ext.InlineKeyboardMarkup{&replyButtons})
				return err
			}
		}
		_, err := b.AnswerCallbackQuery(msg.Id)
		return err
	}
	return gotgbot.ContinueGroups{}
}

func closeQuery(b ext.Bot, u *gotgbot.Update) error {
	msg := u.CallbackQuery
	user := msg.From
	chat := msg.Message.Chat

	if msg != nil {
		if chat.Type == "supergroup" {
			if chat_status.IsUserAdmin(chat, user.Id) == true {
				_, err := msg.Message.Delete()
				return err
			}
		} else {
			_, err := b.AnswerCallbackQuery(msg.Id)
			return err
		}
	}
	return gotgbot.ContinueGroups{}
}

func settingQuery(b ext.Bot, u *gotgbot.Update) error {
	msg := u.CallbackQuery
	user := msg.From
	chat := msg.Message.Chat

	if msg != nil {
		if chat.Type == "supergroup" {
			if chat_status.IsUserAdmin(chat, user.Id) == true {
				switch msg.Data {
				case "mk_utama":
					replyText, _, kn := mainControlMenu(chat.Id)
					_, err := b.EditMessageTextMarkup(chat.Id, msg.Message.MessageId, replyText,
						parsemode.Html, &ext.InlineKeyboardMarkup{&kn})
					return err
				case "mk_reset":
					sql.UpdatePicture(chat.Id, "true", "mute", "-", "true")
					sql.UpdateUsername(chat.Id, "true", "mute", "-", "true")
					sql.UpdateEnforceGban(chat.Id, "true")
					sql.UpdateVerify(chat.Id, "true", "-", "true")
					sql.UpdateSetting(chat.Id, "5m", "true")
					sql.UpdateLang(chat.Id, "id")
					caching.REDIS.Set(fmt.Sprintf("lang_%v", chat.Id), "en", 0)
					caching.REDIS.BgSave()

					err := updateUserControl(b, u)
					return err
				case "mk_spam":
					replyText, kn := mainSpamMenu(chat.Id)
					_, err := b.EditMessageTextMarkup(chat.Id, msg.Message.MessageId,
						replyText, parsemode.Html, &ext.InlineKeyboardMarkup{&kn})
					return err
				default:
					repCb := b.NewSendableAnswerCallbackQuery(msg.Id)
					repCb.CacheTime = 200
					_, err := repCb.Send()
					return err
				}
			}
		}
	}
	return gotgbot.ContinueGroups{}
}

func spamControlQuery(b ext.Bot, u *gotgbot.Update) error {
	var err error = nil
	msg := u.CallbackQuery
	user := msg.From
	chat := msg.Message.Chat

	if msg != nil {
		if chat.Type == "supergroup" {
			if chat_status.IsUserAdmin(chat, user.Id) == true {
				pattern, _ := regexp.MatchString("^mo_toggle$", msg.Data)
				if pattern {
					if strings.Split(msg.Data, "_toggle")[0] == "mo" {
						if sql.GetEnforceGban(chat.Id).Option == "true" {
							sql.UpdateEnforceGban(chat.Id, "false")
						} else {
							sql.UpdateEnforceGban(chat.Id, "true")
						}

						replyText, replyButtons := mainSpamMenu(chat.Id)
						_, err = b.EditMessageTextMarkup(chat.Id, msg.Message.MessageId,
							replyText, parsemode.Html, &ext.InlineKeyboardMarkup{&replyButtons})
						return err
					}
				}
			}
		}
	}
	return nil
}

func userControlQuery(b ext.Bot, u *gotgbot.Update) error {
	msg := u.CallbackQuery
	user := msg.From
	chat := msg.Message.Chat

	if msg != nil {
		if chat.Type == "supergroup" {
			if chat_status.IsUserAdmin(chat, user.Id) {
				// Grab Data From DB
				username := sql.GetUsername(chat.Id)
				profilePicture := sql.GetPicture(chat.Id)
				time := sql.GetSetting(chat.Id)
				ver := sql.GetVerify(chat.Id)
				warn := sql.GetWarnSetting(strconv.Itoa(chat.Id))
				aspam := sql.GetAntispam(chat.Id)

				// Separating Queries
				patternDel, _ := regexp.MatchString("^m[cdef]_del$", msg.Data)
				patternUsername, _ := regexp.MatchString("^mc_(kick|ban|mute|warn)$", msg.Data)
				patternPicture, _ := regexp.MatchString("^md_(kick|ban|mute|warn)$", msg.Data)
				patternToggle, _ := regexp.MatchString("^m[cdeo]_toggle$", msg.Data)
				patternTime, _ := regexp.MatchString("^mf_(plus|minus|duration|waktu)$", msg.Data)
				patternWarn, _ := regexp.MatchString("^mb_(plus|minus|warn)$", msg.Data)
				patternProtect, _ := regexp.MatchString("^as_(toggle|fwd|link|ar)$", msg.Data)

				if patternUsername == true {
					// Username Control Panel
					sql.UpdateUsername(chat.Id, username.Option, strings.Split(msg.Data, "mc_")[1], "-", username.Deletion)
					err := updateUserControl(b, u)
					return err
				} else if patternPicture == true {
					// Profile Photo Control Panel
					sql.UpdatePicture(chat.Id, profilePicture.Option, strings.Split(msg.Data, "md_")[1], "-", profilePicture.Deletion)
					err := updateUserControl(b, u)
					return err
				} else if patternTime == true {
					// Time Control Panel
					switch strings.Split(msg.Data, "mf_")[1] {
					case "duration":
						lastLetter := time.Time[len(time.Time)-1:]
						lastLetter = strings.ToLower(lastLetter)
						re := regexp.MustCompile(`[mhd]`)

						switch lastLetter {
						case "m":
							sql.UpdateSetting(chat.Id, fmt.Sprintf("%vh", re.Split(time.Time, -1)[0]), time.Deletion)
						case "h":
							sql.UpdateSetting(chat.Id, fmt.Sprintf("%vd", re.Split(time.Time, -1)[0]), time.Deletion)
						case "d":
							sql.UpdateSetting(chat.Id, fmt.Sprintf("%vm", re.Split(time.Time, -1)[0]), time.Deletion)
						}

						err := updateUserControl(b, u)
						return err
					case "plus":
						lastLetter := time.Time[len(time.Time)-1:]
						lastLetter = strings.ToLower(lastLetter)
						duration := time.Time[:len(time.Time)-1]
						finalDur, _ := strconv.Atoi(duration)
						finalDur++
						sql.UpdateSetting(chat.Id, fmt.Sprintf("%v%v", finalDur, lastLetter), time.Deletion)
						err := updateUserControl(b, u)
						return err
					case "minus":
						lastLetter := time.Time[len(time.Time)-1:]
						lastLetter = strings.ToLower(lastLetter)
						duration := time.Time[:len(time.Time)-1]
						finalDur, _ := strconv.Atoi(duration)
						finalDur--
						sql.UpdateSetting(chat.Id, fmt.Sprintf("%v%v", finalDur, lastLetter), time.Deletion)
						err := updateUserControl(b, u)
						return err
					case "waktu":
						replyCallback := b.NewSendableAnswerCallbackQuery(msg.Id)
						replyCallback.Text = "Time settings for all actions."
						replyCallback.ShowAlert = true
						replyCallback.CacheTime = 200
						_, err := replyCallback.Send()
						return err
					}
				} else if patternWarn == true {
					// Warn Settings
					if strings.Split(msg.Data, "mb_")[1] == "plus" {
						num := warn + 1
						sql.SetWarnLimit(strconv.Itoa(chat.Id), num)
						err := updateUserControl(b, u)
						return err
					} else if strings.Split(msg.Data, "mb_")[1] == "minus" {
						num := warn - 1
						sql.SetWarnLimit(strconv.Itoa(chat.Id), num)
						err := updateUserControl(b, u)
						return err
					} else if strings.Split(msg.Data, "mb_")[1] == "warn" {
						replyCallback := b.NewSendableAnswerCallbackQuery(msg.Id)
						replyCallback.Text = "Punishment settings for warn action."
						replyCallback.ShowAlert = true
						replyCallback.CacheTime = 200
						_, err := replyCallback.Send()
						return err
					}
				} else if patternToggle == true {
					// On/Off Toggles
					switch strings.Split(msg.Data, "_toggle")[0] {
					case "mc":
						if username.Option == "true" {
							sql.UpdateUsername(chat.Id, "false", username.Action, "-", username.Deletion)
						} else {
							sql.UpdateUsername(chat.Id, "true", username.Action, "-", username.Deletion)
						}
					case "md":
						if profilePicture.Option == "true" {
							sql.UpdatePicture(chat.Id, "false", username.Action, "-", username.Deletion)
						} else {
							sql.UpdatePicture(chat.Id, "true", username.Action, "-", username.Deletion)
						}
					case "me":
						if ver.Option == "true" {
							sql.UpdateVerify(chat.Id, "false", "-", ver.Deletion)
						} else {
							sql.UpdateVerify(chat.Id, "true", "-", ver.Deletion)
						}
					}

					err := updateUserControl(b, u)
					return err
				} else if patternDel == true {
					// On/Off Deletion
					switch strings.Split(msg.Data, "_del")[0] {
					case "mc":
						if username.Deletion == "true" {
							sql.UpdateUsername(chat.Id, username.Option, username.Action, "-", "false")
						} else {
							sql.UpdateUsername(chat.Id, username.Option, username.Action, "-", "true")
						}
					case "md":
						if profilePicture.Deletion == "true" {
							sql.UpdatePicture(chat.Id, profilePicture.Option, profilePicture.Action, "-", "false")
						} else {
							sql.UpdatePicture(chat.Id, profilePicture.Option, profilePicture.Action, "-", "true")
						}
					case "me":
						if ver.Deletion == "true" {
							sql.UpdateVerify(chat.Id, ver.Option, "-", "false")
						} else {
							sql.UpdateVerify(chat.Id, ver.Option, "-", "true")
						}
					case "mf":
						if time.Deletion == "true" {
							sql.UpdateSetting(chat.Id, time.Time, "false")
						} else {
							sql.UpdateSetting(chat.Id, time.Time, "true")
						}
					}

					err := updateUserControl(b, u)
					return err
				} else if patternProtect == true {
					switch strings.Split(msg.Data, "as_")[1] {
					case "toggle":
						if aspam.Deletion == "true" {
							sql.UpdateAntispam(chat.Id, aspam.Arabs, "false", aspam.Forward, aspam.Link)
						} else {
							sql.UpdateAntispam(chat.Id, aspam.Arabs, "true", aspam.Forward, aspam.Link)
						}
						err := updateUserControl(b, u)
						return err
					case "ar":
						if aspam.Arabs == "true" {
							sql.UpdateAntispam(chat.Id, "false", aspam.Deletion, aspam.Forward, aspam.Link)
						} else {
							sql.UpdateAntispam(chat.Id, "true", aspam.Deletion, aspam.Forward, aspam.Link)
						}
						err := updateUserControl(b, u)
						return err
					case "fwd":
						if aspam.Forward == "true" {
							sql.UpdateAntispam(chat.Id, aspam.Arabs, aspam.Deletion, "false", aspam.Link)
						} else {
							sql.UpdateAntispam(chat.Id, aspam.Arabs, aspam.Deletion, "true", aspam.Link)
						}
						err := updateUserControl(b, u)
						return err
					case "link":
						if aspam.Link == "true" {
							sql.UpdateAntispam(chat.Id, aspam.Arabs, aspam.Deletion, aspam.Forward, "false")
						} else {
							sql.UpdateAntispam(chat.Id, aspam.Arabs, aspam.Deletion, aspam.Forward, "true")
						}
						err := updateUserControl(b, u)
						return err
					}
				}
			}
		}
	}
	return gotgbot.ContinueGroups{}
}

func updateUserControl(b ext.Bot, u *gotgbot.Update) error {
	var err error = nil
	msg := u.CallbackQuery
	chat := msg.Message.Chat

	_, err = b.AnswerCallbackQuery(msg.Id)
	err_handler.HandleErr(err)

	opsiSama := "Bad Request: message is not modified: specified new message content and " +
		"reply markup are exactly the same as the current " +
		"content and reply markup of the message"

	replyText, _, replyButtons := mainControlMenu(chat.Id)
	_, err = b.EditMessageTextMarkup(chat.Id, msg.Message.MessageId,
		replyText, parsemode.Html, &ext.InlineKeyboardMarkup{&replyButtons})
	if err != nil {
		if err.Error() == opsiSama {
			_, err := b.AnswerCallbackQuery(msg.Id)
			replyCallback := b.NewSendableAnswerCallbackQuery(msg.Id)
			replyCallback.CacheTime = 200
			_, err = replyCallback.Send()
			return err
		}
		_, err := b.AnswerCallbackQuery(msg.Id)
		return err
	}
	_, err = b.AnswerCallbackQuery(msg.Id)
	return err
}

func mainControlMenu(chatId int) (string, [][]string, [][]ext.InlineKeyboardButton) {
	emoji := getEmoji(chatId)
	if emoji != nil {
		replyText := function.GetStringf(chatId, "modules/helpers/function.go:13",
			map[string]string{"1": emoji[0][0], "2": emoji[1][0], "3": emoji[2][0], "4": emoji[0][1], "5": emoji[1][1],
				"6": emoji[2][1], "7": emoji[0][2], "8": emoji[2][3], "9": emoji[3][0], "10": strconv.Itoa(sql.GetWarnSetting(strconv.Itoa(chatId))),
				"11": emoji[0][4], "12": emoji[5][0], "13": emoji[5][1], "14": emoji[5][2]})

		kn := make([][]ext.InlineKeyboardButton, 0)

		ki := make([]ext.InlineKeyboardButton, 6)
		ki[0] = ext.InlineKeyboardButton{Text: emoji[0][0], CallbackData: "mc_toggle"}
		ki[1] = ext.InlineKeyboardButton{Text: "🔇", CallbackData: "mc_mute"}
		ki[2] = ext.InlineKeyboardButton{Text: "🚷", CallbackData: "mc_kick"}
		ki[3] = ext.InlineKeyboardButton{Text: "⛔", CallbackData: "mc_ban"}
		ki[4] = ext.InlineKeyboardButton{Text: "❗", CallbackData: "mc_warn"}
		ki[5] = ext.InlineKeyboardButton{Text: "🗑", CallbackData: "mc_del"}
		kn = append(kn, ki)

		kd := make([]ext.InlineKeyboardButton, 6)
		kd[0] = ext.InlineKeyboardButton{Text: emoji[0][1], CallbackData: "md_toggle"}
		kd[1] = ext.InlineKeyboardButton{Text: "🔇", CallbackData: "md_mute"}
		kd[2] = ext.InlineKeyboardButton{Text: "🚷", CallbackData: "md_kick"}
		kd[3] = ext.InlineKeyboardButton{Text: "⛔", CallbackData: "md_ban"}
		kd[4] = ext.InlineKeyboardButton{Text: "❗", CallbackData: "md_warn"}
		kd[5] = ext.InlineKeyboardButton{Text: "🗑", CallbackData: "md_del"}
		kn = append(kn, kd)

		ke := make([]ext.InlineKeyboardButton, 4)
		ke[0] = ext.InlineKeyboardButton{Text: emoji[0][4], CallbackData: "as_toggle"}
		ke[1] = ext.InlineKeyboardButton{Text: "Ar", CallbackData: "as_ar"}
		ke[2] = ext.InlineKeyboardButton{Text: "🔗", CallbackData: "as_link"}
		ke[3] = ext.InlineKeyboardButton{Text: "➡", CallbackData: "as_fwd"}
		kn = append(kn, ke)

		kj := make([]ext.InlineKeyboardButton, 2)
		kj[0] = ext.InlineKeyboardButton{Text: emoji[0][2], CallbackData: "me_toggle"}
		kj[1] = ext.InlineKeyboardButton{Text: "🗑", CallbackData: "me_del"}
		kn = append(kn, kj)

		kk := make([]ext.InlineKeyboardButton, 3)
		kk[0] = ext.InlineKeyboardButton{Text: "❗", CallbackData: "mb_warn"}
		kk[1] = ext.InlineKeyboardButton{Text: "➕", CallbackData: "mb_plus"}
		kk[2] = ext.InlineKeyboardButton{Text: "➖", CallbackData: "mb_minus"}
		kn = append(kn, kk)

		ku := make([]ext.InlineKeyboardButton, 5)
		ku[0] = ext.InlineKeyboardButton{Text: "🕑", CallbackData: "mf_waktu"}
		ku[1] = ext.InlineKeyboardButton{Text: "➕", CallbackData: "mf_plus"}
		ku[2] = ext.InlineKeyboardButton{Text: "➖", CallbackData: "mf_minus"}
		ku[3] = ext.InlineKeyboardButton{Text: emoji[3][0], CallbackData: "mf_duration"}
		ku[4] = ext.InlineKeyboardButton{Text: "🗑", CallbackData: "mf_del"}
		kn = append(kn, ku)

		kg := make([]ext.InlineKeyboardButton, 2)
		kg[0] = ext.InlineKeyboardButton{Text: "🔙", CallbackData: "back_main"}
		kg[1] = ext.InlineKeyboardButton{Text: "✖", CallbackData: "close"}
		kn = append(kn, kg)

		return replyText, emoji, kn
	}
	return "", nil, nil
}

func mainSpamMenu(chatId int) (string, [][]ext.InlineKeyboardButton) {
	emoji := getEmoji(chatId)
	if emoji != nil {
		replyText := function.GetStringf(chatId, "modules/helpers/function.go:66", map[string]string{"1": emoji[0][3]})

		kn := make([][]ext.InlineKeyboardButton, 0)

		ki := make([]ext.InlineKeyboardButton, 1)
		ki[0] = ext.InlineKeyboardButton{Text: emoji[0][3], CallbackData: "mo_toggle"}
		kn = append(kn, ki)

		kg := make([]ext.InlineKeyboardButton, 2)
		kg[0] = ext.InlineKeyboardButton{Text: "🔙", CallbackData: "back_main"}
		kg[1] = ext.InlineKeyboardButton{Text: "✖", CallbackData: "close"}
		kn = append(kn, kg)

		return replyText, kn
	}
	return "", nil
}

func mainMenu(chatId int) (string, [][]ext.InlineKeyboardButton) {

	replyText := function.GetString(chatId, "modules/helpers/function.go:85")

	kn := make([][]ext.InlineKeyboardButton, 0)

	ki := make([]ext.InlineKeyboardButton, 2)
	ki[0] = ext.InlineKeyboardButton{Text: function.GetString(chatId, "modules/helpers/function.go:91"), CallbackData: "mk_utama"}
	ki[1] = ext.InlineKeyboardButton{Text: function.GetString(chatId, "modules/helpers/function.go:92"), CallbackData: "mk_spam"}
	kn = append(kn, ki)

	kz := make([]ext.InlineKeyboardButton, 2)
	kz[0] = ext.InlineKeyboardButton{Text: function.GetString(chatId, "modules/helpers/function.go:96"), CallbackData: "mk_media"}
	kz[1] = ext.InlineKeyboardButton{Text: function.GetString(chatId, "modules/helpers/function.go:97"), CallbackData: "mk_pesan"}
	kn = append(kn, kz)

	kd := make([]ext.InlineKeyboardButton, 1)
	kd[0] = ext.InlineKeyboardButton{Text: function.GetString(chatId, "modules/helpers/function.go:101"), CallbackData: "mk_reset"}
	kn = append(kn, kd)

	kk := make([]ext.InlineKeyboardButton, 1)
	kk[0] = ext.InlineKeyboardButton{Text: function.GetString(chatId, "modules/helpers/function.go:105"), CallbackData: "close"}
	kn = append(kn, kk)

	return replyText, kn
}

func getEmoji(chatId int) [][]string {
	chat := sql.GetUsername(chatId)
	pic := sql.GetPicture(chatId)
	ver := sql.GetVerify(chatId)
	tim := sql.GetSetting(chatId)
	spm := sql.GetEnforceGban(chatId)
	aspm := sql.GetAntispam(chatId)

	lastLetter := "m"

	lst := make([][]string, 0)
	opt := make([]string, 5)
	act := make([]string, 2)
	del := make([]string, 4)
	ti := make([]string, 1)
	gu := make([]string, 1)
	hrr := make([]string, 3)

	if tim != nil {
		lastLetter = tim.Time[len(tim.Time)-1:]
		lastLetter = strings.ToLower(lastLetter)
	}

	if chat != nil {
		if chat.Option == "true" {
			chat.Option = "🔵"
		} else {
			chat.Option = "⚪"
		}

		if chat.Deletion == "true" {
			chat.Deletion = "+ 🗑"
		} else {
			chat.Deletion = ""
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

		opt[0] = chat.Option
		act[0] = chat.Action
		del[0] = chat.Deletion
	} else {
		go sql.UpdateUsername(chatId, "true", "mute", "-", "true")
	}

	if pic != nil {
		if pic.Option == "true" {
			pic.Option = "🔵"
		} else {
			pic.Option = "⚪"
		}

		if pic.Deletion == "true" {
			pic.Deletion = "+ 🗑"
		} else {
			pic.Deletion = ""
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

		opt[1] = pic.Option
		act[1] = pic.Action
		del[1] = pic.Deletion
	} else {
		go sql.UpdatePicture(chatId, "true", "mute", "-", "true")
	}

	if ver != nil {
		if ver.Option == "true" {
			ver.Option = "🔵"
		} else {
			ver.Option = "⚪"
		}

		if ver.Deletion == "true" {
			ver.Deletion = "+ 🗑"
		} else {
			ver.Deletion = ""
		}

		opt[2] = ver.Option
		del[3] = ver.Deletion
	} else {
		go sql.UpdateVerify(chatId, "true", "-", "true")
	}

	if spm != nil {
		if spm.Option == "true" {
			spm.Option = "🔵"
		} else {
			spm.Option = "⚪"
		}

		opt[3] = spm.Option
	} else {
		go sql.UpdateEnforceGban(chatId, "true")
	}

	if tim != nil {
		if tim.Deletion == "true" {
			tim.Deletion = "+ 🗑"
		} else {
			tim.Deletion = ""
		}

		ti[0] = tim.Time
		del[2] = tim.Deletion
	} else {
		go sql.UpdateSetting(chatId, "5m", "true")
	}

	if aspm != nil {
		if aspm.Deletion == "true" {
			aspm.Deletion = "🔵"
		} else {
			aspm.Deletion = "⚪"
		}

		if aspm.Arabs == "true" {
			aspm.Arabs = "+ Ar"
		} else {
			aspm.Arabs = ""
		}

		if aspm.Link == "true" {
			aspm.Link = "+ 🔗"
		} else {
			aspm.Link = ""
		}

		if aspm.Forward == "true" {
			aspm.Forward = "+ ➡"
		} else {
			aspm.Forward = ""
		}

		opt[4] = aspm.Deletion
		hrr[0] = aspm.Arabs
		hrr[1] = aspm.Forward
		hrr[2] = aspm.Link
	}

	gu[0] = lastLetter

	lst = append(lst, opt, act, del, ti, gu, hrr)
	return lst
}

func LoadSettingListener(u *gotgbot.Updater) {
	defer logrus.Info("Setting Listener Loaded...")
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("settings", []rune{'/', '.'}, panel))
	u.Dispatcher.AddHandler(handlers.NewCallback(
		"^(m[cdefgb]|as)_(toggle|warn|kick|ban|mute|reset|plus|minus|duration|waktu|del|warn|fwd|ar|link)",
		userControlQuery))
	u.Dispatcher.AddHandler(handlers.NewCallback("mo_toggle", spamControlQuery))
	u.Dispatcher.AddHandler(handlers.NewCallback("mk_", settingQuery))
	u.Dispatcher.AddHandler(handlers.NewCallback("close", closeQuery))
	u.Dispatcher.AddHandler(handlers.NewCallback("back_", backQuery))
}
