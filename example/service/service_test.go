package service

import (
	service_init_tool "github.com/lie-flat-planet/service-init-tool"
	"github.com/lie-flat-planet/service-init-tool/component"
	"testing"
)

type Config struct {
	Mysql *component.Mysql
	Name  string `env:""`
	Age   uint   `env:""`

	Goods *Goods
}

type Goods struct {
	Number string `env:""`
	Size   int    `env:""`
}

var Setting = &Config{
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

func TestInit(t *testing.T) {
	err := service_init_tool.Init("demo", "./", Setting)
	if err != nil {
		t.Fatal(err)
	}
}
