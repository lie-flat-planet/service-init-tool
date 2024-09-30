package util

import (
	"os"
	"path/filepath"
)

func CreateFile(filename string, content []byte) error {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	return os.WriteFile(filename, content, os.ModePerm)
}
