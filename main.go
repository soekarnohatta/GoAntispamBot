package main

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/jumatberkah/antispambot/bot"
	"github.com/jumatberkah/antispambot/bot/modules/commands/admins"
	"github.com/jumatberkah/antispambot/bot/modules/commands/user/help"
	"github.com/jumatberkah/antispambot/bot/modules/commands/user/info"
	"github.com/jumatberkah/antispambot/bot/modules/commands/user/private"
	"github.com/jumatberkah/antispambot/bot/modules/commands/user/setting"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/caching"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/function"
	"github.com/jumatberkah/antispambot/bot/modules/listener"
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

	admins.LoadAdmins(updater)
	setting.LoadSetting(updater)
	private.LoadPm(updater)
	info.LoadInfo(updater)
	help.LoadHelp(updater)

	listener.LoadUserListener(updater)
	listener.LoadSettingListener(updater)
	listener.LoadHelpListener(updater)
	listener.LoadStartListener(updater)

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
