package main

import (
	"log"

	"github.com/teimurjan/go-p2p/notify"

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

	notificator := notify.NewNotificator(config.Port, logger)

	c := client.NewClient(notificator, logger)
	c.Start()
}
