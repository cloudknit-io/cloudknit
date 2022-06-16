package web

import (
	"fmt"
	"net/http"

	"github.com/compuzest/zlifecycle-event-service/app/web/middleware"
	"github.com/justinas/alice"
	"github.com/sirupsen/logrus"

	"github.com/compuzest/zlifecycle-event-service/app/services"

	"github.com/pkg/errors"

	swaggermiddleware "github.com/go-openapi/runtime/middleware"

	"github.com/compuzest/zlifecycle-event-service/app/apm"
	"github.com/compuzest/zlifecycle-event-service/app/env"
	"github.com/compuzest/zlifecycle-event-service/app/web/controllers"
	"github.com/compuzest/zlifecycle-event-service/app/zlog"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"
)

//  Event Service API:
//    version: 0.0.1
//    title: Event Service API
//  Schemes: http, https
//  Host: localhost:8081
//  BasePath: /
//  Produces:
//    - application/json
//
// swagger:meta

const (
	apiPort       = 8081
	streamingPort = 8082
)

func NewStreamingServer(svcs *services.Services) {
	r, err := initStreamingRouter(svcs)
	if err != nil {
		zlog.PlainLogger().Fatalf("error initializing streaming router: %v", err)
	}

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", streamingPort),
		Handler: r,
	}

	zlog.PlainLogger().WithFields(logrus.Fields{"port": streamingPort}).Infof("Started zlifecycle-event-service streaming server on port %d", streamingPort)
	if err := s.ListenAndServe(); err != nil {
		zlog.PlainLogger().Fatalf("error from streaming server: %v", err)
	}
}

func initStreamingRouter(svcs *services.Services) (*mux.Router, error) {
	r := mux.NewRouter()

	if env.Config().DevMode != "true" && env.Config().EnableNewRelic == "true" {
		zlog.PlainLogger().Info("Initializing streaming router with APM")
		zlog.PlainLogger().Info("Initializing NewRelic APM")
		app, err := apm.Init()
		if err != nil {
			return nil, errors.Wrap(err, "error initializing New Relic APM")
		}
		r.Handle(newrelic.WrapHandle(app, "/", controllers.StreamHandler(svcs)))
	} else {
		zlog.PlainLogger().Info("Initializing streaming router without APM")
		r.Handle("/", controllers.StreamHandler(svcs))
	}

	return r, nil
}

func NewServer(svcs *services.Services) {
	errorChain := alice.New(
		middleware.TimeoutHandler,
		middleware.EnforceJSONHandler,
		middleware.LoggerHandler,
		middleware.RecoverHandler,
	)

	r, err := initRESTRouter(svcs)
	if err != nil {
		zlog.PlainLogger().Fatalf("error initializing rest router: %v", err)
	}

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", apiPort),
		Handler: errorChain.Then(r),
	}

	zlog.PlainLogger().WithFields(logrus.Fields{"port": apiPort}).Infof("Started zlifecycle-event-service REST server on port %d", apiPort)
	if err := s.ListenAndServe(); err != nil {
		zlog.PlainLogger().Fatalf("error from rest server: %v", err)
	}
}

func initRESTRouter(svcs *services.Services) (*mux.Router, error) {
	r := mux.NewRouter()

	r.Handle("/swagger.yml", http.FileServer(http.Dir("./")))

	opts := swaggermiddleware.SwaggerUIOpts{SpecURL: "/swagger.yml"}
	sh := swaggermiddleware.SwaggerUI(opts, nil)

	if env.Config().DevMode != "true" && env.Config().EnableNewRelic == "true" {
		zlog.PlainLogger().Info("Initializing REST router with APM")
		zlog.PlainLogger().Info("Initializing NewRelic APM")
		app, err := apm.Init()
		if err != nil {
			return nil, errors.Wrap(err, "error initializing New Relic APM")
		}
		r.HandleFunc(newrelic.WrapHandleFunc(app, "/events", controllers.EventsHandler(svcs)))
		r.HandleFunc(newrelic.WrapHandleFunc(app, "/status", controllers.StatusHandler(svcs)))
		r.HandleFunc(newrelic.WrapHandleFunc(app, "/health/liveness", controllers.HealthHandler(svcs, false)))
		r.HandleFunc(newrelic.WrapHandleFunc(app, "/health/readiness", controllers.HealthHandler(svcs, true)))
		// documentation for developers
		r.Handle(newrelic.WrapHandle(app, "/docs", sh))
	} else {
		zlog.PlainLogger().Info("Initializing REST router without APM")
		r.HandleFunc("/events", controllers.EventsHandler(svcs))
		r.HandleFunc("/status", controllers.StatusHandler(svcs))
		r.HandleFunc("/health/liveness", controllers.HealthHandler(svcs, false))
		r.HandleFunc("/health/readiness", controllers.HealthHandler(svcs, true))
		// documentation for developers
		r.Handle("/docs", sh)
	}

	r.NotFoundHandler = http.HandlerFunc(controllers.NotFoundHandler)

	return r, nil
}
