package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func GetFilenames(directory, extend string) []string {
	var files []string

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return nil
	}

	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if strings.HasPrefix(info.Name(), extend) {
				filename, _ := filepath.Abs(path)
				files = append(files, filename)
			}
		}
		return nil
	})
	return files
}

func IsFileExists(filename string) bool {
	fd, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !fd.IsDir()
}

func IsFolderExists(directory string) bool {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return false
	}
	return true
}
