package util

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/google/go-github/v42/github"
)

func CloseBody(body io.ReadCloser) {
	if body != nil {
		_ = body.Close()
	}
}

func LogBody(log *logrus.Entry, body io.ReadCloser) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		log.WithError(err).Error("error deserializing response body")
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
	body, err := io.ReadAll(stream)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func CreateMockResponse(code int) *http.Response {
	return &http.Response{Body: io.NopCloser(bytes.NewReader([]byte{})), StatusCode: code}
}

func CreateMockGithubResponse(code int) *github.Response {
	return &github.Response{
		Response: CreateMockResponse(code),
	}
}
