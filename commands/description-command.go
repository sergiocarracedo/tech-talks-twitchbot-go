package commands

import (
	"encoding/json"
	"github.com/gempir/go-twitch-irc/v2"
	"io/ioutil"
	"log"
	"os"
)

type DescriptionCommand CommandInterface

var _ DescriptionCommand = new(Command)

func NewDescriptionCommand(client *twitch.Client) *Command {
	jsonFile, err := os.Open("data/description.json")
	defer jsonFile.Close()

	if err != nil {
		log.Println(err)
	}

	var byteValue []byte
	byteValue, err = ioutil.ReadAll(jsonFile)

	var description string

	if err == nil {
		json.Unmarshal(byteValue, &description)
	}


	return &Command{
		Id:     "description",
		Name:   "descripcion",
		client: client,
		handler: func(client *twitch.Client, message twitch.PrivateMessage) error {
			log.Println("Description command:", message.Channel, message.User.Name, description)
			client.Say(message.Channel, description)

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


