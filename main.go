package main

import (
	"fmt"
	twitch "github.com/gempir/go-twitch-irc/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sergiocarracedo.es/streambot-go/commands"
	"sergiocarracedo.es/streambot-go/vigotech"
	"strconv"
	"strings"
	"time"
)

func help(client *twitch.Client, channel string, commandList []*commands.Command) {
	var commandNames []string

	for _, command := range commandList {
		commandNames = append(commandNames, "!" + command.Name)
	}

	log.Println("Sending help message")
	client.Say(channel, "Comandos disponibles: " + strings.Join(commandNames, ", ") + ". Solo se muestra respuesta de cada comando cada " + os.Getenv("COMMAND_COLD_DOWN_TIME") + "segundos")
}


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	lastRunTime := make(map[string]int64)

	client := twitch.NewClient(os.Getenv("BOT_USERNAME"), os.Getenv("TMI_OAUTH_TOKEN"))

	commandList := commands.GetCommands(client)

	vigotech.GetJson()

	fmt.Println("Stating bot...")
	coldDownTime, _ := strconv.ParseInt(os.Getenv("COMMAND_COLD_DOWN_TIME"), 10, 64)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Println("PrivateMessage", message.Channel, message.User, message.Message)

		for _, command := range commandList {
			if message.Message == "!" + command.Name {
				commandLastRunTime, ok := lastRunTime[command.Id]
				log.Println(commandLastRunTime, ok, time.Now().Unix() + coldDownTime)
				if !ok || time.Now().Unix() >= commandLastRunTime + coldDownTime {
					command.Handler(message)
					lastRunTime[command.Id] = time.Now().Unix()
				} else {
					log.Println("Cold down", message.User, message.Message)
				}
			}
		}
	})

	client.OnUserStateMessage(func(message twitch.UserStateMessage) {
		fmt.Println("UserStateMessage", message.Channel, message.User, message.Message)
	})

	client.OnUserNoticeMessage(func(message twitch.UserNoticeMessage) {
		fmt.Println("UserNoticeMessage", message.Channel, message.User, message.Message)
	})

	client.OnUserJoinMessage(func(message twitch.UserJoinMessage) {
		fmt.Println("UserJoinMessage", message.Channel, message.User)
	})

	client.OnUserPartMessage(func(message twitch.UserPartMessage) {
		fmt.Println("UserPartMessage", message.Channel, message.User)
	})

	client.OnNoticeMessage(func(message twitch.NoticeMessage) {
		fmt.Println("NoticeMessage", message.Channel)
	})

	client.Join(os.Getenv("CHANNEL"))

	sendHelpMessageEvery, _ := strconv.ParseInt(os.Getenv("COMMAND_HELP_EVERY"), 10, 64)

	client.OnConnect(func() {
		help(client, os.Getenv("CHANNEL"), commandList)
		go func () {
			for _ = range time.NewTicker(time.Duration(sendHelpMessageEvery) * time.Second).C {
				help(client, os.Getenv("CHANNEL"), commandList)
			}
		}()
	})

	err = client.Connect()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected")





}
