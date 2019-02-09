package main

import (
	"log"

	"github.com/teimurjan/go-p2p/models"
	"github.com/teimurjan/go-p2p/notify"
	"github.com/teimurjan/go-p2p/protocol"
	"github.com/teimurjan/go-p2p/server"

	"github.com/teimurjan/go-p2p/imstorage"

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

	storage := imstorage.NewRedisStorage(config.RedisUrl)

	n := notify.NewNotifier(config.Port, storage, logger)
	n.Start()

	c := client.NewClient(storage, logger)
	c.Start()

	storage.AddNotificationToSend(
		&models.Notification{Req: &protocol.Request{Code: protocol.NewPeerCode}},
	)

	s := server.NewServer(config.Port, logger)
	s.Start()
}
