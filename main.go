package main

import (
	"github.com/teimurjan/go-p2p/client"
	"github.com/teimurjan/go-p2p/protocol"
)

func main() {
	c := client.NewClient("3000")
	go c.StartNotifier()
	go c.StartNotificationListener()
	c.NotifyNetwork(&protocol.Message{
		ID: "/connected",
	})
}
