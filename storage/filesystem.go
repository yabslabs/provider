package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const (
	readPermission = 0444
)

type FileSystem struct{}

func (fs *FileSystem) Save(path, filename string, data []byte) error {
	if err := os.Mkdir(path, os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}

	return ioutil.WriteFile(generateFilepath(path, filename), data, readPermission)
}

func generateFilepath(path, filename string) string {
	now := time.Now().Format("2006-01-02T150405")
	filename = fmt.Sprintf("%v.%v", filename, now)
	return fmt.Sprint(path, string(os.PathSeparator), filename)
}
