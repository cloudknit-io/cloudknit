package common

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
)

func Success() {
	os.Exit(0)
}

func Failure(exitCode int) {
	os.Exit(exitCode)
}

func HandleError(err error, exitCode int) {
	if err != nil {
		Failure(exitCode)
	}
}

func GetHTTPClient() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	t.TLSHandshakeTimeout = 5 * time.Second

	return &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}
}

func ReadBody(stream io.ReadCloser) ([]byte, error) {
	body, err := io.ReadAll(stream)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func ToJSON(data interface{}) ([]byte, error) {
	jsoned, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return jsoned, nil
}

func FromJSON(s interface{}, jsonData []byte) error {
	err := json.Unmarshal(jsonData, s)
	if err != nil {
		return err
	}

	return nil
}

func CloseBody(body io.ReadCloser) {
	if body != nil {
		_ = body.Close()
	}
}
