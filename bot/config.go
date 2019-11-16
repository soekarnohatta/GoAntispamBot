package bot

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

// Config -> Config values
type Config struct {
	ApiKey        string
	OwnerId       int
	LogEvent      int
	LogBan        int
	MainGrp       int
	SudoUsers     []string
	SqlUri        string
	WebhookUrl    string
	WebhookPath   string
	WebhookServe  string
	WebhookPort   int
	RedisAddress  string
	RedisPassword string
	CleanPolling  string
}

// BotConfig -> Return config
var BotConfig Config

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error can't load .env file!")
	}
	returnConfig := Config{}
	var ok bool

	returnConfig.ApiKey, ok = os.LookupEnv("BOT_API_KEY")
	if !ok {
		log.Fatal("Missing API Key")
	}
	returnConfig.OwnerId, err = strconv.Atoi(os.Getenv("OWNER_ID"))
	if err != nil {
		log.Fatal("Missing Owner ID")
	}
	returnConfig.LogEvent, err = strconv.Atoi(os.Getenv("LOG_EVENT"))
	if err != nil {
		log.Fatal("Missing Log Group ID")
	}
	returnConfig.MainGrp, err = strconv.Atoi(os.Getenv("MAIN_GRP"))
	if err != nil {
		log.Fatal("Missing Main Group ID")
	}
	returnConfig.LogBan, err = strconv.Atoi(os.Getenv("LOG_BAN"))
	if err != nil {
		log.Fatal("Missing Ban Log Group ID")
	}
	returnConfig.SudoUsers = strings.Split(os.Getenv("SUDO_USERS"), " ")
	returnConfig.SqlUri, ok = os.LookupEnv("DATABASE_URI")
	if !ok {
		log.Fatal("Missing PostgreSQL URI")
	}
	returnConfig.WebhookUrl, ok = os.LookupEnv("WEBHOOK_URL")
	if !ok {
		returnConfig.WebhookUrl = ""
	}
	returnConfig.WebhookPath, ok = os.LookupEnv("WEBHOOK_PATH")
	if !ok {
		returnConfig.WebhookPath = "api/bot"
	}
	returnConfig.WebhookServe, ok = os.LookupEnv("WEBHOOK_SERVE")
	if !ok {
		returnConfig.WebhookServe = "localhost"
	}
	returnConfig.WebhookPort, err = strconv.Atoi(os.Getenv("WEBHOOK_PORT"))
	if err != nil {
		returnConfig.WebhookPort = 5000
	}
	returnConfig.RedisAddress, ok = os.LookupEnv("REDIS_ADDRESS")
	if !ok {
		returnConfig.RedisAddress = "localhost:6379"
	}
	returnConfig.RedisPassword, ok = os.LookupEnv("REDIS_PASSWORD")
	if !ok {
		returnConfig.RedisPassword = ""
	}
	returnConfig.CleanPolling, ok = os.LookupEnv("CLEAN_POLLING")
	if !ok {
		returnConfig.RedisPassword = "false"
	}

	BotConfig = returnConfig
}
