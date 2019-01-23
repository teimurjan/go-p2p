package imstorage

import "github.com/teimurjan/go-p2p/models"

// Storage is a basic imstorage interface
type Storage interface {
	AddNotificationToSend(notification *models.Notification) error
	AddNotificationToHandle(notification *models.Notification) error
	SubscribeToNotificationsToHandle() error
	SubscribeToNotificationsToSend() error
	GetNotificationsToHandle() chan models.Notification
	GetNotificationsToSend() chan models.Notification
}
