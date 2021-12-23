package http

import (
	"fmt"
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

func NewVerboseError(class string, method string, endpoint string, err error) *VerboseError {
	return &VerboseError{
		Class:         class,
		Method:        method,
		Endpoint:      endpoint,
		Message:       err.Error(),
		OriginalError: err,
	}
}

func (v *VerboseError) Error() string {
	return fmt.Sprintf("%s %s: %s", v.Method, v.Endpoint, v.Message)
}
