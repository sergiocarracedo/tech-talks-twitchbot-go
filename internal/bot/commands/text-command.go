package commands

import (
	"encoding/json"
	"github.com/gempir/go-twitch-irc/v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type TextCommand CommandInterface

type TextCommandData struct {
	Command string `json:"command"`
	Text    string `json:"text"`
}

func NewTextCommand(client *twitch.Client) []*Command {
	jsonFile, err := os.Open("data/text.json")
	defer jsonFile.Close()

	if err != nil {
		log.Println(err)
	}

	var byteValue []byte
	byteValue, err = ioutil.ReadAll(jsonFile)

	var textCommands []TextCommandData
	var commands []*Command

	if err == nil {
		err = json.Unmarshal(byteValue, &textCommands)

		if err != nil {
			log.Println(err.Error())
			return commands
		}

		for _, textCommand := range textCommands {
			commands = append(commands, &Command{
				Id:     textCommand.Command,
				Name:   textCommand.Command,
				client: client,
				handler: func(client *twitch.Client, message twitch.PrivateMessage) error {
					messageContent := strings.Trim(strings.ToLower(message.Message), " !")
					log.Println("Text command ("+messageContent+"):", message.Message, message.Channel, message.User.Name)
					for _, textCommand := range textCommands {
						if strings.ToLower(textCommand.Command) == messageContent {
							client.Say(message.Channel, textCommand.Text)
							break
						}
					}
					return nil
				},
			})
		}
	}

	return commands
}
