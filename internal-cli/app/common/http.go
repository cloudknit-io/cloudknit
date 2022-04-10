package common

import (
	"io"
	"net/http"
	"time"
)

func NewHTTPClient() *http.Client {
	transport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return http.DefaultClient
	}
	t := transport.Clone()
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

func CloseBody(body io.ReadCloser) {
	if body != nil {
		_ = body.Close()
	}
}
