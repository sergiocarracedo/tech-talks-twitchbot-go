package bootstrap

import (
	"github.com/gempir/go-twitch-irc/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"sergiocarracedo.es/streambot-go/internal/bot"
	chatLogger "sergiocarracedo.es/streambot-go/internal/chat-logger"
	"sergiocarracedo.es/streambot-go/internal/notifications"
	_ "sergiocarracedo.es/streambot-go/internal/messages"
	"sergiocarracedo.es/streambot-go/internal/platform/storage"
	"sergiocarracedo.es/streambot-go/internal/server"
	"sergiocarracedo.es/streambot-go/internal/utils"
	"strconv"
)

const (
	host = "localhost"
	port = 4000
)

func Run() error {
	// Database
	log.Printf("Inicializating chat database...")
	db, err := sqlx.Open("sqlite3", "data/db.sqlite")
	if err != nil {
		log.Printf(utils.Colors.Red, err.Error(), utils.Colors.Reset)
		return err
	}

	// Web server
	server := server.New(host, port)

	// Twitch client
	client := twitch.NewClient(
		os.Getenv("BOT_USERNAME"),
		os.Getenv("TMI_OAUTH_TOKEN"))

	// Repositories
	notificationsRepository := storage.NewNotificationsRepository(db)
	messagesRepository := storage.NewMessagesRepository(db)

	err = notificationsRepository.CreateTables()
	if err != nil {
		return err
	}

	messagesRepository.CreateTables()
	if err != nil {
		return err
	}


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
			notificationsRepository,
			messagesRepository,
		)
		if err != nil {
			log.Printf(utils.Colors.Red, err.Error(), utils.Colors.Reset)
			return err
		}
	}

	// Notifications
	notifications.SetupRoutes(server)
	go notifications.SendMessages()


	// Run web server
	go func() {
		server.Run()
	}()



	// On Private Message
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		log.Println(utils.Colors.Blue + "Private Message", message.Channel, message.User.Name, message.Message, utils.Colors.Reset)

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

	log.Println("Join Twitch client to channel:", os.Getenv("CHANNEL"))
	client.Join(os.Getenv("CHANNEL"))

	return client.Connect()
}
