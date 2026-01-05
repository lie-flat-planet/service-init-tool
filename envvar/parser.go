package envvar

import (
	"github.com/lie-flat-planet/service-init-tool/util"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"reflect"
)

type Parser struct {
	setting         any
	flattenedEnvVar map[string]any
}

func NewParser(setting any) *Parser {
	return &Parser{
		setting:         setting,
		flattenedEnvVar: make(map[string]any),
	}
}

func (p *Parser) GetFlattenedEnvVar() map[string]any {
	return p.flattenedEnvVar
}

func (p *Parser) GetFlattenedEnvVarKeys() map[string]struct{} {
	m := make(map[string]struct{})

	for k := range p.flattenedEnvVar {
		m[k] = struct{}{}
	}

	return m
}

func (p *Parser) GenerateEnvVarTemplate(dir string) error {
	bytes, err := structEnvVar(p.setting)
	if err != nil {
		return err
	}

	var values = make(map[string]any)
	if err = yaml.Unmarshal(bytes, &values); err != nil {
		return err
	}

	util.FlattenMap("", values, p.flattenedEnvVar)

	if err = p.generateDevYML(dir); err != nil {
		return err
	}

	return nil
}

func (p *Parser) generateDevYML(dir string) error {
	filename := "dev.yml"

	filename = filepath.Join(dir, filename)

	b, err := yaml.Marshal(p.GetFlattenedEnvVar())
	if err != nil {
		return err
	}

	return util.CreateFile(filename, b)

}

// 将带有 env tag 的字段结构化
func structEnvVar(v any) ([]byte, error) {
	val := reflect.ValueOf(v).Elem()
	typ := val.Type()

	data := make(map[string]any)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if fieldType.Anonymous {
			anonymousData, _ := structEnvVar(field.Addr().Interface())
			var anonymousMap map[string]any
			_ = yaml.Unmarshal(anonymousData, &anonymousMap)

			for anonymousK, anonymousV := range anonymousMap {
				data[anonymousK] = anonymousV
			}
			continue
		}

		if _, skipOK := fieldType.Tag.Lookup("skip"); skipOK {
			continue
		}

		_, tagOK := fieldType.Tag.Lookup("env")
		if field.Kind() != reflect.Ptr && field.Kind() != reflect.Struct && !tagOK {
			continue
		}

		fieldName := fieldType.Name

		switch field.Kind() {
		case reflect.Struct:
			nestedData, _ := structEnvVar(field.Addr().Interface())
			var nestedMap map[string]any
			_ = yaml.Unmarshal(nestedData, &nestedMap)
			data[fieldName] = nestedMap
		case reflect.Ptr:
			if field.IsNil() {
				field.Set(reflect.New(fieldType.Type.Elem()))
			}
			nestedData, _ := structEnvVar(field.Interface())
			var nestedMap map[string]any
			_ = yaml.Unmarshal(nestedData, &nestedMap)
			data[fieldName] = nestedMap
		default:
			data[fieldName] = field.Interface()
		}
	}

	return yaml.Marshal(data)
}
