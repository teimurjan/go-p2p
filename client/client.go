package client

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"

	"github.com/teimurjan/go-p2p/config"
	"github.com/teimurjan/go-p2p/protocol"
)

type Client interface {
	StartNotifier()
	StartNotificationListener()
	NotifyNetwork(message *protocol.Message)
}

type client struct {
	Port                  string
	notificationsReceived chan protocol.Message
	notificationsToSend   chan protocol.Message
}

func NewClient(port string) Client {
	notificationsReceived := make(chan protocol.Message, 10)
	notificationsToSend := make(chan protocol.Message, 10)
	return &client{
		port,
		notificationsReceived,
		notificationsToSend,
	}
}

func (c *client) StartNotifier() {
	destinationAddress, _ := net.ResolveUDPAddr("udp", config.BroadcastAddress+":"+c.Port)
	localAddress, _ := net.ResolveUDPAddr("udp", "")
	connection, err := net.DialUDP("udp", localAddress, destinationAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer connection.Close()

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	for {
		message := <-c.notificationsToSend
		encoder.Encode(message)
		connection.Write(buffer.Bytes())
		buffer.Reset()
	}

}

func (c *client) StartNotificationListener() {
	localAddress, _ := net.ResolveUDPAddr("udp", ":"+c.Port)
	connection, err := net.ListenUDP("udp", localAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error1: %v\n", err)
		os.Exit(1)
	}
	defer connection.Close()

	var message protocol.Message
	for {
		inputBytes := make([]byte, 4096)
		length, _, _ := connection.ReadFromUDP(inputBytes)
		buffer := bytes.NewBuffer(inputBytes[:length])
		decoder := gob.NewDecoder(buffer)
		decoder.Decode(&message)
		c.notificationsReceived <- message
		fmt.Println(message)
	}
}

func (c *client) NotifyNetwork(message *protocol.Message) {
	c.notificationsToSend <- *message
}
