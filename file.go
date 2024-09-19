package service_init_tool

import (
	"os"
	"path/filepath"
)

func GenerateDevYML() {

}

func createFile(filename string, content []byte) error {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	return os.WriteFile(filename, content, os.ModePerm)
}
