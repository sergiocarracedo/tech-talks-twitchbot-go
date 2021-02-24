package notifications

import (
	"time"
)

type Notification struct {
	Id int `json:"id" db:"id" goqu:"skipinsert,skipupdate"`
	TwitchId string `json:"twitch_id" db:"twitch_id"`
	Channel string `json:"channel" db:"channel"`
	RoomID string `json:"room_id" db:"room_id"`
	UserName string `json:"username" db:"username"`
	UserId string `json:"user_id" db:"user_id"`
	UserDisplayName string `json:"user_display_name" db:"user_display_name"`
	UserColor string `json:"user_color" db:"user_color"`
	Message string `json:"message" db:"message"`
	MessageId string `json:"message_id" db:"message_id"`
	SystemMsg string `json:"system_message" db:"system_message"`
	Broadcaster int `json:"broadcaster" db:"broadcaster"`
	Premium int `json:"premium" db:"premium"`
	Subscriber int `json:"susbcriber" db:"subscriber"`
	NotificationType string `json:"type" db:"type"`
	Time time.Time `json:"time" db:"time"`
	Notified int `json:"notified" db:"notified"`
}

type NotificationsRepository interface {
	CreateTables() error
	Save(notification Notification) error
	FindLastNotNotified() (Notification, error)
	SetNotified(id int) error
}