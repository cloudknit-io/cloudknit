package common

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// CopyFile copies file from src path to dst path.
// If dst path is a folder, it will extract the filename of the src path and append it to dst.
func CopyFile(src string, dst string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	safeDst := dst
	isDir, err := IsDir(safeDst)
	if err != nil {
		return err
	}
	if isDir {
		name := ExtractNameFromPath(src)
		safeDst = filepath.Join(dst, name)
	}

	err = ioutil.WriteFile(safeDst, input, 0o644)
	if err != nil {
		return err
	}

	return nil
}

// CopyDirContent copies content of the src directory to the dst directory.
func CopyDirContent(src string, dst string) error {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, file := range files {
		fileSrc := filepath.Join(src, file.Name())
		fileDst := filepath.Join(dst, file.Name())
		if err := CopyFile(fileSrc, fileDst); err != nil {
			return err
		}
	}

	return nil
}

func IsDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}
