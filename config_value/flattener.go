package config_value

import (
	"reflect"
	"strings"
)

// FlattenMap 把有结构的yaml内容铺平
func FlattenMap(prefix string, input map[string]any, result map[string]any) {
	for k, v := range input {
		newKey := prefix + k
		val := reflect.ValueOf(v)
		switch val.Kind() {
		case reflect.Map:
			FlattenMap(newKey+".", v.(map[string]any), result)
		default:
			result[newKey] = v
		}
	}
}

// ParseFlattenedMap 将铺平的map解析成有结构的
func ParseFlattenedMap(flattened map[string]any) map[string]any {
	nested := make(map[string]any)
	for key, value := range flattened {
		keys := strings.Split(key, "_")
		setNestedMap(nested, keys, value)
	}
	return nested
}

func setNestedMap(data map[string]any, keys []string, value any) {
	if len(keys) == 1 {
		data[keys[0]] = value
		return
	}

	if _, ok := data[keys[0]]; !ok {
		data[keys[0]] = make(map[string]any)
	}

	nextMap, ok := data[keys[0]].(map[string]any)
	if !ok {
		nextMap = make(map[string]any)
		data[keys[0]] = nextMap
	}

	setNestedMap(nextMap, keys[1:], value)
}
