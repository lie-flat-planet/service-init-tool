package service_init_tool

import (
	"github.com/lie-flat-planet/service-init-tool/enum"
	"github.com/lie-flat-planet/service-init-tool/log"
)

type Server struct {
	Name string
	// service code
	Code int
	// LogLevel 默认是 debug 等级
	LogLevel string `env:""`
	HttpPort uint   `env:""`
	RunMode  string `env:""`
}

func (s *Server) GetHttpPort() uint {
	return s.HttpPort
}

func (s *Server) Init() error {
	l := &log.Log{
		Name:  s.Name,
		Level: s.LogLevel,
	}
	if l.Level == "" {
		l.Level = enum.Debug
	}

	l.SetDefaults().Build()

	return nil
}
