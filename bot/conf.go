// Package bot contains bot configurations to be used by other funcs.
package bot

import (
	"GoAntispamBot/bot/helpers/errHandler"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

// BConfig is a set of bot's configurations.
type BConfig struct {
	// Bot configs
	BotAPIKey    string `env:"BOT_API_KEY,required"`
	CleanPolling bool   `env:"CLEAN_POLLING,required"`
	WebhookURL   string `env:"WEBHOOK_URL,required"`
	WebhookPath  string `env:"WEBHOOK_PATH,required"`
	WebhookServe string `env:"WEBHOOK_SERVE,required"`
	WebhookPort  int    `env:"WEBHOOK_PORT,required"`

	// User configs
	OwnerID   int   `env:"OWNER_ID,required"`
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
	err = env.Parse(returnConfig)                           // Parse the env vars.
	errHandler.Fatal(err)                                   // Creates fatal error when env vars is empty.
	BotConfig = returnConfig                                // Assign filled struct to the variable.
}
