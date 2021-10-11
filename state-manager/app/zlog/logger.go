package zlog

import (
	log "github.com/sirupsen/logrus"
)

var Logger = initLogger()

func initLogger() *log.Logger {
	l := log.New()
	l.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return l
}
