package config_value

import (
	"fmt"
	"github.com/lie-flat-planet/service-init-tool/config_source"
)

func MergeConfigValue(svcFlattenedEnvKey map[string]struct{}, sources ...config_source.ISource) (structuralConfigInfo map[string]any, err error) {
	mergedValue := make(map[string]any)

	for _, src := range sources {
		var kv map[string]any
		kv, err = src.GetFlattenedConfigInfo()
		if err != nil {
			return nil, fmt.Errorf("get flattened config info error.err:%w", err)
		}

		for k, v := range kv {
			if _, ok := svcFlattenedEnvKey[k]; !ok {
				continue
			}

			mergedValue[k] = v
		}
	}

	return ParseFlattenedMap(mergedValue), nil
}
