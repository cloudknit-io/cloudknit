package util

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// CopyFile copies file from src path to dst path.
// If dst path is a folder, it will extract the filename of the src path and append it to dst.
func CopyFile(src string, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	safeDst := dst
	if IsDir(safeDst) {
		name := ExtractNameFromPath(src)
		safeDst = filepath.Join(dst, name)
	}

	err = os.WriteFile(safeDst, input, 0o600)
	if err != nil {
		return err
	}

	return nil
}

// CopyDirContent copies content of the src directory to the dst directory.
func CopyDirContent(src string, dst string) error {
	if !IsDir(src) {
		return errors.New("source is not a directory")
	}

	files, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, file := range files {
		absoluteSrc := filepath.Join(src, file.Name())
		// skip subfolders
		if IsDir(absoluteSrc) {
			continue
		}
		absoluteDst := filepath.Join(dst, file.Name())
		if err := CopyFile(absoluteSrc, absoluteDst); err != nil {
			return err
		}
	}

	return nil
}

func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

func IsFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !fileInfo.IsDir()
}
