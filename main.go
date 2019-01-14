package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/teimurjan/go-p2p/client"
	"github.com/teimurjan/go-p2p/config"
	"github.com/teimurjan/go-p2p/logging"
	"github.com/teimurjan/go-p2p/protocol"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger := logging.NewLogger(config)

	c := client.NewClient(config.Port, logger)
	go c.StartNotifier()
	go c.StartNotificationListener()

	for {
		notification := <-c.GetReceivedNotifications()
		if notification.ID == protocol.ConnectedID {
			logger.Println("A new client is connected")
		}
	}
}
