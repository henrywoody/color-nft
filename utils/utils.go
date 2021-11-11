package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func RecursiveMkDirIfNotExists(dirPath string) error {
	dirs := strings.Split(dirPath, string(filepath.Separator))
	currentDir := ""
	for _, dir := range dirs {
		currentDir = filepath.Join(currentDir, dir)
		if err := MkDirIfNotExists(currentDir); err != nil {
			return fmt.Errorf("error creating dir '%s': %v", currentDir, err)
		}
	}
	return nil
}

func MkDirIfNotExists(dirPath string) error {
	exists, err := checkDirExists(dirPath)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return os.Mkdir(dirPath, 0777)
}

func checkDirExists(dirPath string) (bool, error) {
	_, err := os.Stat(dirPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
