package service

import (
	service_init_tool "github.com/lie-flat-planet/service-init-tool"
	"github.com/lie-flat-planet/service-init-tool/component"
	"github.com/lie-flat-planet/service-init-tool/util"
	"github.com/sirupsen/logrus"
	"testing"
)

type Config struct {
	Server *service_init_tool.Server
	Mysql  *component.Mysql
	Name   string `env:""`
	Age    uint   `env:""`

	Goods *Goods
}

type Goods struct {
	Number string `env:""`
	Size   int    `env:""`
}

var Setting = &Config{
	Server: &service_init_tool.Server{
		Name:     "demo",
		LogLevel: "DEBUG",
		HttpPort: 80,
	},

	Mysql: &component.Mysql{
		MySqlConfig: component.MySqlConfig{
			Host:        "11",
			User:        "22",
			Password:    "33",
			DbName:      "44",
			MaxIdleConn: 5,
			MaxOpenConn: 10,
		},
	},
	Name: "xiaoxlm",
	Age:  30,
	Goods: &Goods{
		Number: "qa123445678",
		Size:   100,
	},
}

// service_init_tool.Init 会生成环境变量的 kv 模版 dev.yml。
// 增加环境变量，环境变量的值会覆盖 Setting 中的内容。
// local.yml 里面的内容会覆盖 Setting 中的内容。即覆盖优先级 local.yml > 环境变量 > Setting
func TestInit(t *testing.T) {
	err := service_init_tool.Init("./", Setting)
	if err != nil {
		t.Fatal(err)
	}

	logrus.Info("哈哈哈")

	util.LogJSON(Setting)
}
