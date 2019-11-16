package backups

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/jumatberkah/antispambot/bot"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/chat_status"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/err_handler"
	"github.com/keighl/barkup"
	"github.com/robfig/cron"
	"time"
)

func Backup(b ext.Bot, u *gotgbot.Update) {
	c := cron.New()
	_ = c.AddFunc("@hourly", func() {
		f := pg_backup()
		_ = f.To("dump.sql.tar.gz", nil)
		_, _ = b.SendDocumentCaptionStr(bot.BotConfig.MainGrp, f.Filename(), time.Now().String())
	})
	c.Start()
}

func backup(b ext.Bot, u *gotgbot.Update) error {
	if chat_status.IsOwner(u.EffectiveUser.Id) == false {
		return gotgbot.EndGroups{}
	}

	_, err := u.EffectiveMessage.ReplyText("Starting Backup...")
	err_handler.HandleTgErr(b, u, err)
	postgres := &barkup.Postgres{}

	// Writes a file `./bu_DBNAME_TIMESTAMP.sql.tar.gz`
	result := postgres.Export()
	err_handler.HandleErr(result.Error)
	_, err = b.SendDocumentCaptionStr(bot.BotConfig.MainGrp, result.Filename(), time.Now().String())
	err_handler.HandleTgErr(b, u, err)
	_, err = b.SendDocumentStr(bot.BotConfig.MainGrp, result.Filename())
	err_handler.HandleTgErr(b, u, err)
	return err
}

func pg_backup() *barkup.ExportResult {
	postgres := &barkup.Postgres{}

	// Writes a file `./bu_DBNAME_TIMESTAMP.sql.tar.gz`
	result := postgres.Export()
	err_handler.HandleErr(result.Error)
	return result
}

// LoadBackups -> Register Handler
func LoadBackups(u *gotgbot.Updater) {
	go u.Dispatcher.AddHandler(handlers.NewPrefixCommand("backup", []rune{'/', '.'}, backup))
}
