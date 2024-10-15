package service_init_tool

import (
	"fmt"
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
	if err := s.check(); err != nil {
		return err
	}

	s.init()

	return nil
}

func (s *Server) init() {
	if s.LogLevel == "" {
		s.LogLevel = enum.Debug
	}

	if s.HttpPort == 0 {
		s.HttpPort = 80
	}

	if s.RunMode == "" {
		s.RunMode = enum.Debug
	}

	// log
	(&log.Log{
		Name:  s.Name,
		Level: s.LogLevel,
	}).SetDefaults().Build()
}

func (s *Server) check() error {
	if s.Name == "" {
		return fmt.Errorf("server Name cannot be empty")
	}

	if s.Code == 0 {
		return fmt.Errorf("server Code cannot be zero")
	}

	return nil
}
