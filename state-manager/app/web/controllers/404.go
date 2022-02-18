package controllers

import (
	"net/http"

	http2 "github.com/compuzest/zlifecycle-state-manager/app/web/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http2.ErrorResponse(w, "endpoint not implemented", 404)
}
