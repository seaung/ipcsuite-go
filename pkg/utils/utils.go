package utils

import (
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	letterIdxBits = 6
	letterIdXMask = 1<<letterIdxBits - 1
	letterIdXMax  = 63 / letterIdxBits
)

func RandomString(src *rand.Rand, letter string, number int) string {
	rbytes := make([]byte, 1)

	for i, cache, remain := number-1, src.Int63(), letterIdXMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdXMax
		}

		if idx := int(cache & letterIdXMask); idx < len(letter) {
			rbytes[i] = letter[idx]
			i--
		}

		cache >>= letterIdxBits
		remain--
	}
	return string(rbytes)
}

func RandomLowercase(number int) string {
	lowercase := "abcdefghijklmnopqrstuvwxyz"
	src := rand.New(rand.NewSource(time.Now().Unix()))
	return RandomString(src, lowercase, number)
}

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
