package bot

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ApiKey        string
	OwnerId       int
	LogEvent      int
	LogBan        int
	SudoUsers     []string
	SqlUri        string
	WebhookUrl    string
	WebhookPath   string
	WebhookPort   int
	RedisAddress  string
	RedisPassword string
	CleanPolling  string
}

var BotConfig Config

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error can't load .env file!")
	}
	returnConfig := Config{}

	// Assign config struct values by loading them from the env
	var ok bool

	returnConfig.ApiKey, ok = os.LookupEnv("BOT_API_KEY") // If env var is empty
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
	returnConfig.LogBan, err = strconv.Atoi(os.Getenv("LOG_BAN"))
	if err != nil {
		log.Fatal("Missing Ban Log Group ID")
	}
	returnConfig.SudoUsers = strings.Split(os.Getenv("SUDO_USERS"), " ")
	returnConfig.SqlUri, ok = os.LookupEnv("DATABASE_URI")
	// If env var is empty
	if !ok {
		log.Fatal("Missing PostgreSQL URI")
	}
	returnConfig.WebhookUrl, ok = os.LookupEnv("WEBHOOK_URL")
	// If env var is empty
	if !ok {
		returnConfig.WebhookUrl = ""
	}
	returnConfig.WebhookPath, ok = os.LookupEnv("WEBHOOK_PATH")
	// If env var is empty
	if !ok {
		returnConfig.WebhookPath = "bot"
	}
	returnConfig.WebhookPort, err = strconv.Atoi(os.Getenv("WEBHOOK_PORT"))
	// If env var is empty
	if err != nil {
		returnConfig.WebhookPort = 5000
	}
	returnConfig.RedisAddress, ok = os.LookupEnv("REDIS_ADDRESS")
	// If env var is empty
	if !ok {
		returnConfig.RedisAddress = "localhost:6379"
	}
	returnConfig.RedisPassword, ok = os.LookupEnv("REDIS_PASSWORD")
	// If env var is empty
	if !ok {
		returnConfig.RedisPassword = ""
	}
	returnConfig.CleanPolling, ok = os.LookupEnv("CLEAN_POLLING")
	// If env var is empty
	if !ok {
		returnConfig.RedisPassword = "false"
	}

	BotConfig = returnConfig
}
