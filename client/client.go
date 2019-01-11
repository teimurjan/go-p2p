package client

import (
	"bytes"
	"encoding/gob"
	"net"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/teimurjan/go-p2p/protocol"
)

// Client is a P2P client interface
type Client interface {
	StartNotifier()
	StartNotificationListener()
	NotifyNetwork(message *protocol.Message)
}

type client struct {
	Port                  string
	logger                *logrus.Logger
	notificationsReceived chan protocol.Message
	notificationsToSend   chan protocol.Message
}

// NewClient creates new client instance
func NewClient(port string, logger *logrus.Logger) Client {
	notificationsReceived := make(chan protocol.Message, 10)
	notificationsToSend := make(chan protocol.Message, 10)
	return &client{
		port,
		logger,
		notificationsReceived,
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

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	for {
		message := <-c.notificationsToSend
		encoder.Encode(message)
		connection.Write(buffer.Bytes())
		c.logger.Println("Notification sent: ", message)
		buffer.Reset()
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

	var message protocol.Message
	for {
		inputBytes := make([]byte, 4096)
		length, _, _ := connection.ReadFromUDP(inputBytes)

		buffer := bytes.NewBuffer(inputBytes[:length])
		decoder := gob.NewDecoder(buffer)
		decoder.Decode(&message)

		c.logger.Println("Notification received: ", message)

		c.notificationsReceived <- message
	}
}

func (c *client) NotifyNetwork(message *protocol.Message) {
	c.notificationsToSend <- *message
}
