package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "embed"

	"github.com/cloudknit-io/cloudknit/event-service/internal/db"
	"github.com/cloudknit-io/cloudknit/event-service/internal/services"
	"github.com/cloudknit-io/cloudknit/event-service/internal/web"
	"github.com/cloudknit-io/cloudknit/event-service/internal/zlog"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

var (
	//go:embed .version
	Version    string
	baseLogger = zlog.NewPlainEntry()
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	baseLogger = baseLogger.WithField("version", Version)

	if len(os.Args) > 1 {
		if err := migrate(baseLogger); err != nil {
			panic(err)
		}
		return
	}

	errWg, errCtx := errgroup.WithContext(ctx)
	l := baseLogger.WithContext(ctx)

	l.Info("Initializing zlifecycle-event-service...")
	svcs, err := services.NewServices(l)
	if err != nil {
		l.Fatalf("error instantiating services: %v", err)
	}

	ss, err := web.NewStreamingServer(svcs, l)
	if err != nil {
		l.Fatalf("error creating streaming server: %v", err)
	}
	errWg.Go(func() error {
		if err := ss.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return errors.Wrap(err, "streaming server error")
		}
		return nil
	})

	rs, err := web.NewServer(svcs, l)
	if err != nil {
		l.Fatalf("error creating rest server: %v", err)
	}
	errWg.Go(func() error {
		if err := rs.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return errors.Wrap(err, "rest server error")
		}
		return nil
	})
	errWg.Go(func() error {
		<-errCtx.Done()
		stop()
		l.Info("Shutting down streaming server")
		if err := ss.Shutdown(errCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		l.Info("Shutting down REST server")
		if err := rs.Shutdown(errCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	err = errWg.Wait()
	if err == nil || errors.Is(err, context.Canceled) {
		l.Info("Streaming and REST server gracefully shutdown")
	} else if err != nil {
		l.Fatalf("unknown error occurred: %v", err)
	}
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

func setVersion() error {
	v, err := os.ReadFile(".version")
	if err != nil {
		return errors.Wrap(err, "error reading .version file")
	}
	Version = string(v)
	return nil
}
