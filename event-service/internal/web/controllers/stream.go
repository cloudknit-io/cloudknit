package controllers

import (
	"fmt"
	"net/http"

	"github.com/compuzest/zlifecycle-event-service/internal/services"
	"github.com/compuzest/zlifecycle-event-service/internal/stream"
	"github.com/compuzest/zlifecycle-event-service/internal/zlog"
	"github.com/pkg/errors"
)

func StreamHandler(svcs *services.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := zlog.NewCtxEntry(r.Context())
		// Make sure that the writer supports flushing.
		//
		flusher, ok := w.(http.Flusher)
		if !ok {
			err := errors.New("connection does not support streaming")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			l.Errorf("%+v", err)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Each connection registers its own message channel with the Broker's connections registry
		messageChan := svcs.SSEBroker.NewMessageChannel()

		// Signal the broker that we have a new connection
		svcs.SSEBroker.NewConnection(messageChan)

		// Remove this client from the map of connected clients
		// when this handler exits.
		defer func() {
			svcs.SSEBroker.CloseConnection(messageChan)
		}()

		// Listen to connection close and un-register messageChan
		// notify := rw.(http.CloseNotifier).CloseNotify()
		notify := r.Context().Done()
		go func() {
			<-notify
			svcs.SSEBroker.CloseConnection(messageChan)
		}()

		for {
			// Write to the ResponseWriter
			// Server Sent Events compatible
			data := <-messageChan
			_, err := fmt.Fprintf(w, stream.SSEDataFormat, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				l.Errorf("%+v", err)
			}

			// Flush the data immediately instead of buffering it for later.
			flusher.Flush()
		}
	}
}
