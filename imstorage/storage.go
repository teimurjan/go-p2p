package imstorage

import (
	"github.com/teimurjan/go-p2p/models"
)

type subscribeCallback func()

// Storage is a basic imstorage interface
type Storage interface {
	AddNotificationToSend(notification *models.Notification) error
	AddNotificationToHandle(notification *models.Notification) error
	SubscribeToNotificationsToHandle(cb subscribeCallback) error
	SubscribeToNotificationsToSend(cb subscribeCallback) error
	GetNotificationsToHandle() chan models.Notification
	GetNotificationsToSend() chan models.Notification
}
