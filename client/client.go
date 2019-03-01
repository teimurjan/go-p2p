package client

import (
	"errors"

	"github.com/teimurjan/go-p2p/fileutils"
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
	peersWithFile, err := getActivePeers(c.peers)
	if err != nil {
		c.logger.Error(err)
		return err
	}

	activePeersCount := int64(len(peersWithFile))
	if activePeersCount < 1 {
		return errors.New("No peers with file available")
	}

	fileInfo := getFileInfo(peersWithFile[0])

	chunksCount := fileInfo.FileSize / ChunkSize
	chunks := make([][]byte, 0, chunksCount)

	for i, peer := range peersWithFile {
		request := &protocol.Request{
			Code: protocol.GetChunkCode,
			Info: protocol.RequestInfo{
				FileName:   fileInfo.FileName,
				ChunkIndex: int64(i),
				ChunkSize:  ChunkSize,
			},
		}

		response, err := process(peer, request)
		if err != nil {
			c.logger.Error(err)
			return err
		}

		if response.Status == protocol.ChunkSentStatus {
			chunks[i] = response.Data
		}
	}

	fileData := make([]byte, fileInfo.FileSize)
	for chunkIndex, chunk := range chunks {
		for chunkPieceIndex, chunkPiece := range chunk {
			fileData[chunkIndex*chunkPieceIndex] = chunkPiece
		}
	}

	fileutils.SaveFile(path, fileData)

	return nil
}
