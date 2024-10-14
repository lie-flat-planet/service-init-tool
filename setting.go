package service_init_tool

import (
	"fmt"
	"reflect"
	"strconv"
)

// 设置字段
func setSettingFields(setting any, configValues map[string]any) {
	structValue := reflect.ValueOf(setting).Elem()
	structType := structValue.Type()

	for i := 0; i < structValue.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)

		if field.Anonymous {
			setSettingFields(fieldValue.Addr().Interface(), configValues)
			continue
		}

		_, tagOK := field.Tag.Lookup("env")
		if fieldValue.Kind() != reflect.Ptr && fieldValue.Kind() != reflect.Struct && !tagOK {
			continue
		}

		fieldName := field.Name

		if val, ok := configValues[fieldName]; ok {
			v := reflect.ValueOf(val)

			switch fieldValue.Kind() {
			case reflect.Ptr:
				if fieldValue.IsNil() {
					fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
				}
				if mapVal, ok := val.(map[string]any); ok {
					strMapVal := make(map[string]any)
					for k, v := range mapVal {
						strMapVal[k] = v
					}
					setSettingFields(fieldValue.Interface(), strMapVal)
				}
			case reflect.Struct:
				if mapVal, ok := val.(map[string]any); ok {
					strMapVal := make(map[string]any)
					for k, v := range mapVal {
						strMapVal[k] = v
					}
					setSettingFields(fieldValue.Addr().Interface(), strMapVal)
				}
			case reflect.Slice, reflect.Array:
				panic(fmt.Sprintf(`field "%s" don't use array or slice type. please use string'`, fieldName))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if v.Type().Kind() == reflect.String {
					res, err := strconv.ParseUint(v.String(), 10, 64)
					if err != nil {
						panic(err)
					}
					fieldValue.SetUint(res)
					continue
				}
				fieldValue.SetUint(uint64(v.Int()))
			case reflect.Float64, reflect.Float32:
				if v.Type().Kind() == reflect.String {
					res, err := strconv.ParseFloat(v.String(), 64)
					if err != nil {
						panic(err)
					}

					fieldValue.SetFloat(res)
					continue
				}
				fieldValue.SetFloat(v.Float())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if v.Type().Kind() == reflect.String {
					res, err := strconv.ParseInt(v.String(), 10, 64)
					if err != nil {
						panic(err)
					}

					fieldValue.SetInt(res)
					continue
				}
				fieldValue.SetInt(v.Int())
			case reflect.Bool:
				if v.Type().Kind() == reflect.String {
					res, err := strconv.ParseBool(v.String())
					if err != nil {
						panic(err)
					}
					fieldValue.SetBool(res)
					continue
				}
			default:
				if fieldValue.Type() == v.Type() {
					fieldValue.Set(v)
				} else {
					panic(fmt.Sprintf(`type mismatch for field "%s"`, fieldName))
				}
			}
		}
	}
}

func initSetting(setting any) {
	rv := reflect.Indirect(reflect.ValueOf(setting))

	for i := 0; i < rv.NumField(); i++ {
		value := rv.Field(i)

		component, ok := value.Interface().(interface{ Init() error })
		if ok {
			if err := component.Init(); err != nil {
				panic(err)
			}
		}
	}
}
