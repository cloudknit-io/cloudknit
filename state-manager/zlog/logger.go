package zlog

import (
	log "github.com/sirupsen/logrus"
)

var Logger = initLogger()

func initLogger() *log.Logger {
	l := log.New()
	l.SetFormatter(&log.TextFormatter{})

	return l
}
