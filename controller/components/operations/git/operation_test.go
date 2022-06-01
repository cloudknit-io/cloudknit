package git

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepRepo(t *testing.T) {
	t.Parallel()

	pattern := "operator-testing"

	tempDir, err := ioutil.TempDir("", pattern)
	assert.Nil(t, err)

	_, err = os.Create(path.Join(tempDir, "main.tf"))
	assert.Nil(t, err)

	err = os.Mkdir(path.Join(tempDir, ".git"), 0750)
	assert.Nil(t, err)

	err = os.Mkdir(path.Join(tempDir, "modules"), 0750)
	assert.Nil(t, err)

	_, err = os.Create(path.Join(tempDir, "modules", "main.tf"))
	assert.Nil(t, err)

	err = PrepRepo(tempDir)

	assert.Nil(t, err)

	files, err := os.ReadDir(tempDir)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(files))
}
