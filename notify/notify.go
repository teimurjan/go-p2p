package notify

import (
	"encoding/json"
	"net"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/teimurjan/go-p2p/imstorage"
	"github.com/teimurjan/go-p2p/models"
	"github.com/teimurjan/go-p2p/protocol"
	"github.com/teimurjan/go-p2p/utils"
)

// Notifier is base notifier interface
type Notifier interface {
	Start()
}

type notifier struct {
	port    string
	storage imstorage.Storage
	logger  *logrus.Logger
}

// NewNotifier creates new notifier instance
func NewNotifier(port string, storage imstorage.Storage, logger *logrus.Logger) Notifier {
	return &notifier{
		port,
		storage,
		logger,
	}
}

func (n *notifier) Start() {
	go n.startNotifier()
	go n.startNotificationListener()
}

func (n *notifier) startNotifier() {
	destinationAddress, _ := net.ResolveUDPAddr("udp", "255.255.255.255:"+n.port)
	connection, err := net.DialUDP("udp", nil, destinationAddress)
	if err != nil {
		n.logger.Error(err)
		os.Exit(1)
	}
	defer connection.Close()

	n.logger.Println("Notifier has started.")

	for {
		notification := <-n.storage.GetNotificationsToSend()

		encodedRequest, err := json.Marshal(notification.Req)
		if err != nil {
			n.logger.Error(err)
			continue
		}

		_, err = connection.Write(encodedRequest)
		if err != nil {
			n.logger.Error(err)
		}
	}

}

func (n *notifier) startNotificationListener() {
	localAddress, _ := net.ResolveUDPAddr("udp", ":"+n.port)
	connection, err := net.ListenUDP("udp", localAddress)
	if err != nil {
		n.logger.Error(err)
		os.Exit(1)
	}
	defer connection.Close()

	n.logger.Println("Notification listener has started.")

	var notification models.Notification
	var request protocol.Request
	for {
		inputBytes := make([]byte, 4096)
		length, addr, err := connection.ReadFromUDP(inputBytes)
		if err != nil {
			n.logger.Error(err)
			continue
		}

		if addr.IP.Equal(utils.GetLocalIP()) {
			continue
		}

		err = json.Unmarshal(inputBytes[:length], &request)
		if err != nil {
			n.logger.Error(err)
			continue
		}

		notification.Req = &request
		notification.FromAddr = addr

		err = n.storage.AddNotificationToHandle(&notification)
		if err != nil {
			n.logger.Error(err)
		}
	}
}
