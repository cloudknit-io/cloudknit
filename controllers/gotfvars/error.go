package gotfvars

import "fmt"

type FileNotExists struct {
	Name string
}

func (m *FileNotExists) Error() string {
	return fmt.Sprintf("file does not exist: %s", m.Name)
}
