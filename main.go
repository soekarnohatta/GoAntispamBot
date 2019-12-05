package main

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/jumatberkah/antispambot/bot"
	"github.com/jumatberkah/antispambot/bot/modules/admins"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/caching"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
	"github.com/jumatberkah/antispambot/bot/modules/language"
	"github.com/jumatberkah/antispambot/bot/modules/listener"
	"github.com/jumatberkah/antispambot/bot/modules/private"
	"github.com/jumatberkah/antispambot/bot/modules/setting"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
	"github.com/sirupsen/logrus"
)

func main() {
	updater, err := gotgbot.NewUpdater(bot.BotConfig.ApiKey)
	err_handler.FatalError(err)

	function.LoadAllLang()
	caching.InitRedis()
	caching.InitCache()
	sql.InitDb()

	language.LoadLang(updater)
	admins.LoadAdmins(updater)
	setting.LoadSetting(updater)
	setting.LoadSettingPanel(updater)
	private.LoadPm(updater)
	listener.LoadListeners(updater)

	if bot.BotConfig.WebhookUrl != "" {
		logrus.Warn("Using Webhook...")
		var web gotgbot.Webhook
		web.URL = bot.BotConfig.WebhookUrl
		web.MaxConnections = 30
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
}
