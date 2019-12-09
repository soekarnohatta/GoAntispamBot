package caching

import (
	"github.com/allegro/bigcache"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/err_handler"
	"time"
)

var CACHE *bigcache.BigCache

func InitCache() {
	config := bigcache.Config{Shards: 1024,
		LifeWindow:         2 * time.Hour,
		CleanWindow:        5 * time.Minute,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       500,
		HardMaxCacheSize:   512,
		OnRemove:           nil,
		OnRemoveWithReason: nil,
	}

	cache, err := bigcache.NewBigCache(config)
	err_handler.HandleErr(err)
	CACHE = cache
}
