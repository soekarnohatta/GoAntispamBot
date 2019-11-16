package function

import (
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/extraction"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
	"strconv"
)

func MainControlMenu(chatId int) (string, [][]string, [][]ext.InlineKeyboardButton) {
	a := extraction.GetEmoji(chatId)
	teks := GetStringf(chatId, "modules/helpers/function.go:13",
		map[string]string{"1": a[0][0], "2": a[1][0], "3": a[2][0], "4": a[0][1], "5": a[1][1], "6": a[2][1], "7": a[0][2],
			"8": a[2][3], "9": a[3][0], "10": strconv.Itoa(sql.GetWarnSetting(strconv.Itoa(chatId)))})

	// Create Button(s)
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

func MainSpamMenu(chatId int) (string, [][]string, [][]ext.InlineKeyboardButton) {
	a := extraction.GetEmoji(chatId)
	teks := GetStringf(chatId, "modules/helpers/function.go:66", map[string]string{"1": a[0][3]})

	// Create Button(s)
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

func MainMenu(chatId int) (string, [][]string, [][]ext.InlineKeyboardButton) {
	a := extraction.GetEmoji(chatId)
	teks := GetString(chatId, "modules/helpers/function.go:85")

	// Create Button(s)
	var kn = make([][]ext.InlineKeyboardButton, 0)

	ki := make([]ext.InlineKeyboardButton, 2)
	ki[0] = ext.InlineKeyboardButton{Text: GetString(chatId, "modules/helpers/function.go:91"), CallbackData: "mk_utama"}
	ki[1] = ext.InlineKeyboardButton{Text: GetString(chatId, "modules/helpers/function.go:92"), CallbackData: "mk_spam"}
	kn = append(kn, ki)

	kz := make([]ext.InlineKeyboardButton, 2)
	kz[0] = ext.InlineKeyboardButton{Text: GetString(chatId, "modules/helpers/function.go:96"), CallbackData: "mk_media"}
	kz[1] = ext.InlineKeyboardButton{Text: GetString(chatId, "modules/helpers/function.go:97"), CallbackData: "mk_pesan"}
	kn = append(kn, kz)

	kd := make([]ext.InlineKeyboardButton, 1)
	kd[0] = ext.InlineKeyboardButton{Text: GetString(chatId, "modules/helpers/function.go:101"), CallbackData: "mk_reset"}
	kn = append(kn, kd)

	kk := make([]ext.InlineKeyboardButton, 1)
	kk[0] = ext.InlineKeyboardButton{Text: GetString(chatId, "modules/helpers/function.go:105"), CallbackData: "close"}
	kn = append(kn, kk)

	return teks, a, kn
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
