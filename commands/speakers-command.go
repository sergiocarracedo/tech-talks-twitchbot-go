package commands

import (
	"encoding/json"
	"github.com/gempir/go-twitch-irc/v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// ========================== APIKeyService =========================
type SpeakersCommand CommandInterface

var _ SpeakersCommand = new(Command)

type Speaker struct {
	Name string `json:"name"`
	Bio string `json:"bio"`
	Social []string `json:"social"`
}

func NewSpeakersCommand(client *twitch.Client) *Command{
	jsonFile, err := os.Open("data/speakers.json")
	defer jsonFile.Close()

	if err != nil {
		log.Println(err)
	}

	var byteValue []byte
	byteValue, err = ioutil.ReadAll(jsonFile)

	var speakers []Speaker

	if err == nil {
		json.Unmarshal(byteValue, &speakers)
	}

	return &Command{
		Id: "speakers",
		Name: "ponentes",
		client: client,
		handler:  func(client *twitch.Client, message twitch.PrivateMessage) error {
			log.Printf("Speaker command: Channel %s - User %s", message.Channel, message.User.Name)
			for _, speaker := range speakers {
				var messageContent  []string

				messageContent = append(messageContent,"🗣️ " + speaker.Name)
				messageContent = append(messageContent, " => ")
				messageContent = append(messageContent, speaker.Bio)

				messageContent = append(messageContent, strings.Join(speaker.Social, " - 🎯 "))

				client.Say(message.Channel, strings.Join(messageContent, " "))
			}
			return nil
		},
	}
}

