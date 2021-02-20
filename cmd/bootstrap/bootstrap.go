package bootstrap

import (
	"github.com/gempir/go-twitch-irc/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"sergiocarracedo.es/streambot-go/internal/bot"
	chatLogger "sergiocarracedo.es/streambot-go/internal/chat-logger"
	"sergiocarracedo.es/streambot-go/internal/utils"
	"strconv"
)

func Run() error {
	// Database
	log.Printf("Inicializating chat database...")
	db, err := sqlx.Open("sqlite3", "data/db.sqlite")
	if err != nil {
		log.Printf(utils.Colors.Red, err.Error(), utils.Colors.Reset)
		return err
	}

	// Twitch client
	client := twitch.NewClient(
		os.Getenv("BOT_USERNAME"),
		os.Getenv("TMI_OAUTH_TOKEN"))

	// Bot
	disableBot, _ := utils.GetEnvBool("DISABLE_BOT")
	var botService bot.Service
	if !disableBot {
		disabledCommands, _ := utils.GetEnvArrayStr("DISABLED_COMMANDS")
		coldDownTime, _ := utils.GetEnvInt64("COMMAND_COLD_DOWN_TIME")
		sendHelpMessageEvery, _ := strconv.ParseInt(os.Getenv("COMMAND_HELP_EVERY"), 10, 64)

		botService = bot.New(
			client,
			os.Getenv("CHANNEL"),
			disabledCommands,
			coldDownTime,
			sendHelpMessageEvery,
		)

		log.Println("Stating bot...")
		log.Println("Getting messages from " + os.Getenv("CHANNEL"))

		log.Printf("Available commands: %v", botService.GetCommandListNames())
		log.Printf("Disabled commands: %v", disabledCommands)
	}

	// Chat Logger
	var chatLoggerService *chatLogger.Service
	disableChatLogger, _ := utils.GetEnvBool("DISABLE_CHAT_LOGGER")
	if !disableChatLogger {
		chatLoggerService, err = chatLogger.New(
			client,
			db)
		if err != nil {
			log.Printf(utils.Colors.Red, err.Error(), utils.Colors.Reset)
			return err
		}
	}

	// On Private Message
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		log.Println(utils.Colors.Blue, "Private Message", message.Channel, message.User.Name, message.Message, utils.Colors.Reset)

		if !disableBot {
			published := botService.OnPrivateMessage(message)
			if !published {
				log.Println(utils.Colors.Gray, "Cold down", message.User.Name, message.Message, utils.Colors.Reset)
			}
		}

		if !disableChatLogger {
			err := chatLoggerService.OnPrivateMessage(message)
			if err != nil {
				log.Printf(utils.Colors.Red, err.Error(), utils.Colors.Reset)
			}
		}
	})

	// On User Notice Message
	client.OnUserNoticeMessage(func(message twitch.UserNoticeMessage) {
		log.Println(utils.Colors.Green, "User Notice Message", message.Channel, message.User, message.Message, utils.Colors.Reset)
		if !disableChatLogger {
			err := chatLoggerService.OnUserNoticeMessage(message)
			if err != nil {
				log.Printf(utils.Colors.Red, err.Error(), utils.Colors.Reset)
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
