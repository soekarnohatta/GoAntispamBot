package admins

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/jumatberkah/antispambot/bot"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/chat_status"
	"github.com/jumatberkah/antispambot/bot/modules/helpers/err_handler"
)

type PostgreSQL struct {
	host        string
	port        string
	database    string
	username    string
	password    string
	dumpCommand string
}

func backupDb(b ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage

	if !chat_status.RequireOwner(msg, msg.From.Id) {
		return nil
	}
	ctx := PostgreSQL{
		host:     "",
		port:     "5432",
		username: "",
		password: "",
		database: "",
	}

	if err := ctx.prepare(); err != nil {
		return err
	}

	err := ctx.dump()
	err_handler.HandleErr(err)

	f, _ := os.Open("data/backup/nayef.sql")
	defer f.Close()
	sendDoc := b.NewSendableDocument(msg.Chat.Id, "BACKUP_"+time.Now().String())
	sendDoc.Name = "BACKUP_" + fmt.Sprint(time.Now().Unix())
	sendDoc.Reader = f
	_, err = sendDoc.Send()
	err_handler.HandleErr(err)
	return err
}

func CronBackupDb(b *gotgbot.Updater) error {
	ctx := PostgreSQL{
		host:     "",
		port:     "5432",
		username: "",
		password: "",
		database: "",
	}

	if err := ctx.prepare(); err != nil {
		return err
	}

	err := ctx.dump()
	err_handler.HandleErr(err)

	f, _ := os.Open("data/backup/nayef.sql")
	defer f.Close()
	sendDoc := b.Bot.NewSendableDocument(bot.BotConfig.OwnerId, "BACKUP_"+time.Now().String())
	sendDoc.Name = "BACKUP_" + fmt.Sprint(time.Now().Unix())
	sendDoc.Reader = f
	_, err = sendDoc.Send()
	err_handler.HandleErr(err)
	return err
}

func (ctx *PostgreSQL) prepare() (err error) {
	// mysqldump command
	dumpArgs := []string{}
	if len(ctx.database) == 0 {
		return fmt.Errorf("PostgreSQL database config is required")
	}
	if len(ctx.host) > 0 {
		dumpArgs = append(dumpArgs, "--host="+ctx.host)
	}
	if len(ctx.port) > 0 {
		dumpArgs = append(dumpArgs, "--port="+ctx.port)
	}
	if len(ctx.username) > 0 {
		dumpArgs = append(dumpArgs, "--username="+ctx.username)
	}

	ctx.dumpCommand = "pg_dump " + strings.Join(dumpArgs, " ") + " " + ctx.database
	return nil
}

func (ctx *PostgreSQL) dump() error {
	dumpFilePath := path.Join("data/backup/" + ctx.database + ".sql")
	logrus.Info("-> Dumping PostgreSQL...")
	if len(ctx.password) > 0 {
		_ = os.Setenv("PGPASSWORD", ctx.password)
	}
	_, err := Exec(ctx.dumpCommand, "-f", dumpFilePath)
	if err != nil {
		return err
	}
	logrus.Info("dump path:", dumpFilePath)
	return nil
}

func Exec(command string, args ...string) (output string, err error) {
	spaceRegexp := regexp.MustCompile("[\\s]+")
	commands := spaceRegexp.Split(command, -1)
	command = commands[0]
	commandArgs := []string{}
	if len(commands) > 1 {
		commandArgs = commands[1:]
	}
	if len(args) > 0 {
		commandArgs = append(commandArgs, args...)
	}

	fullCommand, err := exec.LookPath(command)
	if err != nil {
		return "", fmt.Errorf("%s cannot be found", command)
	}

	cmd := exec.Command(fullCommand, commandArgs...)
	cmd.Env = os.Environ()

	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr

	// logger.Debug(fullCommand, " ", strings.Join(commandArgs, " "))

	out, err := cmd.Output()
	if err != nil {
		logrus.Debug(fullCommand, " ", strings.Join(commandArgs, " "))
		err = errors.New(stdErr.String())
		return
	}

	output = strings.Trim(string(out), "\n")
	return output, err
}
