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

var Prefix = []rune{'!', '/'}

func registerHandlers(u *gotgbot.Updater) {
	//Message handler(s)
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, handler.Handler{}.UpdateChat))
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, handler.Handler{}.GbanHandler))

	//Regex handler(s)
	u.Dispatcher.AddHandler(handlers.NewRegex("(^ping|/ping)", handler.Handler{}.PingHandler))

	//Command handler(s)
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("lang", Prefix, commands.CommandLang{}.SetLang))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("fban", Prefix, commands.CommandAdmin{}.BanUser))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("unfban", Prefix, commands.CommandAdmin{}.UnBanUser))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("dbg", Prefix, commands.CommandAdmin{}.Debug))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("start", Prefix, commands.CommandStart{}.Start))
}

func run() {
	updater, err := gotgbot.NewUpdater(bot.BotConfig.BotApiKey)
	errHandler.Fatal(err)
	registerHandlers(updater) // Register all defined handler(s)
	err = updater.StartCleanPolling()
	errHandler.Fatal(err)
}

func main() {
	run()
}
