# GoAntispamBot

GoAntispamBot is a Telegram bot written in golang to protect your chat from spammers, intruders, disturbers, etc. This bot can take action(s) for those who haven't set usernames, profile pictures, or even sharing links and non sense arabic words. 

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on your system.

### Prerequisites

What things you need to install or you have to maintain or deploy this software. 

```
-Golang Version 1.12 and above.
-PostgreSQL Version 12 and above.
-Redis.
-Windows/Linux, Linux is recommended
-Telegram Bot Token
```

### Setting Up Enviroment

You have to set up some variables or configurations in order for this bot to work properly. Below are the `.example.env` that located in `/data/configurations` explanation.

 - `BOT_API_KEY `: Your bot token, get it at BotFather
 - `LOG_EVENT`: Your channel id to track all events happened
 - `LOG_BAN`: Your channel id to track all ban activities happened
 - `MAIN_GRP`: Your group id to track all events happened
 - `OWNER_ID`: Your id in telegram
 - `SUDO_USERS`: A space separated list of user_ids which should be considered sudo users
 - `DATABASE_URI`: Your database URL
 - `WEBHOOK_PORT`: Your webhook port
 - `WEBHOOK_URL`: Setting this to ANYTHING will disable webhooks
 - `WEBHOOK_PATH`: The path your webhook should connect to (only needed for webhook mode)
 - `WEBHOOK_SERVE`: A space separated list of user_ids which should be considered sudo users
 - `REDIS_ADDRESS` : Your Redis Address
 - `REDIS_PASSWORD` : Your Redis Password
 - `REDIS_DB` : Your selected Redis DB
 - `CLEAN_POLLING` : Setting this to `true` will enable clean polling
 - `BOT_VERSION` : Your bot version

The last thing is, rename your `.example.env` to `.env` and go ahead to compiling section.
### Compiling

A step by step series of examples that tell you how to compile the code
```
go build -mod vendor
```
Ends by running your bot then type `/start` in your bot.  

### Running the bot

Run your bot by typing `./antispambot` for linux, or `antispambot.exe` for windows in the command line. Open Telegram and type `/start` in your bot. If there is no respond, you have to stepback and correct what have you done wrong.

## License

This project is licensed under the GNU GPL V3 License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

* Hat tip to anyone whose code was used.
* Telegram Community.