package setting

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/caching"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/chat_status"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/extraction"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
)

func panel(b ext.Bot, u *gotgbot.Update) error {
	var err error
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat_status.RequireSupergroup(chat, msg) == true {
		if chat_status.IsUserAdmin(chat, user.Id) == true {
			teks, _, kn := mainMenu(chat.Id)
			reply := b.NewSendableMessage(chat.Id, teks)
			reply.ReplyMarkup = &ext.InlineKeyboardMarkup{&kn}
			reply.ParseMode = parsemode.Html
			reply.ReplyToMessageId = msg.MessageId
			_, err = reply.Send()
			return err
		}
	}
	return nil
}

func backquery(b ext.Bot, u *gotgbot.Update) error {
	var err error
	msg := u.CallbackQuery
	user := msg.From
	chat := msg.Message.Chat

	if msg != nil {
		if chat.Type == "supergroup" {
			if chat_status.IsUserAdmin(chat, user.Id) == true {
				teks, _, kn := mainMenu(chat.Id)
				_, err = b.EditMessageTextMarkup(chat.Id, msg.Message.MessageId, teks, parsemode.Html,
					&ext.InlineKeyboardMarkup{&kn})
				return err
			}
		} else {
			_, err = b.AnswerCallbackQuery(msg.Id)
			return err
		}
	}
	return gotgbot.ContinueGroups{}
}

func closequery(b ext.Bot, u *gotgbot.Update) error {
	var err error
	msg := u.CallbackQuery
	user := msg.From
	chat := msg.Message.Chat

	if msg != nil {
		if chat.Type == "supergroup" {
			if chat_status.IsUserAdmin(chat, user.Id) == true {
				_, err = msg.Message.Delete()
				return err
			}
		} else {
			_, err = b.AnswerCallbackQuery(msg.Id)
			return err
		}
	}
	return gotgbot.ContinueGroups{}
}

func settingquery(b ext.Bot, u *gotgbot.Update) error {
	var err error
	msg := u.CallbackQuery
	user := msg.From
	chat := msg.Message.Chat

	if msg != nil {
		if chat.Type == "supergroup" {
			if chat_status.IsUserAdmin(chat, user.Id) == true {
				if msg.Data == "mk_utama" {
					teks, _, kn := mainControlMenu(chat.Id)
					_, err = b.EditMessageTextMarkup(chat.Id, msg.Message.MessageId,
						teks, "HTML", &ext.InlineKeyboardMarkup{&kn})
					return err
				} else if msg.Data == "mk_reset" {
					err = sql.UpdatePicture(chat.Id, "true", "mute", "-", "true")
					err_handler.HandleCbErr(b, u, err)
					err = sql.UpdateUsername(chat.Id, "true", "mute", "-", "true")
					err_handler.HandleCbErr(b, u, err)
					err = sql.UpdateEnforceGban(chat.Id, "true")
					err_handler.HandleCbErr(b, u, err)
					err = sql.UpdateVerify(chat.Id, "true", "-", "true")
					err_handler.HandleCbErr(b, u, err)
					err = sql.UpdateSetting(chat.Id, "5m", "true")
					err_handler.HandleCbErr(b, u, err)
					err = sql.UpdateLang(chat.Id, "id")
					err_handler.HandleCbErr(b, u, err)
					caching.REDIS.Set(fmt.Sprintf("lang_%v", chat.Id), "en", 0)
					caching.REDIS.BgSave()

					err = updateusercontrol(b, u)
					return err
				} else if msg.Data == "mk_spam" {
					teks, _, kn := mainSpamMenu(chat.Id)
					_, err = b.EditMessageTextMarkup(chat.Id, msg.Message.MessageId,
						teks, "HTML", &ext.InlineKeyboardMarkup{&kn})
					return err
				} else {
					_, err = b.AnswerCallbackQuery(msg.Id)
				}
			}
		}
	}
	return gotgbot.ContinueGroups{}
}

