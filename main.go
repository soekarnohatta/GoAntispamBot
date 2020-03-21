package main

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/handlers/Filters"

	"GoAntispamBot/bot"
	"GoAntispamBot/bot/helpers/errHandler"
	"GoAntispamBot/bot/model"
)

var Prefix = []rune{'!', '/'}

func registerHandlers(u *gotgbot.Updater) {
	//Message handler(s)
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, model.Message{}.UpdateChat))
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, model.Message{}.GbanHandler))

	//Regex handler(s)
	u.Dispatcher.AddHandler(handlers.NewRegex("(^ping|/ping)", model.Message{}.GbanHandler))

	//Command handler(s)
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("lang", Prefix, model.Command{}.SetLang))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("fban", Prefix, model.Command{}.BanUser))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("unfban", Prefix, model.Command{}.UnBanUser))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("dbg", Prefix, model.Command{}.Debug))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("start", Prefix, model.Command{}.Start))
}

func main() {
	updater, err := gotgbot.NewUpdater(bot.BotConfig.BotApiKey)
	errHandler.Fatal(err)
	registerHandlers(updater) // Register all defined handler(s)
	err = updater.StartCleanPolling()
	errHandler.Fatal(err)
}
