package app

import (
	"sync"

	"github.com/sirupsen/logrus"
	clientt "github.com/teimurjan/go-p2p/client"
	"github.com/teimurjan/go-p2p/config"
	"github.com/teimurjan/go-p2p/imstorage"
	"github.com/teimurjan/go-p2p/logging"
	"github.com/teimurjan/go-p2p/notify"
	serverr "github.com/teimurjan/go-p2p/server"
)

// App is the whole app interface
type App interface {
	Start()
}

type application struct {
	config   *config.Config
	logger   *logrus.Logger
	storage  imstorage.Storage
	notifier notify.Notifier
	client   clientt.Client
	server   serverr.Server
}

// NewApp creates new application instance
func NewApp(config *config.Config) App {
	logger := logging.NewLogger(config)
	storage := imstorage.NewRedisStorage(config.RedisURL)
	notifier := notify.NewNotifier(config.UDPPort, storage, logger)
	client := clientt.NewClient(config.HTTPPort, config.TCPPort, config.FileSourceDir, storage, logger)
	server := serverr.NewServer(config.TCPPort, config.FileSourceDir, storage, logger)

	return &application{
		config,
		logger,
		storage,
		notifier,
		client,
		server,
	}
}

func (app *application) Start() {
	app.startImstorage()

	app.notifier.Start()

	app.client.Start()

	app.server.Start()
}

func (app *application) startImstorage() {
	var wg sync.WaitGroup
	wg.Add(2)

	go app.storage.SubscribeToNotificationsToSend(func() { wg.Done() })
	go app.storage.SubscribeToNotificationsToHandle(func() { wg.Done() })

	wg.Wait()
}
