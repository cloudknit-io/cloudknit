package web

import (
	"fmt"
	"net/http"

	"github.com/cloudknit-io/cloudknit/event-service/internal/web/middleware"
	"github.com/justinas/alice"
	"github.com/sirupsen/logrus"

	"github.com/cloudknit-io/cloudknit/event-service/internal/services"

	"github.com/pkg/errors"

	swaggermiddleware "github.com/go-openapi/runtime/middleware"

	"github.com/cloudknit-io/cloudknit/event-service/internal/apm"
	"github.com/cloudknit-io/cloudknit/event-service/internal/env"
	"github.com/cloudknit-io/cloudknit/event-service/internal/web/controllers"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"
)

const (
	apiPort       = 8081
	streamingPort = 8082
)

func NewStreamingServer(svcs *services.Services, l *logrus.Entry) (*http.Server, error) {
	r, err := initStreamingRouter(svcs, l)
	if err != nil {
		return nil, errors.Wrap(err, "error initializing streaming router")
	}

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", streamingPort),
		Handler: r,
	}

	l.WithFields(logrus.Fields{"port": streamingPort}).Infof("Starting zlifecycle-event-service streaming server on port %d", streamingPort)
	return s, nil
}

func initStreamingRouter(svcs *services.Services, l *logrus.Entry) (*mux.Router, error) {
	r := mux.NewRouter()

	if env.Config().DevMode != "true" && env.Config().EnableNewRelic == "true" {
		l.Info("Initializing streaming router with APM")
		l.Info("Initializing NewRelic APM")
		app, err := apm.Init()
		if err != nil {
			return nil, errors.Wrap(err, "error initializing New Relic APM")
		}
		r.Handle(newrelic.WrapHandle(app, "/", controllers.StreamHandler(svcs)))
	} else {
		l.Info("Initializing streaming router without APM")
		r.Handle("/", controllers.StreamHandler(svcs))
	}

	return r, nil
}

func NewServer(svcs *services.Services, l *logrus.Entry) (*http.Server, error) {
	errorChain := alice.New(
		middleware.TimeoutHandler,
		middleware.EnforceJSONHandler,
		middleware.LoggerHandler,
		middleware.RecoverHandler,
	)

	r, err := initRESTRouter(svcs, l)
	if err != nil {
		return nil, errors.Wrap(err, "error initializing rest router")
	}

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", apiPort),
		Handler: errorChain.Then(r),
	}

	l.WithFields(logrus.Fields{"port": apiPort}).Infof("Starting zlifecycle-event-service REST server on port %d", apiPort)
	return s, nil
}

func initRESTRouter(svcs *services.Services, l *logrus.Entry) (*mux.Router, error) {
	r := mux.NewRouter()

	r.Handle("/swagger.yml", http.FileServer(http.Dir("./")))

	opts := swaggermiddleware.SwaggerUIOpts{SpecURL: "/swagger.yml"}
	sh := swaggermiddleware.SwaggerUI(opts, nil)

	if env.Config().DevMode != "true" && env.Config().EnableNewRelic == "true" {
		l.Info("Initializing REST router with APM")
		l.Info("Initializing NewRelic APM")
		app, err := apm.Init()
		if err != nil {
			return nil, errors.Wrap(err, "error initializing New Relic APM")
		}
		r.HandleFunc(newrelic.WrapHandleFunc(app, "/admin/db", controllers.AdminDatabaseHandler(svcs)))
		r.HandleFunc(newrelic.WrapHandleFunc(app, "/events", controllers.EventsHandler(svcs)))
		r.HandleFunc(newrelic.WrapHandleFunc(app, "/status", controllers.StatusHandler(svcs, l)))
		r.HandleFunc(newrelic.WrapHandleFunc(app, "/health/liveness", controllers.HealthHandler(svcs, false)))
		r.HandleFunc(newrelic.WrapHandleFunc(app, "/health/readiness", controllers.HealthHandler(svcs, true)))
		// documentation for developers
		r.Handle(newrelic.WrapHandle(app, "/docs", sh))
	} else {
		l.Info("Initializing REST router without APM")
		r.HandleFunc("/admin/db", controllers.AdminDatabaseHandler(svcs))
		r.HandleFunc("/events", controllers.EventsHandler(svcs))
		r.HandleFunc("/status", controllers.StatusHandler(svcs, l))
		r.HandleFunc("/health/liveness", controllers.HealthHandler(svcs, false))
		r.HandleFunc("/health/readiness", controllers.HealthHandler(svcs, true))
		// documentation for developers
		r.Handle("/docs", sh)
	}

	r.NotFoundHandler = http.HandlerFunc(controllers.NotFoundHandler)

	return r, nil
}
