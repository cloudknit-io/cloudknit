package util

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

// DirExists returns whether the given file or directory exists
func DirExists(dir string) (bool, error) {
	_, err := os.Stat(dir)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func IsDirEmpty(dir string) (bool, error) {
	f, err := os.Open(dir)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if errors.Is(err, io.EOF) {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func ReadBody(stream io.ReadCloser) ([]byte, error) {
	body, err := io.ReadAll(stream)
	if err != nil {
		return nil, err
	}

	return body, nil
}
