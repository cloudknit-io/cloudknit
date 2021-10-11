package controllers

import (
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "endpoint not implemented", 404)
}
