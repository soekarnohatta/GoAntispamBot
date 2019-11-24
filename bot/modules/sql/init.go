package sql

import (
	"github.com/jinzhu/gorm"
	"github.com/jumatberkah/antispambot/bot"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/err_handler"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var SESSION *gorm.DB

func InitDb() {
	conn, err := pq.ParseURL(bot.BotConfig.SqlUri)
	err_handler.FatalError(err)

	db, err := gorm.Open("postgres", conn)
	err_handler.FatalError(err)
	SESSION = db

	db.AutoMigrate(&User{}, &Chat{}, &UserSpam{}, &ChatSpam{}, &Setting{}, &Verify{}, &Picture{}, &Username{},
		&EnforceGban{}, &Lang{}, &Warns{}, &WarnSettings{}, &Notification{}, &Antispam{}, &NewUser{})
	logrus.Info("Database has been connected & Auto-migrated database schema")
}
