package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/teimurjan/go-p2p/client"
	"github.com/teimurjan/go-p2p/config"
	"github.com/teimurjan/go-p2p/logging"
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
		<-c.ReceivedNotifications

	}
}
