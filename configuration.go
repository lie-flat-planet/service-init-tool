package service_init_tool

import (
	"fmt"
	"github.com/lie-flat-planet/service-init-tool/config_source"
	"github.com/lie-flat-planet/service-init-tool/envvar"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	EnvDev string = "DEV"
	EnvPro string = "PRO"
)

var c *Configuration

func Init(svcName, dir string, setting any) {

}

func ListServiceUpstream() {
	// return c.listServiceUpstream()
}

func GetEnv() string {
	env := os.Getenv("ENV")
	return strings.ToUpper(env)
}

type Configuration struct {
	dir string
	// 当前环境
	env string
	// 解析出拍平的环境变量
	parser *envvar.Parser
	// 解析出依赖
	// merge环境变量 done
	merger *envvar.Merger
	// 将环境变量注入进结构体 (initSetting done)

	svc *Service
}

func (conf *Configuration) initSetting(setting any) error {
	tpe := reflect.TypeOf(setting)
	if tpe.Kind() != reflect.Ptr {
		return fmt.Errorf("please pass ptr for setting value")
	}

	// 解析并生成开发环境的环境变量参考文件
	if err := conf.parser.FlattenEnvVar(conf.dir); err != nil {
		return err
	}

	// 获取合并后的环境变量最终值
	configValue, err := conf.getConfigValue()
	if err != nil {
		return err
	}

	// 设置字段
	setFields(setting, configValue)

	return nil
}

func (conf *Configuration) getConfigValue() (configValue map[string]any, err error) {
	m := envvar.NewMerger(conf.parser.GetFlattenedEnvVarKeys(), config_source.NewEnvVar(), conf.getLocalYMLSource())

	return m.Action()
}

func (conf *Configuration) getLocalYMLSource() config_source.ISource {
	if conf.env == EnvDev || conf.env == "" {
		return config_source.NewYamlFile(conf.dir + "/local.yml")
	}

	return nil
}

// TODO
func (conf *Configuration) listServiceUpstream() {

}

// 设置字段
func setFields(setting any, configValues map[string]any) {
	structValue := reflect.ValueOf(setting).Elem()
	structType := structValue.Type()

	for i := 0; i < structValue.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)

		if field.Anonymous {
			setFields(fieldValue.Addr().Interface(), configValues)
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
					setFields(fieldValue.Interface(), strMapVal)
				}
			case reflect.Struct:
				if mapVal, ok := val.(map[string]any); ok {
					strMapVal := make(map[string]any)
					for k, v := range mapVal {
						strMapVal[k] = v
					}
					setFields(fieldValue.Addr().Interface(), strMapVal)
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
