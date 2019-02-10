package server

import (
	"fmt"
	"os"

	"github.com/teimurjan/go-p2p/config"
	"github.com/teimurjan/go-p2p/fileutils"
)

var configuration, _ = config.NewConfig()

func init() {
	fileSourceDir := configuration.FileSourceDir

	if !fileutils.IsFileExists(fileSourceDir) {
		err := os.Mkdir(fileSourceDir, os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fatal: Error was occurred while trying to create directory %v\n", err)
			os.Exit(1)
		}
	}
}
