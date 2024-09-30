package envvar

import (
	"fmt"
	"github.com/lie-flat-planet/service-init-tool/config_source"
	"github.com/lie-flat-planet/service-init-tool/util"
)

type Merger struct {
	envVarKeys map[string]struct{}
	sources    []config_source.ISource
}

func NewMerger(envVarKeys map[string]struct{}, sources ...config_source.ISource) *Merger {
	return &Merger{
		envVarKeys: envVarKeys,
		sources:    sources,
	}
}

func (m *Merger) Action() (structuralConfigInfo map[string]any, err error) {
	mergedValue := make(map[string]any)

	for _, src := range m.sources {
		if src == nil {
			continue
		}

		var kv map[string]any
		kv, err = src.GetFlattenedConfigInfo()
		if err != nil {
			return nil, fmt.Errorf("get flattened config info error.err:%w", err)
		}

		for k, v := range kv {
			if _, ok := m.envVarKeys[k]; !ok {
				continue
			}

			mergedValue[k] = v
		}
	}

	return util.ParseFlattenedMap(mergedValue), nil
}
