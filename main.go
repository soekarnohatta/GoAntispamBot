package main

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/jumatberkah/antispambot/bot"
	"github.com/jumatberkah/antispambot/bot/helpers/backups"
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
	"github.com/sirupsen/logrus"
)

func main() {
	sql.InitDb()

	updater, err := gotgbot.NewUpdater(bot.BotConfig.ApiKey)
	err_handler.FatalError(err)

	go modules.LoadLang(updater)
	go modules.LoadAdmins(updater)
	go modules.LoadSetting(updater)
	go modules.LoadSettingPanel(updater)
	go modules.LoadPm(updater)
	go modules.LoadListeners(updater)

	if bot.BotConfig.WebhookUrl != "" {
		logrus.Warn("Using Webhook...")
		var web gotgbot.Webhook
		web.URL = bot.BotConfig.WebhookUrl
		web.MaxConnections = 20
		web.ServePort = bot.BotConfig.WebhookPort
		web.Serve = bot.BotConfig.WebhookServe
		web.ServePath = bot.BotConfig.WebhookPath
		updater.StartWebhook(web)
		_, err = updater.SetWebhook(bot.BotConfig.WebhookPath, web)
		err_handler.HandleErr(err)
	} else if bot.BotConfig.CleanPolling == "true" {
		logrus.Warn("Using Clean Polling...")
		_ = updater.StartCleanPolling()
	} else {
		logrus.Warn("Using Long Polling...")
		_ = updater.StartPolling()
	}

	updater.Idle()
	backups.Backup()
}
