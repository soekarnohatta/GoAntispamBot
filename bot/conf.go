/*
Package "bot" is a package that defines all bot actions made.
This package should has all bot action(s) to an event.
*/
package bot

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"

	"GoAntispamBot/bot/helpers/errHandler"
)

// Struct "BConfig" is a set of bot's configurations.
type BConfig struct {
	// Bot configs
	BotApiKey    string `env:"BOT_API_KEY,required"`
	CleanPolling bool   `env:"CLEAN_POLLING,required"`
	WebhookUrl   string `env:"WEBHOOK_URL,required"`
	WebhookPath  string `env:"WEBHOOK_PATH,required"`
	WebhookServe string `env:"WEBHOOK_SERVE,required"`
	WebhookPort  int    `env:"WEBHOOK_PORT,required"`

	// User configs
	OwnerId   int   `env:"OWNER_ID,required"`
	LogEvent  int   `env:"LOG_EVENT,required"`
	LogBan    int   `env:"LOG_BAN,required"`
	MainGrp   int   `env:"MAIN_GRP,required"`
	SudoUsers []int `env:"SUDO_USERS,required" envSeparator:":"`

	// Misc
	DatabaseURL   string `env:"DATABASE_URL,required"`
	RedisAddress  string `env:"REDIS_ADDRESS,required"`
	RedisPassword string `env:"REDIS_PASSWORD,required"`
	BotVer        string `env:"BOT_VERSION,required"`
}

// BotConfig will returns "Config" struct
var BotConfig BConfig

// Init function will execute some code to fill the BConfig struct respectively.
func init() {
	returnConfig := BConfig{}                               // Initiate an empty BConfig struct.
	err := godotenv.Load("bot/storage/configurations/.env") // Load env vars from a file
	errHandler.Fatal(err)                                   // Creates fatal error when env file is not found.

	err = env.Parse(returnConfig) // Parse the env vars.
	errHandler.Fatal(err)         // Creates fatal error when env vars is empty.

	BotConfig = returnConfig // Assign filled struct to the variable.
}
