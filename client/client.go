package client

import (
	"encoding/json"
	"net"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/teimurjan/go-p2p/protocol"
)

// Client is a P2P client interface
type Client interface {
	StartNotifier()
	StartNotificationListener()
	NotifyNetwork(message *protocol.Notification)
}

type client struct {
	Port                  string
	logger                *logrus.Logger
	ReceivedNotifications chan protocol.Notification
	notificationsToSend   chan protocol.Notification
}

// NewClient creates new client instance
func NewClient(port string, logger *logrus.Logger) *client {
	ReceivedNotifications := make(chan protocol.Notification, 10)
	notificationsToSend := make(chan protocol.Notification, 10)
	return &client{
		port,
		logger,
		ReceivedNotifications,
		notificationsToSend,
	}
}

func (c *client) StartNotifier() {
	destinationAddress, _ := net.ResolveUDPAddr("udp", "255.255.255.255:"+c.Port)
	connection, err := net.DialUDP("udp", nil, destinationAddress)
	if err != nil {
		c.logger.Error(err)
		os.Exit(1)
	}
	defer connection.Close()

	c.logger.Println("Notifier started")

	for {
		notification := <-c.notificationsToSend

		encodedNotification, err := json.Marshal(notification)
		if err != nil {
			c.logger.Error(err)
			os.Exit(1)
		}
		connection.Write(encodedNotification)

		c.logger.Println("Notification sent: ", notification)
	}

}

func (c *client) StartNotificationListener() {
	localAddress, _ := net.ResolveUDPAddr("udp", ":"+c.Port)
	connection, err := net.ListenUDP("udp", localAddress)
	if err != nil {
		c.logger.Error(err)
	}
	defer connection.Close()

	c.logger.Println("Notifications listener started")

	var notification protocol.Notification
	for {
		inputBytes := make([]byte, 4096)
		length, _, _ := connection.ReadFromUDP(inputBytes)

		err = json.Unmarshal(inputBytes[:length], &notification)
		if err != nil {
			c.logger.Error(err)
			os.Exit(1)
		}

		c.logger.Println("Notification received: ", notification)

		c.ReceivedNotifications <- notification
	}
}

func (c *client) NotifyNetwork(message *protocol.Notification) {
	c.notificationsToSend <- *message
}
