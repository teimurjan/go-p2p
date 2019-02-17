package logging

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/teimurjan/go-p2p/config"
)

// NewLogger creates a new instance of logrus.Logger and set it up
func NewLogger(c *config.Config) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	file, err := getLoggerFile(c.LogFile)
	if err != nil {
		logger.Info("Can't open log file. Using default stdout.")
	} else {
		logger.SetOutput(file)
	}

	return logger
}

func getLoggerFile(fileName string) (*os.File, error) {
	rootDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	filePath := rootDir + "/" + fileName

	return os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0755)
}
