package commands

import (
	"encoding/json"
	"github.com/gempir/go-twitch-irc/v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type SocialCommand CommandInterface

var _ SocialCommand = new(Command)

type Social struct {
	Name   string   `json:"name"`
	Links  []string `json:"links"`
	Spacer string   `json:"spacer"`
}

func NewSocialCommand(client *twitch.Client) *Command {
	jsonFile, err := os.Open("data/social.json")
	defer jsonFile.Close()

	if err != nil {
		log.Println(err)
	}

	var byteValue []byte
	byteValue, err = ioutil.ReadAll(jsonFile)

	var social []Social

	if err == nil {
		json.Unmarshal(byteValue, &social)
	}

	return &Command{
		Id:     "social",
		Name:   "social",
		client: client,
		handler: func(client *twitch.Client, message twitch.PrivateMessage) error {

			var messageContent []string

			for _, item := range social {
				var messageItemContent []string
				if item.Spacer != "" {
					messageItemContent = append(messageItemContent, item.Spacer)
				}
				messageItemContent = append(messageItemContent, item.Name+" => ")
				messageItemContent = append(messageItemContent, strings.Join(item.Links, " "))

				messageContent = append(messageContent, strings.Join(messageItemContent, " "))
			}

			say := strings.Join(messageContent, " ·|· ")
			client.Say(message.Channel, say)

			log.Println("Social command:", message.Channel, message.User.Name, say)
			return nil
		},
	}
}