func spamcontrolquery(b ext.Bot, u *gotgbot.Update) error {
	var err error
	msg := u.CallbackQuery
	user := msg.From
	chat := msg.Message.Chat

	if msg != nil {
		if chat.Type == "supergroup" {
			if chat_status.IsUserAdmin(chat, user.Id) == true {
				g, _ := regexp.MatchString("^mo_toggle$", msg.Data)
				if g {
					if strings.Split(msg.Data, "_toggle")[0] == "mo" {
						if sql.GetEnforceGban(chat.Id).Option == "true" {
							err = sql.UpdateEnforceGban(chat.Id, "false")
							err_handler.HandleCbErr(b, u, err)
						} else {
							err = sql.UpdateEnforceGban(chat.Id, "true")
							err_handler.HandleCbErr(b, u, err)
						}
						teks, _, kn := mainSpamMenu(chat.Id)
						_, err = b.EditMessageTextMarkup(chat.Id, msg.Message.MessageId,
							teks, "HTML", &ext.InlineKeyboardMarkup{&kn})
						return err
					}
				}
			}
		}
		return nil
	}
	return nil
}

func usercontrolquery(b ext.Bot, u *gotgbot.Update) error {
	var err error
	msg := u.CallbackQuery
	user := msg.From
	chat := msg.Message.Chat

	if msg != nil {
		if chat.Type == "supergroup" {
			if chat_status.IsUserAdmin(chat, user.Id) == true {
				// Grab Data From DB
				username := sql.GetUsername(chat.Id)
				fotoprofil := sql.GetPicture(chat.Id)
				waktu := sql.GetSetting(chat.Id)
				ver := sql.GetVerify(chat.Id)
				warn := sql.GetWarnSetting(strconv.Itoa(chat.Id))

				// Separating Queries
				z, _ := regexp.MatchString("^m[cdef]_del$", msg.Data)
				a, _ := regexp.MatchString("^mc_(kick|ban|mute|warn)$", msg.Data)
				f, _ := regexp.MatchString("^md_(kick|ban|mute|warn)$", msg.Data)
				g, _ := regexp.MatchString("^m[cdeo]_toggle$", msg.Data)
				d, _ := regexp.MatchString("^mf_(plus|minus|duration|waktu)$", msg.Data)
				k, _ := regexp.MatchString("^mb_(plus|minus|warn)$", msg.Data)

				// Username Control Panel
				if a == true {
					err = sql.UpdateUsername(chat.Id, username.Option, strings.Split(msg.Data, "mc_")[1], "-", username.Deletion)
					err_handler.HandleCbErr(b, u, err)
					err = updateusercontrol(b, u)
					return err
				} else if f == true {
					// Profile Photo Control Panel
					err = sql.UpdatePicture(chat.Id, fotoprofil.Option, strings.Split(msg.Data, "md_")[1], "-", fotoprofil.Deletion)
					err_handler.HandleCbErr(b, u, err)
					err = updateusercontrol(b, u)
					return err
				} else if d == true {
					// Time Control Panel
					if strings.Split(msg.Data, "mf_")[1] == "duration" {
						lastLetter := waktu.Time[len(waktu.Time)-1:]
						lastLetter = strings.ToLower(lastLetter)
						re := regexp.MustCompile(`[mhd]`)

						t := waktu.Time[:len(waktu.Time)-1]
						_, err := strconv.Atoi(t)
						if err != nil {
							_, err := b.AnswerCallbackQueryText(msg.Id,
								"❌ Invalid time amount specified.", true)
							return err
						}

						if lastLetter == "m" {
							err = sql.UpdateSetting(chat.Id, fmt.Sprintf("%vh", re.Split(waktu.Time, -1)[0]), waktu.Deletion)
							err_handler.HandleCbErr(b, u, err)
						} else if lastLetter == "h" {
							err = sql.UpdateSetting(chat.Id, fmt.Sprintf("%vd", re.Split(waktu.Time, -1)[0]), waktu.Deletion)
							err_handler.HandleCbErr(b, u, err)
						} else if lastLetter == "d" {
							err = sql.UpdateSetting(chat.Id, fmt.Sprintf("%vm", re.Split(waktu.Time, -1)[0]), waktu.Deletion)
							err_handler.HandleCbErr(b, u, err)
						}

						err = updateusercontrol(b, u)
						return err
					} else if strings.Split(msg.Data, "mf_")[1] == "plus" {
						lastLetter := waktu.Time[len(waktu.Time)-1:]
						lastLetter = strings.ToLower(lastLetter)

						t := waktu.Time[:len(waktu.Time)-1]
						j, err := strconv.Atoi(t)
						if err != nil {
							_, err := b.AnswerCallbackQueryText(msg.Id,
								"❌ Invalid time amount specified.", true)
							return err
						}
						j++

						err = sql.UpdateSetting(chat.Id, fmt.Sprintf("%v%v", j, lastLetter), waktu.Deletion)
						err_handler.HandleCbErr(b, u, err)
						err = updateusercontrol(b, u)
						return err
					} else if strings.Split(msg.Data, "mf_")[1] == "minus" {
						lastLetter := waktu.Time[len(waktu.Time)-1:]
						lastLetter = strings.ToLower(lastLetter)
						if strings.ContainsAny(lastLetter, "m & d & h") {
							t := waktu.Time[:len(waktu.Time)-1]
							j, err := strconv.Atoi(t)
							err_handler.HandleCbErr(b, u, err)
							j--

							err = sql.UpdateSetting(chat.Id, fmt.Sprintf("%v%v", j, lastLetter), waktu.Deletion)
							err_handler.HandleCbErr(b, u, err)
							err = updateusercontrol(b, u)
							return err
						}
					} else if strings.Split(msg.Data, "mf_")[1] == "waktu" {
						_, err := b.AnswerCallbackQueryText(msg.Id,
							"🔄 Mengatur tenggat waktu untuk semua aksi.", true)
						return err
					}
				} else if k == true {
					if strings.Split(msg.Data, "mb_")[1] == "plus" {
						num := warn + 1
						sql.SetWarnLimit(strconv.Itoa(chat.Id), num)
						err = updateusercontrol(b, u)
						return err
					} else if strings.Split(msg.Data, "mb_")[1] == "minus" {
						num := warn - 1
						sql.SetWarnLimit(strconv.Itoa(chat.Id), num)
						err = updateusercontrol(b, u)
						return err
					} else if strings.Split(msg.Data, "mb_")[1] == "warn" {
						_, err := b.AnswerCallbackQueryText(msg.Id,
							"🔄 Mengatur hukuman untuk warn.", true)
						return err
					}
				} else if g == true {
					// On/Off Toggles
					if strings.Split(msg.Data, "_toggle")[0] == "mc" {
						if username.Option == "true" {
							err = sql.UpdateUsername(chat.Id, "false", username.Action, "-", username.Deletion)
							err_handler.HandleCbErr(b, u, err)
						} else {
							err = sql.UpdateUsername(chat.Id, "true", username.Action, "-", username.Deletion)
							err_handler.HandleCbErr(b, u, err)
						}
					} else if strings.Split(msg.Data, "_toggle")[0] == "md" {
						if fotoprofil.Option == "true" {
							err = sql.UpdatePicture(chat.Id, "false", username.Action, "-", username.Deletion)
							err_handler.HandleCbErr(b, u, err)
						} else {
							err = sql.UpdatePicture(chat.Id, "true", username.Action, "-", username.Deletion)
							err_handler.HandleCbErr(b, u, err)
						}
					} else if strings.Split(msg.Data, "_toggle")[0] == "me" {
						if ver.Option == "true" {
							err = sql.UpdateVerify(chat.Id, "false", "-", ver.Deletion)
							err_handler.HandleCbErr(b, u, err)
						} else {
							err = sql.UpdateVerify(chat.Id, "true", "-", ver.Deletion)
							err_handler.HandleCbErr(b, u, err)
						}
					}

					err = updateusercontrol(b, u)
					return err
				} else if z == true {
					// On/Off Deletion
					if strings.Split(msg.Data, "_del")[0] == "mc" {
						if username.Deletion == "true" {
							err = sql.UpdateUsername(chat.Id, username.Option, username.Action, "-", "false")
							err_handler.HandleCbErr(b, u, err)
						} else {
							err = sql.UpdateUsername(chat.Id, username.Option, username.Action, "-", "true")
							err_handler.HandleCbErr(b, u, err)
						}
					} else if strings.Split(msg.Data, "_del")[0] == "md" {
						if fotoprofil.Deletion == "true" {
							err = sql.UpdatePicture(chat.Id, fotoprofil.Option, fotoprofil.Action, "-", "false")
							err_handler.HandleCbErr(b, u, err)
						} else {
							err = sql.UpdatePicture(chat.Id, fotoprofil.Option, fotoprofil.Action, "-", "true")
							err_handler.HandleCbErr(b, u, err)
						}
					} else if strings.Split(msg.Data, "_del")[0] == "me" {
						if ver.Deletion == "true" {
							err = sql.UpdateVerify(chat.Id, ver.Option, "-", "false")
							err_handler.HandleCbErr(b, u, err)
						} else {
							err = sql.UpdateVerify(chat.Id, ver.Option, "-", "true")
							err_handler.HandleCbErr(b, u, err)
						}
					} else if strings.Split(msg.Data, "_del")[0] == "mf" {
						if waktu.Deletion == "true" {
							err = sql.UpdateSetting(chat.Id, waktu.Time, "false")
							err_handler.HandleCbErr(b, u, err)
						} else {
							err = sql.UpdateSetting(chat.Id, waktu.Time, "true")
							err_handler.HandleCbErr(b, u, err)
						}
					}

					err = updateusercontrol(b, u)
					return err
				}
			}
		}
	}
	return gotgbot.ContinueGroups{}
}

