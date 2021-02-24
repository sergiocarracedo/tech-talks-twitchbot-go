package chatLogger

import (
	"context"
	"github.com/gempir/go-twitch-irc/v2"
	"log"
	"sergiocarracedo.es/streambot-go/internal/messages"
	"sergiocarracedo.es/streambot-go/internal/notifications"
)

type Service struct {
	client                  *twitch.Client
	notificationsRepository notifications.NotificationsRepository
	messagesRepository      messages.MessagesRepository
}

func New(
	client *twitch.Client,
	notificationsRepository notifications.NotificationsRepository,
	messagesRepository messages.MessagesRepository,
) (*Service, error) {
	return &Service{
		client,
		notificationsRepository,
		messagesRepository,
	}, nil
}

func (s *Service) OnPrivateMessage(message twitch.PrivateMessage) error {
	err := s.messagesRepository.Save(
		context.Background(),
		messages.Message{
			0,
			message.ID,
			message.Channel,
			message.RoomID,
			message.User.Name,
			message.User.ID,
			message.User.DisplayName,
			message.User.Color,
			message.Message,
			message.User.Badges["broadcaster"],
			message.User.Badges["premium"],
			message.RawType,
			message.Time,
			message.Action,
			message.Bits,
		})

	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (s *Service) OnUserNoticeMessage(message twitch.UserNoticeMessage) error {
	err := s.notificationsRepository.Save(
		notifications.Notification{
			0,
			message.ID,
			message.Channel,
			message.RoomID,
			message.User.Name,
			message.User.ID,
			message.User.DisplayName,
			message.User.Color,
			message.Message,
			message.MsgID,
			message.SystemMsg,
			message.User.Badges["broadcaster"],
			message.User.Badges["premium"],
			message.User.Badges["subscriber"],
			message.RawType,
			message.Time,
			0,
		})

	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
