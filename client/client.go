package client

import (
	"net"
	"time"

	"github.com/teimurjan/go-p2p/imstorage"
	"github.com/teimurjan/go-p2p/protocol"

	"github.com/sirupsen/logrus"

	"github.com/teimurjan/go-p2p/utilTypes"
)

// Client is a P2P client interface
type Client interface {
	Start()
	DownloadFile(path string)
}

type client struct {
	peers   utilTypes.UDPAddrsArray
	storage imstorage.Storage
	logger  *logrus.Logger
}

// NewClient creates new client instance
func NewClient(storage imstorage.Storage, logger *logrus.Logger) Client {
	peers := utilTypes.NewUDPAddrsArray()
	return &client{
		peers,
		storage,
		logger,
	}
}

func (c *client) Start() {
	c.logger.Println("Client has started.")
	go c.storage.SubscribeToNotificationsToHandle()
	c.handleNotifications()
}

func (c *client) handleNotifications() {
	for {
		notification := <-c.storage.GetNotificationsToHandle()
		if notification.Req.Code == protocol.NewPeerCode {
			c.logger.Println("A new client is connected " + string(notification.FromAddr.IP))
			c.peers.Add(notification.FromAddr)
		} else if notification.Req.Code == protocol.ExitPeerCode {
			c.logger.Println("The client is disconnected " + string(notification.FromAddr.IP))
			c.peers.Remove(notification.FromAddr)
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
