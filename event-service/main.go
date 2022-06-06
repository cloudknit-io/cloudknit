package main

import (
	_ "embed"
	"os"

	"github.com/compuzest/zlifecycle-event-service/app/db"
	"github.com/compuzest/zlifecycle-event-service/app/services"
	"github.com/compuzest/zlifecycle-event-service/app/web"
	"github.com/compuzest/zlifecycle-event-service/app/zlog"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//go:embed .version
var Version string

func main() {
	log := zlog.PlainLogger().WithField("version", Version)
	log.Info("Initializing zlifecycle-event-service...")
	svcs, err := services.NewServices()
	if err != nil {
		log.Fatal("error instantiating services", "cause", err)
	}
	defer svcs.SSEBroker.Close()

	if len(os.Args) > 1 {
		if err := migrate(log); err != nil {
			panic(err)
		}
		return
	}

	log.Info("Starting zlifecycle-event-service on localhost:8080")
	web.NewServer(svcs)
}

func migrate(log *logrus.Entry) error {
	if len(os.Args) != 3 {
		return errors.Errorf("invalid argument count for migration mode, expected 3, got %d", len(os.Args))
	}
	if os.Args[1] == "--migrate" || os.Args[1] == "-m" {
		action := os.Args[2]
		log.Infof("Executing migrate %s command", action)
		change, err := db.Migrate(db.MigrationCommand(action))
		if err != nil {
			return err
		}
		if change {
			log.Infof("Finished executing migrate %s command", action)
		} else {
			log.Info("Latest migrations already applied, no change needed")
		}

		return nil
	}
	return errors.Errorf("invalid argument provided, expected --migrate or -m, got %s", os.Args[1])
}
