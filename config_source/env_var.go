package config_source

import (
	"fmt"
	"os"
	"strings"
)

type EnvVar struct {
	content map[string]any
}

func NewEnvVar() *EnvVar {
	return &EnvVar{
		content: make(map[string]any),
	}
}

func (ev *EnvVar) GetFlattenedConfigInfo() (map[string]any, error) {
	if err := ev.parse(); err != nil {
		return nil, err
	}

	return ev.content, nil
}

func (ev *EnvVar) parse() error {
	osEnvVars := os.Environ()
	for _, kv := range osEnvVars {
		kvPair := strings.SplitN(kv, "=", 2)
		if len(kvPair) != 2 {
			return fmt.Errorf("the kv format of env var is invalid. it should be like 'a=b'")
		}

		k := kvPair[0]
		v := kvPair[1]

		ev.content[k] = v
	}

	return nil
}
