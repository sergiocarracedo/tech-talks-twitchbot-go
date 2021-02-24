package storage

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"sergiocarracedo.es/streambot-go/internal/notifications"
	"time"
)

type NotificationsRepository struct {
	db *sqlx.DB
}

func NewNotificationsRepository(db *sqlx.DB) *NotificationsRepository {
	return &NotificationsRepository{db }
}

const notificationsTable = "notifications"


func (r *NotificationsRepository) CreateTables() error {
	_, err := r.db.Exec("CREATE TABLE IF NOT EXISTS "+notificationsTable+"(" +
		"id INTEGER PRIMARY KEY," +
		"channel TEXT," +
		"username TEXT," +
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
		"time INTEGER," +
		"notified INTEGER" +
		"notify_time INTEGER" +
		")")
	return err
}

func (r *NotificationsRepository) Save(notification notifications.Notification) error {
	ds := goqu.Insert(notificationsTable).Rows(notification)
	sql, args, _ := ds.ToSQL()

	_, err := r.db.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("error trying to persist notification on database: %v", err)
	}

	return nil
}

func (r *NotificationsRepository) FindLastNotNotified() (notification notifications.Notification, err error) {
	query := sq.Select("*").
		From(notificationsTable).
		Where("notified = ?", 1).
		OrderBy("id ASC").
		Limit(1).
		RunWith(r.db)

	err = query.QueryRow().Scan(&notification)
	return
}

func (r *NotificationsRepository) SetNotified(id int) error {
	query := sq.Update(notificationsTable).
		Set("notified", 1).
		Set("notify_time", time.Now().Unix()).
		Where("id = ?", id).
		RunWith(r.db)

	_, err := query.Exec()

	return err
}