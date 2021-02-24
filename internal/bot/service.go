package bot

import (
	"github.com/gempir/go-twitch-irc/v2"
	"sergiocarracedo.es/streambot-go/internal/bot/commands"
	"sergiocarracedo.es/streambot-go/vigotech"
	"strings"
	"time"
)

type Service struct {
	client               *twitch.Client
	channel              string
	disabledCommands     []string
	lastRunTime          map[string]int64
	coldDownTime         int64
	sendHelpMessageEvery int64
	commandList          []*commands.Command
}

func New(client *twitch.Client, channel string, disabledCommands []string, coldDownTime int64, sendHelpMessageEvery int64) Service {
	vigotech.GetJson()

	return Service{
		client,
		channel,
		disabledCommands,
		make(map[string]int64),
		coldDownTime,
		sendHelpMessageEvery,
		commands.GetCommands(client, disabledCommands),
	}
}

func (b Service) GetCommandListNames() []string {
	return commands.ListNames(b.commandList)
}

func (b Service) OnPrivateMessage(message twitch.PrivateMessage) bool {
	messageContent := strings.Trim(strings.ToLower(message.Message), " ")

	// Check if message is a command and run it
	if messageContent == "!help" || messageContent == "!ayuda" {
		commands.Help(b.client, b.channel, b.commandList)
	} else {
		for _, command := range b.commandList {
			if messageContent == "!"+command.Name {
				commandLastRunTime, ok := b.lastRunTime[command.Id]
				if !ok || time.Now().Unix() >= commandLastRunTime+b.coldDownTime {
					command.Handler(message)
					b.lastRunTime[command.Id] = time.Now().Unix()
					return true
				} else {
					return false
				}
			}
		}
	}

	return true
}

func (b Service) OnConnect() {
	if b.sendHelpMessageEvery > 0 {
		commands.Help(b.client, b.channel, b.commandList)
		go func() {
			for _ = range time.NewTicker(time.Duration(b.sendHelpMessageEvery) * time.Second).C {
				commands.Help(b.client, b.channel, b.commandList)
			}
		}()
	}
}
