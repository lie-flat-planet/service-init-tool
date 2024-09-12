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

	dir string
	svc *Service
}

func (conf *Configuration) initSetting(setting any) (err error) {
	tpe := reflect.TypeOf(setting)
	if tpe.Kind() != reflect.Ptr {
		return fmt.Errorf("please pass ptr for setting value")
	}

	configValue, err := conf.getConfigValue()
	if err != nil {
		return err
	}

	conf.setFields(setting, configValue)

	return
}

// 返回服务的依赖信息
func (conf *Configuration) listServiceUpstream() {

}

// 设置字段
func (conf *Configuration) setFields(setting any, configValue map[string]any) {

}

// 获取结构化的环境变量值
func (conf *Configuration) getConfigValue() (configValue map[string]any, err error) {
	return
}

func GetEnv() string {
	env := os.Getenv("ENV")
	return strings.ToUpper(env)
}
