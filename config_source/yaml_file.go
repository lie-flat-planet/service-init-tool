package config_source

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type YamlFile struct {
	path    string
	content map[string]any
}

func NewYamlFile(path string) *YamlFile {
	return &YamlFile{
		path:    path,
		content: make(map[string]any),
	}
}

func (file *YamlFile) GetFlattenedConfigInfo() (map[string]any, error) {
	if err := file.parse(); err != nil {
		return nil, err
	}

	return file.content, nil
}

func (file *YamlFile) parse() error {
	fileContents, err := os.ReadFile(file.path)
	if err != nil {
		return err
	}
	if len(fileContents) < 1 {
		return fmt.Errorf("the yaml config file's content can't be empty")
	}

	return yaml.Unmarshal(fileContents, &file.content)
}
