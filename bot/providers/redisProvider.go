/*
Package "providers" is a package that provides required things reqired by the bot
to be used by other funcs.
This package should has all providers for the bot.
*/
package providers

import (
	"github.com/go-redis/redis"
	"time"

	"GoAntispamBot/bot"
	"GoAntispamBot/bot/helpers/errHandler"
)

var Redis *redis.Client = nil

func init() {
	Redis = redis.NewClient(
		&redis.Options{
			Addr:         bot.BotConfig.RedisAddress,
			Password:     bot.BotConfig.RedisPassword,
			DB:           0,
			DialTimeout:  time.Second,
			MinIdleConns: 0,
		},
	)
}

func GetRedisKey(key string) string {
	val := Redis.Get(key)
	if val.Err() != nil {
		if !(val.Err() == redis.Nil) {
			return val.Val()
		}
	}
	return ""
}

func SetRedisKey(key string, val string) {
	set := Redis.Set(key, val, 7200)
	errHandler.Error(set.Err())
	SaveRedis()
}

func RemoveRedisKey(key string) {
	rem := Redis.Del(key)
	errHandler.Error(rem.Err())
}

func SaveRedis() {
	save := Redis.BgSave()
	errHandler.Error(save.Err())
}
