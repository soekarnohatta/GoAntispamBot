package main

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/handlers/Filters"

	"GoAntispamBot/bot"
	"GoAntispamBot/bot/handler"
	"GoAntispamBot/bot/handler/commands"
	"GoAntispamBot/bot/helpers/errHandler"
)

var prefix = []rune{'!', '/'}

func registerHandlers(u *gotgbot.Updater) {
	//Message handler(s)
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, handler.Handler{}.UpdateChat))
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, handler.Handler{}.GbanHandler))

	//Regex handler(s)
	u.Dispatcher.AddHandler(handlers.NewRegex("(^ping|/ping)", handler.Handler{}.PingHandler))

	//Command handler(s)
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("lang", prefix, commands.CommandLang{}.SetLang))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("fban", prefix, commands.CommandAdmin{}.BanUser))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("unfban", prefix, commands.CommandAdmin{}.UnBanUser))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("dbg", prefix, commands.CommandAdmin{}.Debug))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("start", prefix, commands.CommandStart{}.Start))
}

func run() {
	updater, err := gotgbot.NewUpdater(bot.BotConfig.BotAPIKey, nil)
	errHandler.Fatal(err)
	registerHandlers(updater) // Register all defined handler(s)
	err = updater.StartCleanPolling()
	errHandler.Fatal(err)
}

func main() {
	run()
}
