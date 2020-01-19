package caching

import (
	"github.com/go-redis/redis"
	"time"

	"github.com/jumatberkah/antispambot/bot"
)

var REDIS *redis.Client = nil

func InitRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:         bot.BotConfig.RedisAddress,
		Password:     bot.BotConfig.RedisPassword,
		DB:           0,
		DialTimeout:  time.Second,
		MinIdleConns: 0,
	})

	REDIS = client
}
