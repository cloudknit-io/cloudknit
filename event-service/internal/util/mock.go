package util

import (
	"bytes"
	"io"
	"net/http"
)

func CreateMockEmptyBody() io.ReadCloser {
	return io.NopCloser(bytes.NewReader([]byte{}))
}

func CreateMockBody(body interface{}) io.ReadCloser {
	return io.NopCloser(bytes.NewReader(ToJSONBytes(body, false)))
}

func CreateMockResponse(code int) *http.Response {
	return &http.Response{Body: CreateMockEmptyBody(), StatusCode: code}
}

func CreateMockResponseWithBody(code int, body interface{}) *http.Response {
	return &http.Response{Body: CreateMockBody(body), StatusCode: code}
}
