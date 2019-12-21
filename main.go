package main

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/jumatberkah/antispambot/bot"
	"github.com/jumatberkah/antispambot/bot/modules/admins"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/caching"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
	"github.com/jumatberkah/antispambot/bot/modules/info"
	"github.com/jumatberkah/antispambot/bot/modules/language"
	"github.com/jumatberkah/antispambot/bot/modules/listener"
	"github.com/jumatberkah/antispambot/bot/modules/private"
	"github.com/jumatberkah/antispambot/bot/modules/setting"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("Starting Bot...")

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
	info.LoadInfo(updater)
	listener.LoadListeners(updater)

	if bot.BotConfig.WebhookUrl != "" {
		logrus.Info("Using Webhook...")
		webHook := gotgbot.Webhook{
			URL:            bot.BotConfig.WebhookUrl,
			MaxConnections: 20,
			Serve:          bot.BotConfig.WebhookServe,
			ServePort:      bot.BotConfig.WebhookPort,
			ServePath:      bot.BotConfig.WebhookPath,
		}
		updater.StartWebhook(webHook)
		_, err = updater.SetWebhook(webHook.ServePath, webHook)
		err_handler.HandleErr(err)
	} else if bot.BotConfig.CleanPolling == "true" {
		logrus.Info("Using Clean Polling...")
		_ = updater.StartCleanPolling()
	} else {
		logrus.Info("Using Long Polling...")
		_ = updater.StartPolling()
	}
	updater.Idle()
}
