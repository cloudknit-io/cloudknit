package main

import (
	"github.com/compuzest/zlifecycle-internal-cli/app/cmd"
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	log2 "github.com/compuzest/zlifecycle-internal-cli/app/log"
)

func main() {
	log := log2.NewLogger()
	log.WithField("version", env.Version).Info("Running zlifecycle-internal-cli")
	cmd.Execute()
}
