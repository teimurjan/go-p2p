package client

import (
	"github.com/sirupsen/logrus"

	"github.com/teimurjan/go-p2p/notify"
	"github.com/teimurjan/go-p2p/protocol"
	"github.com/teimurjan/go-p2p/udpAddrsArray"
)

// Client is a P2P client interface
type Client interface {
	Start()
}

type client struct {
	peers       udpAddrsArray.UDPAddrsArray
	notificator notify.Notificator
	logger      *logrus.Logger
}

// NewClient creates new client instance
func NewClient(notificator notify.Notificator, logger *logrus.Logger) Client {
	peers := udpAddrsArray.NewUDPAddrsArray()
	return &client{
		peers,
		notificator,
		logger,
	}
}

func (c *client) Start() {
	go c.notificator.StartNotifier()
	go c.notificator.StartNotificationListener()
	c.handleNotifications()
}

func (c *client) handleNotifications() {
	for {
		notification := <-c.notificator.GetReceivedNotifications()
		if notification.ID == protocol.ConnectedID {
			c.logger.Println("A new client is connected " + string(notification.FromAddr.IP))
			c.peers = c.peers.Add(notification.FromAddr)
		} else if notification.ID == protocol.DisconnectedID {
			c.logger.Println("The client is disconnected " + string(notification.FromAddr.IP))
			c.peers = c.peers.Remove(notification.FromAddr)
		}
	}
}
