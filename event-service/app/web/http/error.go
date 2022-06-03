package http

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
}

type VerboseError struct {
	Class         string
	Message       string
	Method        string
	Endpoint      string
	OriginalError error
}

func NewVerboseError(class string, r *http.Request, err error) *VerboseError {
	return &VerboseError{
		Class:         class,
		Method:        r.Method,
		Endpoint:      r.URL.Path,
		Message:       err.Error(),
		OriginalError: err,
	}
}

func (v *VerboseError) Error() string {
	return fmt.Sprintf("%s %s: %s", v.Method, v.Endpoint, v.Message)
}

func NewNotFoundError(r *http.Request) *VerboseError {
	return NewVerboseError("NotFoundError", r, errors.New("endpoint not implemented"))
}

func WriteNotFoundError(err error, w http.ResponseWriter, log *logrus.Entry) {
	log.Error(err)
	ErrorResponse(w, err.Error(), http.StatusNotFound)
}

func WriteInternalError(w http.ResponseWriter, verr *VerboseError, r *http.Request, log *logrus.Entry) {
	log.WithError(verr).Errorf("%s handler error", r.URL.Path)
	log.Errorf("%+v", verr.OriginalError)
	ErrorResponse(w, verr.Error(), http.StatusBadRequest)
}
