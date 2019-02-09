package app

import (
	"sync"

	"github.com/teimurjan/go-p2p/client"
	"github.com/teimurjan/go-p2p/config"
	"github.com/teimurjan/go-p2p/imstorage"
	"github.com/teimurjan/go-p2p/logging"
	"github.com/teimurjan/go-p2p/models"
	"github.com/teimurjan/go-p2p/notify"
	"github.com/teimurjan/go-p2p/protocol"
	"github.com/teimurjan/go-p2p/server"
)

// App is the whole app interface
type App interface {
	Start()
}

type application struct {
	config *config.Config
}

// NewApp creates new application instance
func NewApp(config *config.Config) App {
	return &application{config}
}

func (app *application) Start() {
	logger := logging.NewLogger(app.config)

	storage := app.createImstorage()

	n := notify.NewNotifier(app.config.Port, storage, logger)
	n.Start()

	c := client.NewClient(storage, logger)
	c.Start()

	app.greetPeers(storage)

	s := server.NewServer(app.config.Port, logger)
	s.Start()
}

func (app *application) createImstorage() imstorage.Storage {
	var wg sync.WaitGroup
	wg.Add(2)

	storage := imstorage.NewRedisStorage(app.config.RedisUrl)
	go storage.SubscribeToNotificationsToSend(func() { wg.Done() })
	go storage.SubscribeToNotificationsToHandle(func() { wg.Done() })

	wg.Wait()

	return storage
}

func (app *application) greetPeers(storage imstorage.Storage) {
	storage.AddNotificationToSend(
		&models.Notification{
			Req: &protocol.Request{Code: protocol.NewPeerCode},
		},
	)
}
