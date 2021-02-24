package storage

import (
	"context"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"sergiocarracedo.es/streambot-go/internal/messages"
)

type MessagesRepository struct {
	db *sqlx.DB
}

func NewMessagesRepository(db *sqlx.DB) *MessagesRepository {
	return &MessagesRepository{db }
}

const messagesTable = "main.messages"

func (r *MessagesRepository) CreateTables() error {
	_, err := r.db.Exec("CREATE TABLE IF NOT EXISTS "+messagesTable+"(" +
		"id INTEGER PRIMARY KEY," +
		"channel TEXT," +
		"username TEXT," +
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
	return err
}

func (r *MessagesRepository) Save(ctx context.Context, message messages.Message) error {
	ds := goqu.Insert(messagesTable).Rows(message)
	sql, args, _ := ds.ToSQL()

	_, err := r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("error trying to persist message on database: %v", err)
	}

	return nil
}