/* Copyright (C) 2020 CompuZest, Inc. - All Rights Reserved
 *
 * Unauthorized copying of this file, via any medium, is strictly prohibited
 * Proprietary and confidential
 *
 * NOTICE: All information contained herein is, and remains the property of
 * CompuZest, Inc. The intellectual and technical concepts contained herein are
 * proprietary to CompuZest, Inc. and are protected by trade secret or copyright
 * law. Dissemination of this information or reproduction of this material is
 * strictly forbidden unless prior written permission is obtained from CompuZest, Inc.
 */

package file

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"

	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/ghodss/yaml"
	"k8s.io/apimachinery/pkg/util/json"
)

const (
	DefaultFilePermission = 0o600
)

//go:generate mockgen --build_flags=--mod=mod -destination=./mock_file.go -package=file "github.com/compuzest/zlifecycle-il-operator/controller/codegen/file" API
type API interface {
	SaveFileFromString(input string, folderName string, fileName string) error
	SaveFileFromByteArray(input []byte, folderName string, fileName string) error
	SaveYamlFile(obj interface{}, folderName string, fileName string) error
	SaveVarsToFile(variables []*stablev1.Variable, folderName string, fileName string) error
	CreateEmptyDirectory(folderName string) error
	SaveFileFromTemplate(t *template.Template, vars interface{}, folderName string, fileName string) error
	NewFile(folderName string, fileName string) (*os.File, error)
	RemoveAll(path string) error
	CopyDirContent(src string, dst string, mkdir bool) error
	CopyFile(src string, dst string) error
	IsDir(path string) bool
	IsFile(path string) bool
	FileExistsInDir(dir, path string) (bool, error)
	ReadDir(path string) ([]fs.DirEntry, error)
	CleanDir(path string, exclude []string) error
}

type OSFileService struct{}

var _ API = (*OSFileService)(nil)

func NewOSFileService() *OSFileService {
	return &OSFileService{}
}

func (f *OSFileService) CopyDirContent(src string, dst string, mkdir bool) error {
	if mkdir {
		if err := os.MkdirAll(dst, os.ModePerm); err != nil {
			return errors.Wrapf(err, "error creating folder [%s]", dst)
		}
	}
	if !f.IsDir(src) {
		return errors.New("source is not a directory")
	}

	files, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, file := range files {
		absoluteSrc := filepath.Join(src, file.Name())
		absoluteDst := filepath.Join(dst, file.Name())
		// copy subfolders
		if f.IsDir(absoluteSrc) {
			if err := f.CopyDirContent(absoluteSrc, absoluteDst, mkdir); err != nil {
				return err
			}
		}
		if err := f.CopyFile(absoluteSrc, absoluteDst); err != nil {
			return err
		}
	}

	return nil
}

func (f *OSFileService) FileExistsInDir(dir, path string) (bool, error) {
	joined := filepath.Join(dir, path)
	if _, err := os.Stat(joined); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (f *OSFileService) CopyFile(src string, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	safeDst := dst
	if f.IsDir(safeDst) {
		name := f.ExtractNameFromPath(src)
		safeDst = filepath.Join(dst, name)
	}

	err = os.WriteFile(safeDst, input, DefaultFilePermission)
	if err != nil {
		return err
	}

	return nil
}

func (f *OSFileService) IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

func (f *OSFileService) IsFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !fileInfo.IsDir()
}

func (f *OSFileService) ExtractNameFromPath(path string) string {
	tokens := strings.Split(path, "/")
	return tokens[len(tokens)-1]
}

func (f *OSFileService) NewFile(folderName string, fileName string) (*os.File, error) {
	file, err := createFileRecursive(folderName, fileName)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (f *OSFileService) SaveVarsToFile(variables []*stablev1.Variable, folderName string, fileName string) error {
	file, err := createFileRecursive(folderName, fileName)
	if err != nil {
		return err
	}
	defer closeFile(file)

	for _, variable := range variables {
		if _, err := fmt.Fprintf(file, "%s = \"%s\"\n", variable.Name, variable.Value); err != nil {
			return err
		}
	}

	return nil
}

// SaveFileFromTemplate creates a file with variables.
func (f *OSFileService) SaveFileFromTemplate(t *template.Template, vars interface{}, folderName string, fileName string) error {
	file, err := createFileRecursive(folderName, fileName)
	if err != nil {
		return err
	}
	defer closeFile(file)
	err = t.Execute(file, vars)

	if err != nil {
		return err
	}
	return nil
}

func closeFile(file *os.File) {
	_ = file.Close()
}

func createFileRecursive(folderName string, fileName string) (*os.File, error) {
	if err := os.MkdirAll(folderName, os.ModePerm); err != nil {
		return nil, fmt.Errorf("error: failed to create directory: %w", err)
	}
	file, err := os.Create(fmt.Sprintf("%s/%s", folderName, fileName))
	if err != nil {
		return nil, err
	}

	return file, nil
}

// SaveFileFromString Create file.
func (f *OSFileService) SaveFileFromString(input string, folderName string, fileName string) error {
	return saveBytesToFile([]byte(input), folderName, fileName)
}

func (f *OSFileService) SaveFileFromByteArray(input []byte, folderName string, fileName string) error {
	return saveBytesToFile(input, folderName, fileName)
}

// CreateEmptyDirectory creates empty directory with a .keep file.
func (f *OSFileService) CreateEmptyDirectory(folderName string) error {
	return saveBytesToFile(nil, folderName, ".keep")
}

// SaveYamlFile creates file and directory, does not validate yaml.
func (f *OSFileService) SaveYamlFile(obj interface{}, folderName string, fileName string) error {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("error: failed to marshal json: %w", err)
	}

	bytes, err := yaml.JSONToYAML(jsonBytes)
	if err != nil {
		return fmt.Errorf("error: failed to convert json to yaml: %w", err)
	}

	return saveBytesToFile(bytes, folderName, fileName)
}

func (f *OSFileService) ReadDir(path string) ([]fs.DirEntry, error) {
	return os.ReadDir(path)
}

// CleanDir cleans directory prior to generating files from controllers.
func (f *OSFileService) CleanDir(dir string, exclude []string) (err error) {
	if !f.IsDir(dir) {
		return nil
	}
	files, err := f.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		keep := false

		for _, k := range exclude {
			if file.Name() == k {
				keep = true
				break
			}
		}

		if keep {
			continue
		}

		if err := f.RemoveAll(filepath.Join(dir, file.Name())); err != nil {
			return err
		}
	}

	return nil
}

func saveBytesToFile(bytes []byte, folderName string, fileName string) error {
	if err := os.MkdirAll(folderName, os.ModePerm); err != nil {
		return errors.Wrapf(err, "error recursively creating directories for path [%s] with perms [%s]", folderName, os.ModePerm)
	}

	return os.WriteFile(fmt.Sprintf("%s/%s", folderName, fileName), bytes, 0o600)
}

// RemoveAll removes path and any children it contains.
// It removes everything it can but returns the first error
// it encounters. If the path does not exist, RemoveAll
// returns nil (no error).
func (f *OSFileService) RemoveAll(path string) error {
	return os.RemoveAll(path)
}
