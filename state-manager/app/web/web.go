package web

import (
	"fmt"
	"net/http"

	"github.com/compuzest/zlifecycle-state-manager/app/apm"
	"github.com/compuzest/zlifecycle-state-manager/app/env"
	"github.com/compuzest/zlifecycle-state-manager/app/web/controllers"
	"github.com/compuzest/zlifecycle-state-manager/app/web/middleware"
	"github.com/compuzest/zlifecycle-state-manager/app/zlog"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

const (
	port = 8080
)

func NewServer() {
	errorChain := alice.New(
		middleware.TimeoutHandler,
		middleware.EnforceJSONHandler,
		middleware.LoggerHandler,
		middleware.RecoverHandler,
	)

	r, err := initRouter()
	if err != nil {
		zlog.PlainLogger().Fatalf("Error initializing router: %v", err)
	}
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: errorChain.Then(r),
	}

	zlog.PlainLogger().WithFields(logrus.Fields{"port": port}).Info("Starting HTTP server")
	if err := s.ListenAndServe(); err != nil {
		zlog.PlainLogger().Fatalf("Error from webserver: %v", err)
	}
}

func initRouter() (*mux.Router, error) {
	r := mux.NewRouter()

	if env.Config().DevMode != "true" && env.Config().EnableNewRelic == "true" {
		zlog.PlainLogger().Info("Initializing application with APM")
		zlog.PlainLogger().Info("Initializing NewRelic APM")
		app, err := apm.Init()
		if err != nil {
			return nil, err
		}
		r.HandleFunc(newrelic.WrapHandleFunc(app, "/state", controllers.StateHandler))
	} else {
		zlog.PlainLogger().Info("Initializing application without APM")
		r.HandleFunc("/state", controllers.StateHandler)
	}

	r.NotFoundHandler = http.HandlerFunc(controllers.NotFoundHandler)

	return r, nil
}
