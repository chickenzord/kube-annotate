package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	//LogLevel app-wide logger level
	LogLevel logrus.Level

	//LogFormat app-wide log formatter
	LogFormat logrus.Formatter

	//AppLogger app-wide logger
	AppLogger *logrus.Logger
)

func init() {
	LogLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		LogLevel = logrus.InfoLevel
	}

	if format := os.Getenv("LOG_FORMAT"); format == "json" {
		LogFormat = &logrus.JSONFormatter{}
	} else {
		LogFormat = &logrus.TextFormatter{}
	}

	AppLogger = logrus.New()
	AppLogger.SetLevel(LogLevel)
	AppLogger.SetFormatter(LogFormat)
}
