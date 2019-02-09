package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/teimurjan/go-p2p/app"
	"github.com/teimurjan/go-p2p/config"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	a := app.NewApp(config)
	a.Start()
}
