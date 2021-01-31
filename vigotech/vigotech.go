package vigotech

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"time"
)

type VigoTechGroup struct {
	Name      string                   `json:"name"`
	Logo      string                   `json:"logo"`
	Sticker   string                   `json:"sticker"`
	Links     map[string]string        `json:"links"`
	NextEvent VigoTechEvent            `json:"nextEvent"`
	VideoList []VigoTechVideo          `json:"videoList"`
	Members   map[string]VigoTechGroup `json:"members"`
}

type VigoTechEvent struct {
	Title    string `json:"title"`
	Date     int64  `json:"date"`
	Link     string `json:"url"`
	location string `json:"location"`
}

type VigoTechVideo struct {
	Title  string `json:"title"`
	Player string `json:"player"`
	Date   int64  `json:"pubDate"`
	Id     string `json:"id"`
	Link   string
}

type VigoTechVideoSource struct {
	Type      string `json:"type"`
	ChannelId string `json:"channel_id"`
	Source    string `json:"source""`
}

func GetData() (VigoTechGroup, error) {
	return GetJson()
}
func GetGroup(group VigoTechGroup, groupName string) VigoTechGroup {
	if group.Name == groupName {
		return group
	} else {
		for _, member := range group.Members {
			subgroup := GetGroup(member, groupName)
			if subgroup.Name == groupName {
				return subgroup
			}
		}
		return VigoTechGroup{}
	}
}

func GetNextEvents(group VigoTechGroup) []VigoTechEvent {
	var getGroupNextEvents func(group VigoTechGroup) []VigoTechEvent
	getGroupNextEvents = func(group VigoTechGroup) []VigoTechEvent {
		var events []VigoTechEvent
		if (group.NextEvent != VigoTechEvent{}) {
			events = append(events, group.NextEvent)
		}

		if len(group.Members) > 0 {
			for _, subgroup := range group.Members {
				events = append(events, getGroupNextEvents(subgroup)...)
			}
		}
		return events
	}

	return getGroupNextEvents(group)
}

func GetGroupVideos(group VigoTechGroup) []VigoTechVideo {
	videos := group.VideoList

	for i := range videos {
		videos[i].Link = "https://youtube.com/v/" + videos[i].Id
	}

	sort.Slice(videos, func(i, j int) bool {
		return videos[i].Date > videos[j].Date
	})
	return videos
}

func GetJson() (VigoTechGroup, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	r, err := client.Get("https://vigotech.org/vigotech-generated.json")
	if err != nil {
		log.Println(err.Error())
		return VigoTechGroup{}, err
	}
	defer r.Body.Close()

	var target VigoTechGroup

	err = json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		log.Println(r.Body, err.Error())
		return VigoTechGroup{}, err
	}
	return target, nil
}
