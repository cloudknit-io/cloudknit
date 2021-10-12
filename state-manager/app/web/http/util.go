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
	w.Write(jsonResp)
}

func Response(w http.ResponseWriter, body interface{}, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	jsonResp, _ := json.Marshal(body)
	w.Write(jsonResp)
}
