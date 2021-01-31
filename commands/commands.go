package commands

import (
	"github.com/gempir/go-twitch-irc/v2"
	"sergiocarracedo.es/streambot-go/vigotech"
)

type CommandInterface interface {
	Handler(message twitch.PrivateMessage) error
}

type Command struct {
	Id      string
	Name    string
	client  *twitch.Client
	handler func(client *twitch.Client, message twitch.PrivateMessage) error
}

func GetCommands(client *twitch.Client) []*Command {
	vigotechData, _ := vigotech.GetData()
	return []*Command{
		NewDescriptionCommand(client),
		NewSpeakersCommand(client),
		NewVideosCommand(client, vigotechData),
		NewEventsCommand(client, vigotechData),
	}
}

func (c *Command) Handler(message twitch.PrivateMessage) error {
	return c.handler(c.client, message)
}
