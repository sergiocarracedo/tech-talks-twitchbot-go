package commands

import (
	"encoding/json"
	"fmt"
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/goodsign/monday"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sergiocarracedo.es/streambot-go/vigotech"
	"sort"
	"strings"
	"time"
)

type EventsCommand CommandInterface

var _ EventsCommand = new(Command)

type Event struct {
	Title string
	Date  int64
	Link  string
}

type EventsSource struct {
	Link   string  `json:"link"`
	Events []Event `json:"events"`
}

func getLocalDefinedEvents() EventsSource {
	jsonFile, err := os.Open("data/events.json")
	defer jsonFile.Close()

	if err != nil {
		log.Println(err)
	}

	var byteValue []byte
	byteValue, err = ioutil.ReadAll(jsonFile)

	var eventsSource EventsSource

	if err == nil {
		json.Unmarshal(byteValue, &eventsSource)
	}
	return eventsSource
}

func NewEventsCommand(client *twitch.Client, vigotechData vigotech.VigoTechGroup) *Command {
	var events []Event

	// Get local defined events
	localEventsSource := getLocalDefinedEvents()
	events = append(events, localEventsSource.Events...)

	// Get videos from VigoTech data source
	for _, event := range vigotech.GetNextEvents(vigotechData) {
		events = append(events, Event{event.Title, event.Date / 1000, event.Link})
	}

	// Sort videos by date
	sort.Slice(events, func(i, j int) bool {
		return events[i].Date < events[j].Date
	})

	return &Command{
		Id:     "events",
		Name:   "eventos",
		client: client,
		handler: func(client *twitch.Client, message twitch.PrivateMessage) error {
			log.Printf("Videos command: Channel %s - User %s", message.Channel, message.User.Name)

			// Remove pass events
			var nextEvents []Event
			for _, event := range events {
				if event.Date >= time.Now().Unix() {
					nextEvents = append(nextEvents, event)
				}
			}

			numEvents := int(math.Min(float64(len(nextEvents)), 5))

			if numEvents > 1 {
				client.Say(message.Channel, "Próximos eventos")
			} else if numEvents == 1 {
				client.Say(message.Channel, "Próximo evento")
			}
			for i := 0; i < numEvents; i++ {
				event := nextEvents[i]
				var messageContent []string
				tm := time.Unix(event.Date, 0)
				messageContent = append(messageContent, "📅 "+monday.Format(tm, "Mon, _2 Jan 15:04", monday.LocaleEsES))
				messageContent = append(messageContent, " - 🔊 "+event.Title)
				messageContent = append(messageContent, " => ")
				messageContent = append(messageContent, event.Link)

				client.Say(message.Channel, strings.Join(messageContent, " "))
			}

			if localEventsSource.Link != "" {
				client.Say(message.Channel, fmt.Sprintf("ℹ️ Más eventos e info ⏩ %s", localEventsSource.Link))
			}

			return nil

		},
	}
}
