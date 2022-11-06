package http

import (
	"encoding/json"
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	err := Error{Message: message}
	jsonResp, _ := json.Marshal(err)
	_, _ = w.Write(jsonResp)
}

func WriteResponse(w http.ResponseWriter, body any, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	jsonResp, _ := json.Marshal(body)
	_, _ = w.Write(jsonResp)
}
