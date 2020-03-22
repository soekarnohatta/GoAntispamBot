/*
Package "services" is a package that provides services to be used by other funcs.
This package should has all services for the bot.
*/
package langService

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"

	"GoAntispamBot/bot/model"
	"GoAntispamBot/bot/providers/mongoProvider"
	"GoAntispamBot/bot/providers/redisProvider"
)

func UpdateLang(chatID int, language string) {
	// Initiate & fills the empty struct
	var langStruct = model.Lang{}
	langStruct.ChatID = chatID
	langStruct.Language = language

	// Start updating...
	go mongoProvider.Update("lang", langStruct.ChatID, bson.M{"$set": langStruct}, true)
	go updateLangRedis(chatID, language)
}

func RemoveLang(chatID int) {
	// Initiate & fills the empty struct
	var langStruct = model.Lang{}
	langStruct.ChatID = chatID

	// Start removing...
	go mongoProvider.Remove("lang", langStruct.ChatID)
}

func FindLang(chatID int) (lang string) {
	// Initiate & fills the empty struct
	var langStruct = model.Lang{}
	langStruct.ChatID = chatID

	// Start search...
	langcache := findLangRedis(chatID) // Redis
	if langcache != "" {
		lang = langcache
		return
	}

	res := mongoProvider.FindOne("lang", langStruct.ChatID) // MongoDB
	if res != nil {
		lang = string(res.Lookup("language").Value)
		return
	}

	lang = "en-GB"
	return
}

// --------------------------------------------------------------------------------------

func updateLangRedis(chatID int, language string) {
	// Start updating...
	key := fmt.Sprintf("lang_%v", chatID)
	redisProvider.SetRedisKey(key, language)
}

func findLangRedis(chatID int) (lang string) {
	// Start search...
	key := fmt.Sprintf("lang_%v", chatID)
	lang = redisProvider.GetRedisKey(key)
	return
}
