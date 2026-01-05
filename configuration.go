package service_init_tool

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/lie-flat-planet/service-init-tool/config_source"
	"github.com/lie-flat-planet/service-init-tool/enum"
	"github.com/lie-flat-planet/service-init-tool/envvar"
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
	if err := conf.parser.GenerateEnvVarTemplate(conf.dir); err != nil {
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
	var envVarSources []config_source.ISource

	// 优先级，1为最高
	if conf.env == enum.EnvTest {
		envVarSources = append(envVarSources, conf.envVarsFromTestYML()) // 5
	}

	if conf.env == enum.EnvStaging {
		envVarSources = append(envVarSources, conf.envVarsFromStagingYML()) // 4
	}

	envVarSources = append(envVarSources, conf.envVarsFromHotFixYML()) // 3

	envVarSources = append(envVarSources, config_source.NewEnvVar()) // 2

	if conf.env == enum.EnvDev || conf.env == "" {
		envVarSources = append(envVarSources, conf.envVarsFromLocalYML()) // 1
	}

	conf.envVarMerger = envvar.NewMerger(conf.parser.GetFlattenedEnvVarKeys(), envVarSources...)
}

// 1
func (conf *configuration) envVarsFromLocalYML() *config_source.YamlFile {
	localYMLPath := conf.dir + "/local.yml"

	return conf.formYamlFile(localYMLPath)
}

// 2 envvar (线上建议使用 envvar 的配置方式)

// 3
func (conf *configuration) envVarsFromHotFixYML() *config_source.YamlFile {
	featureYMLPath := conf.dir + "/hot-fix.yml"

	return conf.formYamlFile(featureYMLPath)
}

// 4
func (conf *configuration) envVarsFromStagingYML() *config_source.YamlFile {
	stagingYMLPath := conf.dir + "/staging.yml"

	return conf.formYamlFile(stagingYMLPath)
}

// 5
func (conf *configuration) envVarsFromTestYML() *config_source.YamlFile {
	testYMLPath := conf.dir + "/test.yml"

	return conf.formYamlFile(testYMLPath)
}

func (conf *configuration) formYamlFile(filepath string) *config_source.YamlFile {
	_, err := os.Stat(filepath)
	if err == nil {
		return config_source.NewYamlFile(filepath)
	}

	if !os.IsNotExist(err) {
		panic(err)
	}

	return nil
}

// TODO
func (conf *configuration) listServiceUpstream() {

}
