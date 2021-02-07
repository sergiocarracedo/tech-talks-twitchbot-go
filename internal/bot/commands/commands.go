package commands

import (
	"github.com/gempir/go-twitch-irc/v2"
	"log"
	"os"
	"sergiocarracedo.es/streambot-go/vigotech"
	"strings"
)

func Help(client *twitch.Client, channel string, commandList []*Command) {
	commandNames := ListNames(commandList)

	for i, _ := range commandNames {
		commandNames[i] = "!" + commandNames[i]
	}

	log.Println("Sending help message")
	client.Say(channel, "Comandos disponibles: "+strings.Join(commandNames, ", ")+". Solo se muestra respuesta de cada comando cada "+os.Getenv("COMMAND_COLD_DOWN_TIME")+" segundos")
}

func ListNames(commands []*Command) []string {
	var names []string
	for _, command := range commands {
		names = append(names, command.Name)
	}
	return names
}

func ListIds(commands []*Command) []string {
	var ids []string
	for _, command := range commands {
		ids = append(ids, command.Id)
	}
	return ids
}

type CommandInterface interface {
	Handler(message twitch.PrivateMessage) error
}

type Command struct {
	Id      string
	Name    string
	client  *twitch.Client
	handler func(client *twitch.Client, message twitch.PrivateMessage) error
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func GetCommands(client *twitch.Client, disabledCommands []string) []*Command {
	vigotechData, _ := vigotech.GetData()

	availableCommands := []*Command{
		NewDescriptionCommand(client),
		NewSpeakersCommand(client),
		NewVideosCommand(client, vigotechData),
		NewEventsCommand(client, vigotechData),
		NewSocialCommand(client),
	}

	var commands []*Command

	for _, command := range availableCommands {
		if !contains(disabledCommands, command.Id) {
			commands = append(commands, command)
		}
	}

	return commands
}

func (c *Command) Handler(message twitch.PrivateMessage) error {
	return c.handler(c.client, message)
}
