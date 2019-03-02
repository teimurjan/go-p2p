package client

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/teimurjan/go-p2p/fileutils"
	"github.com/teimurjan/go-p2p/imstorage"
	"github.com/teimurjan/go-p2p/protocol"

	"github.com/sirupsen/logrus"

	"github.com/teimurjan/go-p2p/utilTypes"
)

type GUIRequestBody struct {
	path string
}

// ChunkSize is a default size chunk in bytes
const ChunkSize = 1024

// Client is a P2P client interface
type Client interface {
	Start()
	DownloadFile(path string) error
}

type client struct {
	GUIPort string
	peers   utilTypes.UDPAddrsArray
	storage imstorage.Storage
	logger  *logrus.Logger
}

// NewClient creates new client instance
func NewClient(GUIPort string, storage imstorage.Storage, logger *logrus.Logger) Client {
	peers := utilTypes.NewUDPAddrsArray()
	return &client{
		GUIPort,
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

func (c *client) listenGUI() {
	s := &http.Server{Addr: ":" + c.GUIPort}
	http.HandleFunc("/getFile", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			decoder := json.NewDecoder(r.Body)
			var body GUIRequestBody
			err := decoder.Decode(&body)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				c.logger.Error(err)
				return
			}

			c.DownloadFile(body.path)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("{\"msg\":\"Successfully downloaded\"}")
		}
	})

	s.ListenAndServe()
}

func (c *client) DownloadFile(path string) error {
	peersWithFile, err := c.getActivePeers()
	if err != nil {
		c.logger.Error(err)
		return err
	}

	activePeersCount := int64(len(peersWithFile))
	if activePeersCount < 1 {
		return errors.New("No peers with file available")
	}

	var fileInfo *protocol.ResponseInfo
	response, err := process(
		peersWithFile[0],
		&protocol.Request{Code: protocol.CheckFileCode},
	)
	if err != nil {
		c.logger.Error(err)
		return err
	}
	fileInfo = &response.Info

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

func (c *client) getActivePeers() (utilTypes.UDPAddrsArray, error) {
	peersWithFile := utilTypes.NewUDPAddrsArray()

	for _, peer := range c.peers {
		request := &protocol.Request{Code: protocol.CheckFileCode}

		response, err := process(peer, request)

		if err != nil {
			return nil, err
		} else if response.Status == protocol.FileExistStatus {
			peersWithFile.Add(peer)
		}

	}

	return peersWithFile, nil
}
