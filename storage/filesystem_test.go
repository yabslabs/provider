package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

var storage Storage

func init() {
	storage = NewStorage()
}

func TestGenerateFilepath(t *testing.T) {
	data := []byte("hello yabslabs")

	testdir := generateBackupDir()
	err := storage.Save(testdir, "test.back", data)
	assert.NoError(t, err)
	files, err := ioutil.ReadDir(testdir)
	require.NoError(t, err)
	assert.Len(t, files, 1)
	assert.True(t, strings.HasPrefix(files[0].Name(), "test.back"), "filename is should start with test.back but was \"%v\"", files[0].Name)
	toRemove := fmt.Sprint(testdir, string(os.PathSeparator), files[0].Name())
	require.NoError(t, os.Remove(toRemove))
}

func generateBackupDir() string {
	dir := os.TempDir()
	if !strings.HasSuffix(dir, string(os.PathSeparator)) {
		dir += string(os.PathSeparator)
	}
	dir += "backups"
	return dir
}
