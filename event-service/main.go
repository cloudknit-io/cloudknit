package main

import (
	_ "embed"
	"github.com/compuzest/zlifecycle-event-service/app/services"
	"github.com/compuzest/zlifecycle-event-service/app/web"
	"github.com/compuzest/zlifecycle-event-service/app/zlog"
	_ "github.com/go-sql-driver/mysql"
)

//go:embed .version
var Version string

func main() {
	svcs, err := services.NewServices()
	if err != nil {
		zlog.PlainLogger().Fatal("error instantiating services", "cause", err)
	}
	zlog.PlainLogger().WithField("version", Version).Info("Starting zlifecycle event service on localhost:8080")
	web.NewServer(svcs)
}
