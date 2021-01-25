package commands

import (
	"github.com/gempir/go-twitch-irc/v2"
	"log"
	"os"
)

// ========================== APIKeyService =========================
type SpeakersCommand CommandInterface

var _ SpeakersCommand = new(Command)

func NewSpeakersCommand(client *twitch.Client) *Command{
	return &Command{
		Id: "speakers",
		Name: "speakers",
		client: client,
		handler:  func(client *twitch.Client, message twitch.PrivateMessage) error {
			log.Printf("Speaker command: Channel %s - User %s", message.Channel, message.User.Name)
			client.Say(message.Channel, os.Getenv("COMMAND_SPEAKERS"))
			return nil
		},
	}
}

