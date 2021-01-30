package commands

import (
	"fmt"
	"github.com/gempir/go-twitch-irc/v2"
	"log"
	"math"
	"os"
	"sergiocarracedo.es/streambot-go/vigotech"
	"strings"
)

type VideosCommand CommandInterface

var _ VideosCommand = new(Command)


func NewVideosCommand(client *twitch.Client) *Command{
	vigotechData, _ := vigotech.GetData()
	group := vigotech.GetGroup(vigotechData, os.Getenv("VIGOTECH_GROUP"))
	videos := vigotech.GetGroupVideos(group)

	log.Println(videos)

	return &Command{
		Id: "videos",
		Name: "videos",
		client: client,
		handler:  func(client *twitch.Client, message twitch.PrivateMessage) error {
			log.Printf("Videos command: Channel %s - User %s", message.Channel, message.User.Name)
			numVideos := int(math.Min(float64(len(videos)), 1))

			if numVideos > 1 {
				client.Say(message.Channel, fmt.Sprintf("Últimos %d videos", numVideos))
			} else if numVideos == 1 {
				client.Say(message.Channel, "Último video")
			}
			for i := 0; i < numVideos; i++  {
				video := videos[i]
				var messageContent []string
				messageContent = append(messageContent, "📺️ "+video.Title)
				messageContent = append(messageContent, " => ")
				messageContent = append(messageContent, video.Link)

				client.Say(message.Channel, strings.Join(messageContent, " "))
			}

			youtube, ok := group.Links["youtube"]
			if ok {
				client.Say(message.Channel, fmt.Sprintf("Más videos en ⏩ %s", youtube))
			}

			return nil

		},
	}
}

