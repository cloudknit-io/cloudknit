package main

import (
	"github.com/compuzest/zlifecycle-state-manager/app/web"
	"github.com/compuzest/zlifecycle-state-manager/app/zlog"
)

var Version = "v0.0.1"

func main() {
	zlog.PlainLogger().WithField("version", Version).Info("Starting zlifecycle state manager on localhost:8080")
	web.NewServer()
}
