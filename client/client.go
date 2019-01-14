package client

import (
	"net"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/teimurjan/go-p2p/notify"
	"github.com/teimurjan/go-p2p/protocol"
	"github.com/teimurjan/go-p2p/udpAddrsArray"
)

// Client is a P2P client interface
type Client interface {
	Start()
	DownloadFile(path string)
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

func (c *client) DownloadFile(path string) {
	for _, peer := range c.peers {
		tcpAddr := peer.String() + ":" + string(peer.Port)
		_, err := net.DialTimeout("tcp", tcpAddr, time.Second*2)
		if err != nil {
			c.logger.Printf("Connection with %s cannot be established. Removing from the peers list.", tcpAddr)
			c.peers.Remove(peer)
			return
		}
		c.logger.Printf("Connection with %s established.", tcpAddr)
	}
}
