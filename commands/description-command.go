package commands

import (
	"github.com/gempir/go-twitch-irc/v2"
	"log"
	"os"
)

// ========================== APIKeyService =========================

type DescriptionCommand CommandInterface

var _ DescriptionCommand = new(Command)

func NewDescriptionCommand(client *twitch.Client) *Command {
	return &Command{
		Id:     "description",
		Name:   "descripcion",
		client: client,
		handler: func(client *twitch.Client, message twitch.PrivateMessage) error {
			log.Println("Description command:", message.Channel, message.User.Name, os.Getenv("COMMAND_DESCRIPTION"))
			client.Say(message.Channel, os.Getenv("COMMAND_DESCRIPTION"))

			return nil
		},
	}
}

//
//func (c *CommandDescription) Handler(message twitch.PrivateMessage) error{
//	log.Println("Description command:", message.Channel, message.User.Name, os.Getenv("COMMAND_DESCRIPTION"))
//	c.client.Say(message.Channel, os.Getenv("COMMAND_DESCRIPTION"))
//
//	return nil
//}


