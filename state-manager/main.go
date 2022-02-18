package main

import (
	"github.com/compuzest/zlifecycle-state-manager/app/web"
	"github.com/compuzest/zlifecycle-state-manager/app/zlog"
)

var Version = "0.0.3"

func main() {
	zlog.PlainLogger().WithField("version", Version).Info("Starting zlifecycle state manager on localhost:8080")
	web.NewServer()
}
