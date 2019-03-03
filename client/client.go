package client

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"path"
	"time"

	"github.com/teimurjan/go-p2p/fileutils"
	"github.com/teimurjan/go-p2p/imstorage"
	"github.com/teimurjan/go-p2p/protocol"
	"github.com/teimurjan/go-p2p/utils"

	"github.com/rs/cors"
	"github.com/sirupsen/logrus"

	"github.com/teimurjan/go-p2p/utilTypes"
)

// ChunkSize is a default size chunk in bytes
const ChunkSize int64 = 1024

// Client is a P2P client interface
type Client interface {
	Start()
	DownloadFile(pathToFile string) error
}

type client struct {
	GUIPort       string
	TCPPort       string
	fileSourceDir string
	peers         utilTypes.UDPAddrsArray
	storage       imstorage.Storage
	logger        *logrus.Logger
}

// NewClient creates new client instance
func NewClient(GUIPort string, TCPPort string, fileSourceDir string, storage imstorage.Storage, logger *logrus.Logger) Client {
	peers := utilTypes.NewUDPAddrsArray()
	return &client{
		GUIPort,
		TCPPort,
		fileSourceDir,
		peers,
		storage,
		logger,
	}
}

func (c *client) Start() {
	c.logger.Println("Client has started.")

	go c.handleNotifications()
	go c.listenGUI()
}

func (c *client) handleNotifications() {
	for {
		notification := <-c.storage.GetNotificationsToHandle()
		if notification.Req.Code == protocol.NewPeerCode {
			c.logger.Printf("A new client is connected(IP=%s)", notification.FromAddr.IP.String())
			c.peers = c.peers.Add(notification.FromAddr)
		} else if notification.Req.Code == protocol.ExitPeerCode {
			c.logger.Printf("The client is disconnected(IP=%s)", notification.FromAddr.IP.String())
			c.peers = c.peers.Remove(notification.FromAddr)
		}
	}
}

type requestBody struct {
	Path string `json:"path"`
}

func (c *client) listenGUI() {
	mux := http.NewServeMux()
	mux.HandleFunc("/getFile", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {

			b, err := ioutil.ReadAll(r.Body)

			defer r.Body.Close()

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				c.logger.Errorf("Cannot read request body. %e", err)
				return
			}

			var parsedBody requestBody
			err = json.Unmarshal(b, &parsedBody)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				c.logger.Errorf("Cannot marshal request body. %e", err)
				return
			}

			err = c.DownloadFile(parsedBody.Path)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")

			json.NewEncoder(w).Encode("{\"msg\":\"Successfully downloaded\"}")
		}
	})

	handler := cors.Default().Handler(mux)

	http.ListenAndServe(":"+c.GUIPort, handler)
}

func (c *client) DownloadFile(pathToFile string) error {
	peersWithFile, err := c.getPeersWithFile(pathToFile)
	if err != nil {
		c.logger.Error(err)
		return err
	}

	activePeersCount := int64(len(peersWithFile))
	if activePeersCount < 1 {
		return errors.New("No peers with file available")
	}

	var fileInfo *protocol.ResponseInfo
	response, err := c.doRequest(
		peersWithFile[0],
		&protocol.Request{
			Code: protocol.CheckFileCode,
			Info: protocol.RequestInfo{FileName: pathToFile},
		},
	)
	if err != nil {
		c.logger.Error(err)
		return err
	}
	fileInfo = &response.Info

	chunksCount := math.Ceil(float64(fileInfo.FileSize) / float64(ChunkSize))

	chunks := make([][]byte, int(chunksCount))

	for i, peer := range peersWithFile {
		request := &protocol.Request{
			Code: protocol.GetChunkCode,
			Info: protocol.RequestInfo{
				FileName:   fileInfo.FileName,
				ChunkIndex: int64(i),
				ChunkSize:  ChunkSize,
			},
		}

		response, err := c.doRequest(peer, request)
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
			fileData[(chunkIndex+1)*chunkPieceIndex] = chunkPiece
		}
	}

	fullPathToFile := path.Join(c.fileSourceDir, pathToFile)

	fileutils.SaveFile(fullPathToFile, fileData)

	return nil
}

func (c *client) getPeersWithFile(pathToFile string) (utilTypes.UDPAddrsArray, error) {
	peersWithFile := utilTypes.NewUDPAddrsArray()

	for _, peer := range c.peers {
		request := &protocol.Request{
			Code: protocol.CheckFileCode,
			Info: protocol.RequestInfo{FileName: pathToFile},
		}

		response, err := c.doRequest(peer, request)

		if err != nil {
			return nil, err
		} else if response.Status == protocol.FileExistStatus {
			peersWithFile = peersWithFile.Add(peer)
		}

	}

	return peersWithFile, nil
}

func (c *client) doRequest(peer *net.UDPAddr, request *protocol.Request) (*protocol.Response, error) {
	tcpAddr := peer.IP.String() + ":" + c.TCPPort
	conn, err := net.DialTimeout("tcp", tcpAddr, time.Second*2)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	json, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	conn.Write(json)

	response, err := utils.ReadResponseFromTCP(conn)
	if err != nil {
		return nil, err
	}

	return response, nil
}
