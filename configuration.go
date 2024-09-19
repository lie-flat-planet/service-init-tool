package service_init_tool

import (
	"fmt"
	"os"
	"reflect"
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

type Configuration struct {
	env string

	flatKeys  map[string]struct{}
	dir       string
	svc       *Service
	assistant *FileAssistant
}

func (conf *Configuration) initSetting(setting any) error {
	tpe := reflect.TypeOf(setting)
	if tpe.Kind() != reflect.Ptr {
		return fmt.Errorf("please pass ptr for setting value")
	}

	if err := conf.assistant.GenerateDevYML(); err != nil {
		return err
	}

	configValue, err := conf.getConfigValue()
	if err != nil {
		return err
	}

	conf.setFields(setting, configValue)

	return nil
}

// 设置字段
func (conf *Configuration) setFields(setting any, configValue map[string]any) {

}

// 获取结构化的环境变量值
func (conf *Configuration) getConfigValue() (configValue map[string]any, err error) {
	return
}

// TODO
func (conf *Configuration) listServiceUpstream() {

}

func GetEnv() string {
	env := os.Getenv("ENV")
	return strings.ToUpper(env)
}

type FileAssistant struct {
	flatEnvVar map[string]any
}

// TODO
func (assistant *FileAssistant) GenerateDevYML() error {
	return nil
}

// TODO
func (assistant *FileAssistant) GetLocalYMLContent() {

}
