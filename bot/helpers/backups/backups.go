package backups

import (
	"github.com/keighl/barkup"
	"github.com/robfig/cron/v3"
)

func Backup() {
	c := cron.New()
	_, _ = c.AddFunc("@hourly", pg_backup)
}

func pg_backup() {
	postgres := &barkup.Postgres{
		Host: "127.0.0.1",
		Port: "5432",

		// Not necessary if the program runs as an authorized pg user/role
		Username: "postgres",

		// Any extra pg_dump options
		Options: []string{"--no-owner"},
	}

	// Writes a file `./bu_DBNAME_TIMESTAMP.sql.tar.gz`
	result := postgres.Export()
	_ = result.To("backups", nil)

	if result.Error != nil {
		panic(result.Error)
	}
}
