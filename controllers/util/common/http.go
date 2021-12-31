package common

import (
	"bytes"
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
	bodyBytes, err := io.ReadAll(body)
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
		Timeout:   10 * time.Second,
		Transport: t,
	}
}

func ReadBody(stream io.ReadCloser) ([]byte, error) {
	body, err := ioutil.ReadAll(stream)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func CreateMockResponse(code int) *http.Response {
	r := http.Response{Body: ioutil.NopCloser(bytes.NewReader([]byte{})), StatusCode: code}
	return &r
}
