package server

import (
	"fmt"
	"os"

	"github.com/teimurjan/go-p2p/config"
)

var configuration, _ = config.NewConfig()

func init() {
	fileSourceDir := configuration.FileSourceDir
	if _, error := os.Stat(fileSourceDir); os.IsNotExist(error) {
		error = os.Mkdir(fileSourceDir, os.ModePerm)
		if error != nil {
			fmt.Fprintf(os.Stderr, "fatal: Error was occurred while trying to create directory %v\n", error)
			os.Exit(1)
		}
	}
}
