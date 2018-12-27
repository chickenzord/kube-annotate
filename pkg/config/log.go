package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

//CustomLogger logrus logger with helper methods
type CustomLogger struct{ *logrus.Logger }

var (
	//LogLevel app-wide logger level
	LogLevel logrus.Level

	//LogFormat app-wide log formatter
	LogFormat logrus.Formatter

	//AppLogger app-wide logger
	AppLogger *CustomLogger
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

	AppLogger = &CustomLogger{logrus.New()}
	AppLogger.SetLevel(LogLevel)
	AppLogger.SetFormatter(LogFormat)
}

//WithData embed data field
func (l *CustomLogger) WithData(data interface{}) *logrus.Entry {
	return l.WithField("data", data)
}
