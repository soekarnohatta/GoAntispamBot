package main

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"

	"github.com/jumatberkah/antispambot/bot"
	"github.com/jumatberkah/antispambot/bot/handlers/commands/admins"
	"github.com/jumatberkah/antispambot/bot/handlers/commands/user/help"
	"github.com/jumatberkah/antispambot/bot/handlers/commands/user/info"
	"github.com/jumatberkah/antispambot/bot/handlers/commands/user/private"
	"github.com/jumatberkah/antispambot/bot/handlers/commands/user/reporting"
	"github.com/jumatberkah/antispambot/bot/handlers/commands/user/setting"
	"github.com/jumatberkah/antispambot/bot/handlers/listener"
	"github.com/jumatberkah/antispambot/bot/helpers/caching"
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/helpers/function"
	"github.com/jumatberkah/antispambot/bot/sql"
)

func multiInstance() {
	for _, botToken := range bot.BotConfig.ApiKey {
		logrus.Info("Starting Bot...")

		updater, err := gotgbot.NewUpdater(botToken)
		err_handler.FatalError(err)
		registerHandlers(updater)

		if bot.BotConfig.WebhookUrl != "" {
			for _, hookPath := range bot.BotConfig.WebhookPath {
				webHook := gotgbot.Webhook{
					URL:            bot.BotConfig.WebhookUrl,
					MaxConnections: 20,
					Serve:          bot.BotConfig.WebhookServe,
					ServePort:      bot.BotConfig.WebhookPort,
					ServePath:      hookPath,
				}
				updater.StartWebhook(webHook)
				_, err = updater.SetWebhook(webHook.ServePath, webHook)
				err_handler.HandleErr(err)
				logrus.Info("Using Webhook...")
			}
		} else if bot.BotConfig.CleanPolling == "true" {
			logrus.Info("Using Clean Polling...")
			_ = updater.StartCleanPolling()
		} else {
			logrus.Info("Using Long Polling...")
			_ = updater.StartPolling()
		}

		if bot.BotConfig.ApiKey[len(bot.BotConfig.ApiKey)-1] == botToken {
			logrus.Info(fmt.Sprintf("Bot Running On Version: %v - %v", bot.BotConfig.BotVer,
				updater.Bot.UserName))
			updater.Idle()
		} else {
			logrus.Info(fmt.Sprintf("Bot Running On Version: %v - %v", bot.BotConfig.BotVer,
				updater.Bot.UserName))
			go updater.Idle()
			go func() {
				// cron job
				c := cron.New()
				_ = c.AddFunc("@every 30m", func() {
					_ = admins.CronBackupDb(updater)
				})
				c.Start()
			}()
		}
	}
}

func singleInstance() {
	logrus.Info("Starting Bot...")

	updater, err := gotgbot.NewUpdater(bot.BotConfig.ApiKey[0])
	err_handler.FatalError(err)
	registerHandlers(updater)

	if bot.BotConfig.WebhookUrl != "" {
		webHook := gotgbot.Webhook{
			URL:            bot.BotConfig.WebhookUrl,
			MaxConnections: 20,
			Serve:          bot.BotConfig.WebhookServe,
			ServePort:      bot.BotConfig.WebhookPort,
			ServePath:      bot.BotConfig.WebhookPath[0],
		}
		updater.StartWebhook(webHook)
		_, err = updater.SetWebhook(webHook.ServePath, webHook)
		err_handler.HandleErr(err)
		logrus.Info("Using Webhook...")

	} else if bot.BotConfig.CleanPolling == "true" {
		logrus.Info("Using Clean Polling...")
		_ = updater.StartCleanPolling()
	} else {
		logrus.Info("Using Long Polling...")
		_ = updater.StartPolling()
	}

	logrus.Info(fmt.Sprintf("Bot Running On Version: %v - %v", bot.BotConfig.BotVer,
		updater.Bot.UserName))
	updater.Idle()
}

func registerHandlers(updater *gotgbot.Updater) {
	admins.LoadAdmins(updater)
	setting.LoadSetting(updater)
	private.LoadPm(updater)
	info.LoadInfo(updater)
	help.LoadHelp(updater)
	reporting.LoadReport(updater)

	listener.LoadSettingListener(updater)
	listener.LoadHelpListener(updater)
	listener.LoadStartListener(updater)
	listener.LoadReportListener(updater)
	listener.LoadLangListener(updater)
	listener.LoadUserListener(updater)
}

func main() {
	sql.InitDb()
	caching.InitRedis()
	caching.InitCache()
	function.LoadAllLang()

	multiInstance() // This is used if you want multiple bot running on single instance. Be aware that this can take much resources.
	//singleInstance() // This is used if you have only single instance. Do not use both of them!
}
