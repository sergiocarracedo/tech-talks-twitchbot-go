package messages

import (
	"context"
	"time"
)

type Message struct {
	Id int `json:"id" db:"id" goqu:"skipinsert,skipupdate"`
	TwitchId string `json:"twitch_id" db:"twitch_id"`
	Channel string `json:"channel" db:"channel"`
	RoomID string `json:"room_id" db:"room_id"`
	UserName string `json:"username" db:"username"`
	UserId string `json:"user_id" db:"user_id"`
	UserDisplayName string `json:"user_display_name" db:"user_display_name"`
	UserColor string `json:"user_color" db:"user_color"`
	Message string `json:"message" db:"message"`
	Broadcaster int `json:"broadcaster" db:"broadcaster"`
	Premium int `json:"premium" db:"premium"`
	RawType string `json:"type" db:"type"`
	Time time.Time `json:"time" db:"time"`
	Action bool `json:"action" db:"action"`
	Bits int `json:"bits" db:"bits"`
}

type MessagesRepository interface {
	CreateTables() error
	Save(ctx context.Context, message Message) error
}
