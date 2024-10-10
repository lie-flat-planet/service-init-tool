package service_init_tool

import (
	"fmt"
	"github.com/lie-flat-planet/service-init-tool/config_source"
	"github.com/lie-flat-planet/service-init-tool/enum"
	"github.com/lie-flat-planet/service-init-tool/envvar"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

func Init(dir string, setting any) error {
	_, filename, _, _ := runtime.Caller(1)
	dir = filepath.Join(filepath.Dir(filename), dir)

	cfg := &configuration{
		dir:    dir,
		env:    GetEnv(),
		parser: envvar.NewParser(setting),
	}

	return cfg.initSetting(setting)
}

// TODO
func ListServiceUpstream() {
	// return c.listServiceUpstream()
}

func GetEnv() string {
	env := os.Getenv(enum.EnvKey)
	return strings.ToUpper(env)
}

type configuration struct {
	dir string
	env string
	// 解析出拍平的环境变量
	parser *envvar.Parser
	// TODO 解析出依赖
	// merge环境变量
	envVarMerger *envvar.Merger
}

// initSetting 将环境变量注入进结构体
func (conf *configuration) initSetting(setting any) error {
	tpe := reflect.TypeOf(setting)
	if tpe.Kind() != reflect.Ptr {
		return fmt.Errorf("please pass ptr for setting value")
	}

	// 解析并生成开发环境的环境变量参考文件
	if err := conf.parser.FlattenEnvVar(conf.dir); err != nil {
		return err
	}
	conf.injectEnvVarMerger()

	// 获取合并后的环境变量最终值
	configValue, err := conf.getConfigValue()
	if err != nil {
		return err
	}

	// 设置字段
	setSettingFields(setting, configValue)

	// 初始化
	initSetting(setting)

	return nil
}

func (conf *configuration) getConfigValue() (configValue map[string]any, err error) {
	return conf.envVarMerger.Action()
}

func (conf *configuration) injectEnvVarMerger() {
	conf.envVarMerger = envvar.NewMerger(conf.parser.GetFlattenedEnvVarKeys(), conf.listEnvVarSource()...)
}

func (conf *configuration) listEnvVarSource() []config_source.ISource {
	list := []config_source.ISource{config_source.NewEnvVar()}

	if conf.env == enum.EnvDev || conf.env == "" {
		localYMLPath := conf.dir + "/local.yml"

		_, err := os.Stat(localYMLPath)
		if err == nil {
			list = append(list, config_source.NewYamlFile(localYMLPath))
		} else {
			if !os.IsNotExist(err) {
				panic(err)
			}
		}
	}

	return list
}

// TODO
func (conf *configuration) listServiceUpstream() {

}
