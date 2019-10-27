package modules

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
)

func start(_ ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat

	txtStart := "Halo, aku adalah bot untuk membisukan pengguna yang belum pasang username dan memblokir spammer serta" +
		" dilengkapi dengan berbagai fitur keamanan. Berikut beberapa perintah yang bisa kamu gunakan:\n\n" +
		"	1. /username (true/false) - Untuk mematikan atau menyalakan fitur mute user jika belum memasang username.\n" +
		"	2. /verify (true/false) - User harus verifikasi ketika masuk grup.\n" +
		"	3. /time (int m/h/d) - Untuk menentukan tenggat waktu semua perintah.\n" +
		"	4. /profilepicture (true/false) - Untuk mematikan atau menyalakan fitur mute user jika belum pasang pp\n" +
		"	5. /setting - Untuk menyetel semua pengaturan\n\n" +
		"Jika butuh bantuan silahkan masuk ke @polybotsupport. <b>Mohon donasi untuk pengembangan lebih lanjut.</b>"

	if chat.Type == "supergroup" {
		_, err := msg.Delete()
		return err
	} else {
		_, err := msg.ReplyHTML(txtStart)
		return err
	}
}

func LoadPm(u *gotgbot.Updater) {
	u.Dispatcher.AddHandler(handlers.NewCommand("start", start))
}
