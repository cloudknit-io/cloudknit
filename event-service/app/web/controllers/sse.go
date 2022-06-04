package controllers

import (
	"net/http"

	"github.com/compuzest/zlifecycle-event-service/app/services"
	"github.com/compuzest/zlifecycle-event-service/app/util"
	http2 "github.com/compuzest/zlifecycle-event-service/app/web/http"
	"github.com/compuzest/zlifecycle-event-service/app/zlog"
)

func SSEHandler(svcs *services.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := zlog.CtxLogger(r.Context())
		log.Infof("Get handshake from client")

		prepareHeaderForSSE(w)

		flusher, ok := w.(http.Flusher)
		if !ok {
			http2.ErrorResponse(w, "connection does not support streaming", 500)
			return
		}
		// trap the request under loop forever
		for {
			select {
			case message := <-svcs.SSEBroker.Chan():
				_, err := w.Write(util.ToJSONBytes(message, false))
				if err != nil {
					log.Errorf("error writing to channel: %v", err)
					continue
				}
				flusher.Flush()

			// connection is closed then defer will be executed
			case <-r.Context().Done():
				return
			}
		}
	}
}

func prepareHeaderForSSE(w http.ResponseWriter) {
	// prepare the header
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
