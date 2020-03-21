package model

import (
	"GoAntispamBot/bot/handlers"
	"GoAntispamBot/bot/handlers/commands"
	"GoAntispamBot/bot/providers"
)

type (
	Command struct {
		TelegramProvider providers.TelegramProvider
		*commands.CommandHandler
	}

	Message struct {
		TelegramProvider providers.TelegramProvider
		*handlers.UpdateHandler
	}
)
