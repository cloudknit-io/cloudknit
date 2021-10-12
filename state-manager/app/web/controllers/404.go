package controllers

import (
	http2 "github.com/compuzest/zlifecycle-state-manager/app/web/http"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http2.ErrorResponse(w, "endpoint not implemented", 404)
}
