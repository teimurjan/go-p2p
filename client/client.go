package client

import (
	"errors"

	"github.com/teimurjan/go-p2p/imstorage"
	"github.com/teimurjan/go-p2p/protocol"

	"github.com/sirupsen/logrus"

	"github.com/teimurjan/go-p2p/utilTypes"
)

// Client is a P2P client interface
type Client interface {
	Start()
	DownloadFile(path string) error
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

	go c.handleNotifications()
}

func (c *client) handleNotifications() {
	for {
		notification := <-c.storage.GetNotificationsToHandle()
		if notification.Req.Code == protocol.NewPeerCode {
			c.logger.Printf("A new client is connected(IP=%s)", notification.FromAddr.IP.String())
			c.peers.Add(notification.FromAddr)
		} else if notification.Req.Code == protocol.ExitPeerCode {
			c.logger.Printf("The client is disconnected(IP=%s)", notification.FromAddr.IP.String())
			c.peers.Remove(notification.FromAddr)
		}
	}
}

func (c *client) DownloadFile(path string) error {
	peersWithFile := utilTypes.NewUDPAddrsArray()

	for _, peer := range c.peers {
		request := &protocol.Request{Code: protocol.CheckFileCode}

		response, err := process(peer, request)

		if err != nil {
			c.logger.Error(err)
		} else {
			c.logger.Printf("Connection with %s established", peer.IP.String())
			if response.Status == protocol.FileExistStatus {
				peersWithFile.Add(peer)
			}
		}

	}

	activePeersCount := int64(len(peersWithFile))

	if activePeersCount < 1 {
		return errors.New("No peers with file available")
	}

	return nil
}
