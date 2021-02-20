package commands

import (
	"encoding/json"
	"fmt"
	"github.com/gempir/go-twitch-irc/v2"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sergiocarracedo.es/streambot-go/vigotech"
	"sort"
	"strings"
)

type VideosCommand CommandInterface

type Video struct {
	Title string
	Date  int64
	Link  string
}

type VideosSource struct {
	Link   string  `json:"link"`
	Videos []Video `json:"videos"`
}

func getLocalDefinedVideos() VideosSource {
	jsonFile, err := os.Open("data/videos.json")
	defer jsonFile.Close()

	if err != nil {
		log.Println(err)
	}

	var byteValue []byte
	byteValue, err = ioutil.ReadAll(jsonFile)

	var videosSource VideosSource

	if err == nil {
		json.Unmarshal(byteValue, &videosSource)
	}

	return videosSource
}

func NewVideosCommand(client *twitch.Client, vigotechData vigotech.VigoTechGroup) *Command {
	var videos []Video

	// Get local defined videos
	localVideosSource := getLocalDefinedVideos()
	videos = append(videos, localVideosSource.Videos...)

	// Get videos from VigoTech data source
	group := vigotech.GetGroup(vigotechData, os.Getenv("VIGOTECH_GROUP"))
	for _, video := range vigotech.GetGroupVideos(group) {
		videos = append(videos, Video{video.Title, video.Date / 1000, video.Link})
	}

	// Sort videos by date
	sort.Slice(videos, func(i, j int) bool {
		return videos[i].Date > videos[j].Date
	})

	return &Command{
		Id:     "videos",
		Name:   "videos",
		client: client,
		handler: func(client *twitch.Client, message twitch.PrivateMessage) error {
			log.Printf("Videos command: Channel %s - User %s", message.Channel, message.User.Name)
			numVideos := int(math.Min(float64(len(videos)), 1))

			if numVideos > 1 {
				client.Say(message.Channel, fmt.Sprintf("Últimos %d videos", numVideos))
			} else if numVideos == 1 {
				client.Say(message.Channel, "Último video")
			}
			for i := 0; i < numVideos; i++ {
				video := videos[i]
				var messageContent []string
				messageContent = append(messageContent, "📺️ "+video.Title)
				messageContent = append(messageContent, " => ")
				messageContent = append(messageContent, video.Link)

				client.Say(message.Channel, strings.Join(messageContent, " "))
			}

			var videoSourceLink string
			if localVideosSource.Link != "" {
				videoSourceLink = localVideosSource.Link
			} else {
				if youtube, ok := group.Links["youtube"]; ok {
					videoSourceLink = youtube
				}
			}

			if videoSourceLink != "" {
				client.Say(message.Channel, fmt.Sprintf("Más videos en ⏩ %s", videoSourceLink))
			}

			return nil

		},
	}
}
