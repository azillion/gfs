package util

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// FileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// SaveFile creates a new file and saves the data to it
func SaveFile(fileName string, data []byte) error {
	if FileExists(fileName) {
		return fmt.Errorf("file already exists")
	}
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	err = f.Sync()
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}
	logrus.Debugf("saved %s\n", f.Name())
	return nil
}
