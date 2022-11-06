package util

import (
	"io"
)

func CloseBody(body io.ReadCloser) {
	if body != nil {
		_ = body.Close()
	}
}