func updateusercontrol(b ext.Bot, u *gotgbot.Update) error {
	var err error
	msg := u.CallbackQuery
	chat := msg.Message.Chat

	_, err = b.AnswerCallbackQuery(msg.Id)
	err_handler.HandleErr(err)

	opsisama := "Bad Request: message is not modified: specified new message content and " +
		"reply markup are exactly the same as a current " +
		"content and reply markup of the message"

	teks, _, kn := mainControlMenu(chat.Id)
	_, err = b.EditMessageTextMarkup(chat.Id, msg.Message.MessageId,
		teks, "HTML", &ext.InlineKeyboardMarkup{&kn})
	if err != nil {
		if err.Error() == opsisama {
			_, err := b.AnswerCallbackQuery(msg.Id)
			return err
		}
		_, err := b.AnswerCallbackQuery(msg.Id)
		return err
	}
	_, err = b.AnswerCallbackQuery(msg.Id)
	return err
}

func mainControlMenu(chatId int) (string, [][]string, [][]ext.InlineKeyboardButton) {
	a := extraction.GetEmoji(chatId)
	if a != nil {
		teks := function.GetStringf(chatId, "modules/helpers/function.go:13",
			map[string]string{"1": a[0][0], "2": a[1][0], "3": a[2][0], "4": a[0][1], "5": a[1][1], "6": a[2][1], "7": a[0][2],
				"8": a[2][3], "9": a[3][0], "10": strconv.Itoa(sql.GetWarnSetting(strconv.Itoa(chatId)))})

		kn := make([][]ext.InlineKeyboardButton, 0)

		ki := make([]ext.InlineKeyboardButton, 6)
		ki[0] = ext.InlineKeyboardButton{Text: a[0][0], CallbackData: "mc_toggle"}
		ki[1] = ext.InlineKeyboardButton{Text: "🔇", CallbackData: "mc_mute"}
		ki[2] = ext.InlineKeyboardButton{Text: "🚷", CallbackData: "mc_kick"}
		ki[3] = ext.InlineKeyboardButton{Text: "⛔", CallbackData: "mc_ban"}
		ki[4] = ext.InlineKeyboardButton{Text: "❗", CallbackData: "mc_warn"}
		ki[5] = ext.InlineKeyboardButton{Text: "🗑", CallbackData: "mc_del"}
		kn = append(kn, ki)

		kd := make([]ext.InlineKeyboardButton, 6)
		kd[0] = ext.InlineKeyboardButton{Text: a[0][1], CallbackData: "md_toggle"}
		kd[1] = ext.InlineKeyboardButton{Text: "🔇", CallbackData: "md_mute"}
		kd[2] = ext.InlineKeyboardButton{Text: "🚷", CallbackData: "md_kick"}
		kd[3] = ext.InlineKeyboardButton{Text: "⛔", CallbackData: "md_ban"}
		kd[4] = ext.InlineKeyboardButton{Text: "❗", CallbackData: "md_warn"}
		kd[5] = ext.InlineKeyboardButton{Text: "🗑", CallbackData: "md_del"}
		kn = append(kn, kd)

		kj := make([]ext.InlineKeyboardButton, 2)
		kj[0] = ext.InlineKeyboardButton{Text: a[0][2], CallbackData: "me_toggle"}
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
		ku[3] = ext.InlineKeyboardButton{Text: a[3][0], CallbackData: "mf_duration"}
		ku[4] = ext.InlineKeyboardButton{Text: "🗑", CallbackData: "mf_del"}
		kn = append(kn, ku)

		kg := make([]ext.InlineKeyboardButton, 2)
		kg[0] = ext.InlineKeyboardButton{Text: "🔙", CallbackData: "back_main"}
		kg[1] = ext.InlineKeyboardButton{Text: "✖", CallbackData: "close"}
		kn = append(kn, kg)

		return teks, a, kn
	}
	return "", nil, nil
}

