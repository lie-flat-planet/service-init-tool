package service_init_tool

import "testing"

func Test_setFields(t *testing.T) {

	setFields(ple, map[string]any{})
}

var ple = &People{
	Name:   "liming",
	Age:    18,
	Height: 180,
	Belong: "四川",
}

type Config struct {
	Self *People
}

type People struct {
	Name   string `env:""`
	Age    int
	Height int
	Belong string `env:""`
}

type Home struct {
	Province string `env:""`
	City     string `env:""`
	Area     string `env:""`
}
