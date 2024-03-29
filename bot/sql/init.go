package sql

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/jumatberkah/antispambot/bot"
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
)

var SESSION *gorm.DB = nil

func InitDb() {
	conn, err := pq.ParseURL(bot.BotConfig.SqlUri)
	err_handler.FatalError(err)

	db, err := gorm.Open("postgres", conn)
	err_handler.FatalError(err)
	SESSION = db

	db.AutoMigrate(&User{}, &Chat{}, &UserSpam{}, &ChatSpam{}, &Setting{}, &Verify{}, &Picture{}, &Username{},
		&EnforceGban{}, &Lang{}, &Warns{}, &WarnSettings{}, &Notification{}, &Antispam{})
	logrus.Info("Database has been connected & Auto-migrated database schema")
}
