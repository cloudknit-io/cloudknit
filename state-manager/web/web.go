package web

import (
	"fmt"
	"github.com/compuzest/zlifecycle-state-manager/web/controllers"
	"github.com/compuzest/zlifecycle-state-manager/web/middleware"
	"github.com/compuzest/zlifecycle-state-manager/zlog"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
)

const (
	port = 8080
)

func NewServer() {
	errorChain := alice.New(middleware.EnforceJSONHandler, middleware.LoggerHandler, middleware.RecoverHandler)

	r := mux.NewRouter()
	r.HandleFunc("/state", controllers.StateHandler)
	r.NotFoundHandler = http.HandlerFunc(controllers.NotFoundHandler)
	http.Handle("/", errorChain.Then(r))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		zlog.Logger.Fatalf("Error from webserver: %w", err)
	}
}
