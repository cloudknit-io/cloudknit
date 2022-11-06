package controllers

import (
	"net/http"

	http2 "github.com/compuzest/zlifecycle-event-service/internal/web/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http2.ErrorResponse(w, "endpoint not implemented", 404)
}
