package chatLogger

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Service struct {
	client *twitch.Client
	db *sqlx.DB
}

func New(client *twitch.Client, db *sqlx.DB) (*Service, error) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS main.messages(" +
		"id INTEGER PRIMARY KEY," +
		"channel TEXT," +
		"user TEXT," +
		"message TEXT," +
		"twitch_id TEXT, " +
		"room_id TEXT," +
		"user_id TEXT," +
		"user_display_name TEXT," +
		"user_color TEXT," +
		"broadcaster INTEGER," +
		"premium INTEGER," +
		"type TEXT," +
		"time INT," +
		"action INT," +
		"bits INT" +
		")")
	if err != nil {
		log.Printf(err.Error())
		return &Service{}, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS main.notifications(" +
		"id INTEGER PRIMARY KEY," +
		"channel TEXT," +
		"user TEXT," +
		"message_id TEXT," +
		"message TEXT," +
		"system_message TEXT," +
		"twitch_id TEXT, " +
		"room_id TEXT," +
		"user_id TEXT," +
		"user_display_name TEXT," +
		"user_color TEXT," +
		"broadcaster INTEGER," +
		"premium INTEGER," +
		"subscriber INTEGER," +
		"type TEXT," +
		"time INT" +
		")")
	if err != nil {
		log.Printf(err.Error())
		return &Service{}, err
	}

	return &Service{client, db}, nil
}

func (s *Service) OnPrivateMessage(message twitch.PrivateMessage) error {
	query, args, _ := sq.
		Insert("messages").
		Columns(
			"channel",
			"room_id",
			"twitch_id",
			"user_id",
			"user",
			"user_display_name",
			"user_color",
			"message",
			"broadcaster",
			"premium",
			"type",
			"time",
			"action",
			"bits").
		Values(
			message.Channel,
			message.RoomID,
			message.ID,
			message.User.ID,
			message.User.Name,
			message.User.DisplayName,
			message.User.Color,
			message.Message,
			message.User.Badges["broadcaster"],
			message.User.Badges["premium"],
			message.RawType,
			message.Time.Unix(),
			message.Action,
			message.Bits).
		ToSql()
	_, err := s.db.Exec(query, args...)

	if err != nil {
		return err
	}
	return nil
}


func (s *Service) OnUserNoticeMessage(message twitch.UserNoticeMessage) error {
	query, args, _ := sq.
		Insert("messages").
		Columns(
			"channel",
			"room_id",
			"twitch_id",
			"user_id",
			"user",
			"user_display_name",
			"user_color",
			"message",
			"message_id",
			"system_message",
			"broadcaster",
			"premium",
			"type",
			"time",
			"subscriber").
		Values(
			message.Channel,
			message.RoomID,
			message.ID,
			message.User.ID,
			message.User.Name,
			message.User.DisplayName,
			message.User.Color,
			message.Message,
			message.MsgID,
			message.SystemMsg,
			message.User.Badges["broadcaster"],
			message.User.Badges["premium"],
			message.RawType,
			message.Time.Unix(),
			message.User.Badges["subscriber"]).
		ToSql()
	_, err := s.db.Exec(query, args...)

	if err != nil {
		return err
	}
	return nil
}
