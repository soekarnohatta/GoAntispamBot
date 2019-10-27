package main

import (
	"github.com/PaulSonOfLars/goloc"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/jumatberkah/antispambot/bot"
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

func loadlang() {
	files, err := ioutil.ReadDir("trans")
	err_handler.FatalError(err)
	for _, f := range files {
		if f.IsDir() {
			goloc.LoadAll(f.Name())
		}
	}
}

func main() {
	// language
	loadlang()

	// connect DB
	sql.InitDb()

	// initiation
	updater, err := gotgbot.NewUpdater(bot.BotConfig.ApiKey)
	err_handler.FatalError(err)

	// registering handlers
	modules.LoadAdmins(updater)
	modules.LoadLang(updater)
	modules.LoadSetting(updater)
	modules.LoadSettingPanel(updater)
	modules.LoadPm(updater)
	modules.LoadListeners(updater)

	// start clean polling / webhook
	if bot.BotConfig.WebhookUrl != "" {
		logrus.Warn("Using Webhook...")
		var web gotgbot.Webhook
		web.URL = bot.BotConfig.WebhookUrl
		web.MaxConnections = 40
		web.Serve = "localhost"
		web.ServePort = bot.BotConfig.WebhookPort
		_, err = updater.SetWebhook(bot.BotConfig.WebhookPath, web)
		err_handler.HandleErr(err)
		updater.StartWebhook(web)
	} else if bot.BotConfig.CleanPolling == "true" {
		logrus.Warn("Using Long Clean Polling...")
		_ = updater.StartCleanPolling()
	} else {
		logrus.Warn("Using Long Polling...")
		_ = updater.StartPolling()
	}

	// wait
	updater.Idle()
}
