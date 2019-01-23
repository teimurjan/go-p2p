package imstorage

import (
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/teimurjan/go-p2p/models"
)

const (
	// NotificationToSendKey is a notification to send redis key
	NotificationToSendKey = "notificationToSend"
	// NotificationToHandleKey is a notification to handle redis key
	NotificationToHandleKey = "notificationToHandle"
)

type redisStorage struct {
	pool                  *redis.Pool
	notificationsToSend   chan models.Notification
	notificationsToHandle chan models.Notification
}

// NewRedisStorage creates new storage instance
func NewRedisStorage(URL string) Storage {
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", URL)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	notificationsToHandle := make(chan models.Notification, 10)
	notificationsToSend := make(chan models.Notification, 10)
	return &redisStorage{
		pool,
		notificationsToHandle,
		notificationsToSend,
	}
}

func (s *redisStorage) AddNotificationToSend(notification *models.Notification) error {
	json, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	conn := s.pool.Get()
	defer conn.Close()

	_, err = conn.Do("PUBLISH", NotificationToSendKey, json)
	return err
}

func (s *redisStorage) AddNotificationToHandle(notification *models.Notification) error {
	json, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	conn := s.pool.Get()
	defer conn.Close()

	_, err = conn.Do("PUBLISH", NotificationToHandleKey, json)
	return err
}

func (s *redisStorage) SubscribeToNotificationsToHandle() error {
	conn := s.pool.Get()
	defer conn.Close()

	psc := redis.PubSubConn{Conn: conn}
	psc.PSubscribe(NotificationToHandleKey)

	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			notification := models.Notification{}
			err := json.Unmarshal(v.Data, &notification)
			if err == nil {
				s.notificationsToHandle <- notification
			}
		case error:
			return v
		}
	}
}

func (s *redisStorage) SubscribeToNotificationsToSend() error {
	conn := s.pool.Get()
	defer conn.Close()

	psc := redis.PubSubConn{Conn: conn}
	psc.PSubscribe(NotificationToSendKey)
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			notification := models.Notification{}
			err := json.Unmarshal(v.Data, &notification)
			if err == nil {
				s.notificationsToSend <- notification
			}
		case error:
			return v
		}
	}
}

func (s *redisStorage) GetNotificationsToHandle() chan models.Notification {
	return s.notificationsToHandle
}

func (s *redisStorage) GetNotificationsToSend() chan models.Notification {
	return s.notificationsToSend
}
