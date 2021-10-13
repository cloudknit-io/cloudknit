package zlog

import (
	"github.com/newrelic/go-agent/v3/integrations/logcontext/nrlogrusplugin"
	log "github.com/sirupsen/logrus"
)

var Logger = initLogger()

func initLogger() *log.Logger {
	l := log.New()
	l.SetFormatter(nrlogrusplugin.ContextFormatter{})

	return l
}