func mainSpamMenu(chatId int) (string, [][]string, [][]ext.InlineKeyboardButton) {
	a := extraction.GetEmoji(chatId)
	if a != nil {
		teks := function.GetStringf(chatId, "modules/helpers/function.go:66", map[string]string{"1": a[0][3]})

		var kn = make([][]ext.InlineKeyboardButton, 0)

		ki := make([]ext.InlineKeyboardButton, 1)
		ki[0] = ext.InlineKeyboardButton{Text: a[0][3], CallbackData: "mo_toggle"}
		kn = append(kn, ki)

		kg := make([]ext.InlineKeyboardButton, 2)
		kg[0] = ext.InlineKeyboardButton{Text: "🔙", CallbackData: "back_main"}
		kg[1] = ext.InlineKeyboardButton{Text: "✖", CallbackData: "close"}
		kn = append(kn, kg)

		return teks, a, kn
	}
	return "", nil, nil
}

func mainMenu(chatId int) (string, [][]string, [][]ext.InlineKeyboardButton) {
	a := extraction.GetEmoji(chatId)
	if a != nil {
		teks := function.GetString(chatId, "modules/helpers/function.go:85")

		var kn = make([][]ext.InlineKeyboardButton, 0)

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

		return teks, a, kn
	}
	return "", nil, nil
}

func LoadSettingPanel(u *gotgbot.Updater) {
	defer logrus.Info("Setting Panel Module Loaded...")
	go u.Dispatcher.AddHandler(handlers.NewPrefixCommand("settings", []rune{'/', '.'}, panel))
	go u.Dispatcher.AddHandler(handlers.NewCallback(
		"^m[cdefgb]_(toggle|warn|kick|ban|mute|reset|plus|minus|duration|waktu|del|warn)",
		usercontrolquery))
	go u.Dispatcher.AddHandler(handlers.NewCallback("mo_toggle", spamcontrolquery))
	go u.Dispatcher.AddHandler(handlers.NewCallback("mk_", settingquery))
	go u.Dispatcher.AddHandler(handlers.NewCallback("close", closequery))
	go u.Dispatcher.AddHandler(handlers.NewCallback("back_", backquery))
}
