package bootstrap

import (
	"fmt"
	"github.com/gempir/go-twitch-irc/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	env "sergiocarracedo.es/streambot-go/internal"
	"sergiocarracedo.es/streambot-go/internal/bot"
	chatLogger "sergiocarracedo.es/streambot-go/internal/chat-logger"
	"strconv"
)

func Run() error {
	// Database
	log.Printf("Inicializating chat database...")
	db, err := sqlx.Open("sqlite3", "data/db.sqlite")
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	// Twitch client
	client := twitch.NewClient(
		os.Getenv("BOT_USERNAME"),
		os.Getenv("TMI_OAUTH_TOKEN"))

	// Bot
	disableBot, _ := env.GetEnvBool("DISABLE_BOT")
	var botService bot.Service
	if !disableBot {
		disabledCommands, _ := env.GetEnvArrayStr("DISABLED_COMMANDS")
		coldDownTime, _ := env.GetEnvInt64("COMMAND_COLD_DOWN_TIME")
		sendHelpMessageEvery, _ := strconv.ParseInt(os.Getenv("COMMAND_HELP_EVERY"), 10, 64)

		botService = bot.New(
			client,
			os.Getenv("CHANNEL"),
			disabledCommands,
			coldDownTime,
			sendHelpMessageEvery,
			)

		fmt.Println("Stating bot...")
		log.Println("Getting messages from " + os.Getenv("CHANNEL"))

		log.Printf("Available commands: %v", botService.GetCommandListNames())
		log.Printf("Disabled commands: %v", disabledCommands)
	}



	// Logger
	var chatLoggerService *chatLogger.Service
	disableChatLogger, _ := env.GetEnvBool("DISABLE_CHAT_LOGGER")
	if !disableChatLogger {
		chatLoggerService, err = chatLogger.New(
			client,
			db)
		if err != nil {
			log.Printf(err.Error())
			return err
		}
	}


	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Println("Private Message", message.Channel, message.User, message.Message)

		if !disableBot {
			published := botService.OnPrivateMessage(message)
			if !published {
				log.Println("Cold down", message.User, message.Message)
			}
		}

		if !disableChatLogger {
			err := chatLoggerService.OnPrivateMessage(message)
			if err != nil {
				log.Printf(err.Error())
			}
		}
	})


	client.OnUserNoticeMessage(func(message twitch.UserNoticeMessage) {
		log.Println("User Notice Message", message.Channel, message.User, message.Message)
		if !disableChatLogger {
			err := chatLoggerService.OnUserNoticeMessage(message)
			if err != nil {
				log.Printf(err.Error())
			}
		}
	})




	// On connect
	client.OnConnect(func() {
		if !disableBot {
			botService.OnConnect()
		}

	})

	client.Join(os.Getenv("CHANNEL"))

	return client.Connect()
}