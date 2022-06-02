package web

import (
	"fmt"
	"net/http"

	"github.com/compuzest/zlifecycle-event-service/app/services"

	"github.com/pkg/errors"

	"github.com/compuzest/zlifecycle-event-service/app/apm"
	"github.com/compuzest/zlifecycle-event-service/app/env"
	"github.com/compuzest/zlifecycle-event-service/app/web/controllers"
	"github.com/compuzest/zlifecycle-event-service/app/web/middleware"
	"github.com/compuzest/zlifecycle-event-service/app/zlog"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

const (
	port = 8081
)

func NewServer(svcs *services.Services) {
	errorChain := alice.New(
		middleware.TimeoutHandler,
		middleware.EnforceJSONHandler,
		middleware.LoggerHandler,
		middleware.RecoverHandler,
	)

	r, err := initRouter(svcs)
	if err != nil {
		zlog.PlainLogger().Fatalf("Error initializing router: %v", err)
	}
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: errorChain.Then(r),
	}

	zlog.PlainLogger().WithFields(logrus.Fields{"port": port}).Info("Started HTTP server")
	if err := s.ListenAndServe(); err != nil {
		zlog.PlainLogger().Fatalf("Error from webserver: %v", err)
	}
}

func initRouter(svcs *services.Services) (*mux.Router, error) {
	r := mux.NewRouter()

	if env.Config().DevMode != "true" && env.Config().EnableNewRelic == "true" {
		zlog.PlainLogger().Info("Initializing application with APM")
		zlog.PlainLogger().Info("Initializing NewRelic APM")
		app, err := apm.Init()
		if err != nil {
			return nil, errors.Wrap(err, "error initializing New Relic APM")
		}
		r.HandleFunc(newrelic.WrapHandleFunc(app, "/events", controllers.EventsHandler(svcs)))
	} else {
		zlog.PlainLogger().Info("Initializing application without APM")
		r.HandleFunc("/events", controllers.EventsHandler(svcs))
	}

	r.NotFoundHandler = http.HandlerFunc(controllers.NotFoundHandler)

	return r, nil
}
