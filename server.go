package service_init_tool

import "github.com/lie-flat-planet/service-init-tool/log"

type Server struct {
	Name string
	// service code
	Code     int
	LogLevel string `env:""`
	HttpPort uint   `env:""`
	RunMode  string `env:""`
}

func (s *Server) GetHttpPort() uint {
	return s.HttpPort
}

func (s *Server) Init() {
	l := &log.Log{
		Name:  s.Name,
		Level: s.LogLevel,
	}

	l.SetDefaults().Build()
}
