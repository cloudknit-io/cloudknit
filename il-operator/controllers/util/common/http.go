package common

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-logr/logr"
)

func CloseBody(body io.ReadCloser) {
	if body != nil {
		_ = body.Close()
	}
}

func LogBody(log logr.Logger, body io.ReadCloser) {
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Error(err, "Error while deserializing body")
		return
	}
	bodyString := string(bodyBytes)
	log.Info(bodyString)
}

func GetHTTPClient() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	t.TLSHandshakeTimeout = 5 * time.Second

	return &http.Client{
		Timeout:   5 * time.Second,
		Transport: t,
	}
}
